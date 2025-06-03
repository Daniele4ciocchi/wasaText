<template>
    <div class="conversation">
        <div class="header">
            <img :src="conversationImage" alt="Foto profilo" style="width: 40px; height:40px; border-radius: 10px;" />
            <h1 id="convName">{{ conversation.name }}</h1>
            <button class="leave" v-if="conversation.is_group" @click="leaveGroup">
                <svg class="feather">
                    <use href="/feather-sprite-v4.29.0.svg#x" />
                </svg>
                Esci dal gruppo
            </button>
        </div>

        <div class="messages" id="messagesContainer">
            <div v-for="(message, index) in messages" :key="index"
                :class="['message', message.sender === name ? 'user' : 'receiver']">
                <div v-if="message.replied_message_id">
                    <p class="replied-message">
                        {{ getMessageById(message.replied_message_id).sender }}:
                        {{ getMessageById(message.replied_message_id).content }}
                    </p>
                </div>
                <p class="sender">{{ message.sender }}</p>
                <p class="content">{{ message.content }}</p>
                <p class="timestamp">{{ new Date(message.timestamp).toLocaleString() }}</p>

                <div class="message-footer">
                    <div v-if="message.sender === name" class="view">
                        <svg v-if="message.status === 0" class="feather">
                            <use href="/feather-sprite-v4.29.0.svg#check" />
                        </svg>
                        <svg v-if="message.status === 1" class="feather">
                            <use href="/feather-sprite-v4.29.0.svg#check-circle" />
                        </svg>
                        <svg v-if="message.status === 2" class="feather">
                            <use href="/feather-sprite-v4.29.0.svg#check-square" />
                        </svg>

                    </div>

                    <button class="reply-btn" @click="replyMessage(message.message_id, message.content)">
                        <svg class="feather">
                            <use href="/feather-sprite-v4.29.0.svg#repeat" />
                        </svg>
                    </button>
                    <button v-if="message.sender === name" class="delete-btn"
                        @click="deleteMessage(message.message_id)">
                        <svg class="feather">
                            <use href="/feather-sprite-v4.29.0.svg#trash" />
                        </svg>
                    </button>
                    <button class="forward-btn" @click="forwardMessage(message.message_id)">
                        <svg class="feather">
                            <use href="/feather-sprite-v4.29.0.svg#corner-down-right" />
                        </svg>
                    </button>
                </div>

            </div>
        </div>

        <div v-if="replyID" class="reply-info">
            <p>
                risposta a: {{ replyContent }}
            </p>
            <button @click="replyID = null; replyContent = ''">Annulla</button>
        </div>

        <div class="input-area">
            <input v-model="newMessage" type="text" placeholder="Scrivi un messaggio..." @keyup.enter="sendMessage" />
            <button id="photo" @click="sendPhoto">
                <svg class="feather">
                    <use href="/feather-sprite-v4.29.0.svg#image" />
                </svg>
            </button>
            <button id="send-button" @click="sendMessage" :disabled="!newMessage.trim()">
                <svg class="feather">
                    <use href="/feather-sprite-v4.29.0.svg#send" />
                </svg>
            </button>
        </div>

        <div class="user-list" v-if="conversation.is_group">
            <h2>Utenti nella conversazione</h2>
            <ul>
                <li v-for="user in members" :key="user.name">
                    {{ user.name }}
                </li>
            </ul>
        </div>

        <div class="popup" v-if="forward">
            <div class="forward">

                <div class="forward-header">
                    <h3>invia a</h3>
                    <button @click="closeForwardMessage">
                        <svg class="feather">
                            <use href="/feather-sprite-v4.29.0.svg#x" />
                        </svg>
                    </button>
                </div>
                <ul>
                    <li v-for="conv in conversations" class="users">
                        <label>
                            <input type="checkbox" :value="conv" v-model="forwardList" />
                            {{ conv.name }}
                        </label>
                    </li>
                </ul>
                <button @click="forwardMessage">
                    invia
                </button>

            </div>
        </div>

    </div>
</template>

