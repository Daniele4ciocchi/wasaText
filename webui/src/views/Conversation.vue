<template>
    <div class="conversation">
        <div class="header">
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
                <button class="reply-btn" @click="replyMessage(message.message_id, message.content)">
                    <svg class="feather">
                        <use href="/feather-sprite-v4.29.0.svg#corner-down-right" />
                    </svg>
                </button>
            </div>
        </div>

        <div v-if="replyID" class="reply-info">
            risposta a: {{ replyContent }} (ID: {{ replyID }})
        </div>

        <div class="input-area">
            <input v-model="newMessage" type="text" placeholder="Scrivi un messaggio..." @keyup.enter="sendMessage" />
            <button id="send-button" @click="sendMessage">
                <svg class="feather">
                    <use href="/feather-sprite-v4.29.0.svg#send" />
                </svg>
            </button>
        </div>

        <div class="user-list" v-if="conversation.is_group">
            <h2>Utenti nella conversazione</h2>
            <ul>
                <li v-for="user in users" :key="user.name">
                    {{ user.name }}
                </li>
            </ul>
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
            users: [],
            messages: [],
            newMessage: '',
            replyID: null,
            replyContent: '',
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

                if (this.conversation.is_group) fetchUsers()
            } catch (err) {
                console.error('Errore nel caricamento della conversazione:', err)
            }
        },
        async fetchUsers() {
            try {
                const res = await this.$axios.get(`/group/${this.conversationID}`, {
                    headers: { Authorization: `Bearer ${this.token}` }
                })
                this.users = res.data.members
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
                this.messages.push(messageToSend)
                this.newMessage = ''
                this.$nextTick(() => {
                    this.scrollToBottom()
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
        getMessageById(id) {
            return this.messages.find((m) => m.message_id === id) || {}
        },
        async leaveGroup() {
            try {
                await this.$axios.delete(`/group/${conversationID.value}`, {
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
                    headers: { Authorization: `Bearer ${token}` }
                })

                if (res.data != null) {
                    await fetchMessages()
                    scrollToBottom()
                }
            } catch (err) {
                console.error('Errore nel caricamento dei nuovi messaggi:', err)
            }
        }
    },
    mounted() {
        this.getConversation()
        this.fetchMessages().then(this.scrollToBottom)
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
    justify-content: space-between;
    align-items: center;
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
    align-items: flex-end;
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

#timestamp {
    font-size: 0.8em;
    color: #888;
}
</style>