<script>
export default {
    name: 'Conversation',
    data() {
        return {
            conversationID: this.$route.params.conversationID,
            conversation: {},
            conversationImage: '',
            members: [],
            conversations: [],
            messages: [],
            newMessage: '',
            replyID: null,
            replyContent: '',
            forward: false,
            forwardList: [],
            forwardMessageId: null,
            token: localStorage.getItem("token"),
            name: localStorage.getItem("name"),
            username: localStorage.getItem("username"),
            user_id: localStorage.getItem("user_id"),
        }
    },
    methods: {
        async getConversation() {
            try {
                const res = await this.$axios.get(`/conversation/${this.conversationID}`, {
                    headers: { Authorization: `Bearer ${this.token}` }
                })
                this.conversation = res.data
                await this.fetchMembers()
                await this.fetchPhoto()
            } catch (err) {
                console.error('Errore nel caricamento della conversazione:', err)
            }
        },
        async fetchPhoto() {
            if (this.conversation.is_group) {
                try {
                    const res = await this.$axios.get(`/group/${this.conversationID}/photo`, {
                        headers: { Authorization: `Bearer ${this.token}` },
                        responseType: 'blob'
                    })
                    this.conversationImage = URL.createObjectURL(res.data);
                } catch (err) {
                    console.error('Errore nel caricamento della foto del gruppo:', err)
                }
            } else {
                try {
                    console.log("Utenti disponibili:", this.members)

                    const user = this.members.find(user => user.name === this.conversation.name);
                    const res = await this.$axios.get(`/user/${user.user_id}/photo`, {
                        headers: { Authorization: `Bearer ${this.token}` },
                        responseType: 'blob'
                    })
                    this.conversationImage = URL.createObjectURL(res.data);
                } catch (err) {
                    console.error('Errore nel caricamento della foto dell\'utente:', err)
                }
            }
        },
        async fetchMembers() {
            try {
                const res = await this.$axios.get(`/group/${this.conversationID}/member`, {
                    headers: { Authorization: `Bearer ${this.token}` }
                })
                this.members = res.data
            } catch (err) {
                console.error('Errore nel caricamento degli utenti:', err)
            }
        },
        async fetchMessages() {
            try {
                const res = await this.$axios.get(`/conversation/${this.conversationID}/message`, {
                    headers: { Authorization: `Bearer ${this.token}` }
                })
                this.messages = res.data
            } catch (err) {
                console.error('Errore nel caricamento dei messaggi:', err)
            }
        },
        async sendMessage() {

            const messageToSend = {
                replied_message_id: this.replyID,
                content: this.newMessage,
                sender: this.name,
                timestamp: new Date().toISOString()
            }

            try {
                await this.$axios.post(
                    `/conversation/${this.conversationID}/message`,
                    {
                        content: messageToSend.content,
                        replied_message_id: messageToSend.replied_message_id
                    },
                    { headers: { Authorization: `Bearer ${this.token}` } }
                )

                this.newMessage = ''
                this.$nextTick(() => {
                    this.scrollToBottom()
                    this.fetchMessages()
                })
            } catch (err) {
                console.error("Errore durante l'invio del messaggio:", err)
            }

            this.replyID = null
            this.replyContent = ''
        },
        replyMessage(messageId, messageContent) {
            this.replyID = messageId
            this.replyContent = messageContent
        },
        async fetchMyConversations() {
            this.loading = true;
            try {
                const response = await this.$axios.get("/conversation", {
                    headers: {
                        Authorization: `Bearer ${this.token}`,
                    },
                });
                this.conversations = response.data;
            } catch (err) {
                this.error = "Errore nel recupero degli utenti";
            } finally {
                this.loading = false;
            }
            for (const conv of this.conversations) {
                if (conv.name === this.name) {
                    this.users.splice(this.users.indexOf(user), 1);
                }
            }
        },
        async forwardMessage(messageId) {
            if (this.forward === false) {
                this.forward = true
                this.forwardMessageId = messageId
                this.fetchMyConversations()
                return
            } else {
                for (const conv of this.forwardList) {
                    try {
                        await this.$axios.post(
                            `/message/${this.forwardMessageId}/forwarded`,
                            {
                                receiver_id: conv.conversation_id,
                            },
                            { headers: { Authorization: `Bearer ${this.token}` } }
                        )

                        this.newMessage = ''
                        this.$nextTick(() => {
                            this.scrollToBottom()
                            this.fetchMessages()
                        })
                    } catch (err) {
                        console.error("Errore durante l'invio del messaggio:", err)
                    }
                }
                this.forwardList = []
                this.forwardMessageId = null
                this.forward = false

            }

        },
        closeForwardMessage() {
            this.forwardList = []
            this.forwardMessageId = null
            this.forward = false
        },
        async deleteMessage(messageId) {
            try {
                await this.$axios.delete(`/message/${messageId}`, {
                    headers: { Authorization: `Bearer ${this.token}` }
                })
                this.messages = this.messages.filter((m) => m.message_id !== messageId)
            } catch (err) {
                console.error("Errore durante l'eliminazione del messaggio:", err)
            }
            this.fetchMessages()
        },
        getMessageById(id) {
            return this.messages.find((m) => m.message_id === id) || {}
        },
        async leaveGroup() {
            try {
                await this.$axios.delete(`/group/${this.conversationID}`, {
                    headers: { Authorization: `Bearer ${this.token}` }
                })
                alert('Sei uscito dal gruppo')
                window.location.href = '/'
            } catch (err) {
                console.error("Errore durante l'uscita dal gruppo:", err)
            }
        },
        scrollToBottom() {
            const container = document.getElementById('messagesContainer')
            container.scrollTop = container.scrollHeight
        },
        async getNewMessages() {
            try {
                const res = await this.$axios.get('/me/newmessage', {
                    headers: { Authorization: `Bearer ${this.token}` }
                })
                if (res.data != null) {
                    await this.fetchMessages()
                    this.scrollToBottom()
                }
            } catch (err) {
                console.error('Errore nel caricamento dei nuovi messaggi:', err)
            }
        }
    },
    mounted() {
        this.getConversation()
        this.fetchMessages().then(this.scrollToBottom)
        this.intervalID = setInterval(() => {
            this.getNewMessages()
        }, 5000)
    },
    beforeDestroy() {
        clearInterval(this.intervalID)
    }

}
</script>

<style scoped>
#convName {
    font-size: 2em;
    text-align: left;
    margin: 15px 0px;
}

.conversation {
    width: 100%;
    margin: 0 auto;
    font-family: Arial, sans-serif;
}

.header {
    display: flex;
    align-items: center;
    gap: 20px;
}

.reply-info {
    display: flex;
    justify-content: left;
    margin-bottom: 10px;
    gap: 10px;
}

.reply-info p {
    margin: 0;
    padding: 5px;
    background-color: #d1e7dd;
    border-radius: 10px;
    border: #888 1px solid;
}

.leave {
    background-color: #fa716c;
    color: black;
    border-radius: 10px;
    cursor: pointer;
    border: #888 1px solid;
    padding: 5px 10px;
}

.leave:hover {
    background-color: #b02a2a;
}

#send-button {
    border-radius: 50px;
}


.messages {
    border: 1px solid #888;
    padding: 10px;
    height: 400px;
    margin-bottom: 10px;
    background-color: #f4f6f8;
    display: flex;
    flex-direction: column;
    overflow-y: auto;
    border-radius: 23px;
}


.message {
    margin-bottom: 10px;
    padding: 5px 15px;
    border-radius: 15px;
    width: auto;
    border: #888 1px solid;
    display: block;
}

.message-footer {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    gap: 10px;
}

.view {
    width: 15px;
    height: 15px;
    color: #6b6b6b;
    align-self: center;
}

.message.user {
    background-color: #d1e7dd;
    text-align: right;
    align-self: flex-end;
}

.message.receiver {
    background-color: #ffffff;
    text-align: left;
    align-self: flex-start;
}

.replied-message {
    background-color: #e2e3e5;
    border-radius: 10px;
    border: #888 1px solid;
    padding: 5px;
    margin: 5px 0;
}

.sender {
    font-weight: bold;
    margin: 0;
}

.content {
    margin: 5px 0;
}

.input-area {
    display: flex;
    gap: 10px;
}

.user-list {
    margin-top: 20px;
    padding: 10px;
    border: 1px solid #888;
    border-radius: 15px;
    background-color: #f4f6f8;
}

.popup {
    position: absolute;
    width: calc(100vw - 315px);
    height: 100vh;
    background-color: rgba(138, 137, 137, 0.87);
    display: flex;
    justify-content: center;
    align-items: center;
}

.users {
    text-align: center;
    width: 60%;
    display: flex;
    margin-bottom: 10px;
    border: 1px solid #888;
    border-radius: 10px;
    padding: 10px;
    background-color: #d1e7dd;

}

.forward-header {
    display: flex;
    justify-content: space-between;
}

.forward {
    border: 1px solid #888;
    padding: 10px;
    margin: 10px 0px;
    background-color: #f4f6f8;
    display: flex;
    flex-direction: column;
    overflow-y: auto;
    border-radius: 23px;
    position: absolute;
    width: 80%;
    height: 80%;
}

input {
    flex: 1;
    padding: 5px;
    border: 1px solid #888;
    border-radius: 15px;
    background-color: #d1e7dd;
    color: black;
}

button {
    padding: 5px 10px;
    border: none;
    background-color: #d7f2ba;
    color: black;
    border-radius: 10px;
    cursor: pointer;
    border: #888 1px solid;
}

button:hover {
    background-color: #bde4a8;
}

#replies {
    align-items: left;
    width: 100%;
    padding: 0px;
    margin: 0px;
}

.timestamp {
    font-size: 0.8em;
    color: #888;
}
</style>
