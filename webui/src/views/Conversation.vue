<template>

    <div class="conversation">

        <!-- header pagina -->
        <div class="header">
            <img :src="conversationImage" alt="Foto profilo" style="width: 40px; height:40px; border-radius: 10px;" />
            <button v-if="conversation.is_group" @click="groupPhotoPopup = true">
                <svg class="feather">
                    <use href="/feather-sprite-v4.29.0.svg#edit" />
                </svg>
            </button>

            <h1 id="convName">{{ conversation.name }}</h1>
            <button class="leave" v-if="conversation.is_group" @click="leaveGroup">
                <svg class="feather">
                    <use href="/feather-sprite-v4.29.0.svg#x" />
                </svg>
                Esci dal gruppo
            </button>
        </div>

        <!-- conversazione -->
        <div class="messages" id="messagesContainer">
            <div v-if="!forwardPopup && !photoPopup && !groupPhotoPopup && !reactionPopup"
                v-for="(message, index) in messages" :key="index"
                :class="['message', message.sender === name ? 'user' : 'receiver']">
                <div v-if="message.replied_message_id">
                    <p class="replied-message">
                        {{ getMessageById(message.replied_message_id).sender }}:
                        {{ (getMessageById(message.replied_message_id).photo ==
                            false) ? getMessageById(message.replied_message_id).content : "photo" }}
                    </p>
                </div>

                <!-- contenuto del messaggio in caso di testo -->
                <div v-if="!message.photo" class="normal-message">

                    <p class="sender">{{ message.sender }}</p>
                    <p class="content">{{ message.content }}</p>
                    <p class="timestamp">{{ new Date(message.timestamp).toLocaleString() }}</p>


                </div>

                <!-- contenuto del messaggio in caso di immagine -->
                <div v-else class="photo-message">

                    <p class="sender">{{ message.sender }}</p>
                    <img :src="message.photoUrl" alt="Foto del messaggio"
                        style="max-width: 100%; border-radius: 10px;" />
                    <p class="timestamp">{{ new Date(message.timestamp).toLocaleString() }}</p>

                </div>

                <!-- footer del messaggio  -->
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

                    <button class="reply-btn" @click="replyMessage(message)">
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
                    <button>
                        <svg class="feather" @click="reactMessage(message.message_id)">
                            <use href="/feather-sprite-v4.29.0.svg#plus-circle" />
                        </svg>
                    </button>
                </div>
            </div>

            <!-- popup per inoltro messaggi -->
            <div class="backpopup" v-if="forwardPopup">
                <div class="popup">

                    <div class="popup-header">
                        <h3>invia a</h3>
                        <button id="exit-button" @click="forwardPopup = false">
                            <svg class="feather">
                                <use href="/feather-sprite-v4.29.0.svg#x" />
                            </svg>
                        </button>
                    </div>
                    <div class="popup-content">
                        <div v-if="groups.length != 0" class="popup-content-conversations">
                            <h4>Gruppi</h4>
                            <ul>
                                <li v-for="group in groups" class="users">
                                    <label>
                                        <input type="checkbox" :value="conv" v-model="forwardGroupsList" />
                                        {{ group.name }}
                                    </label>
                                </li>
                            </ul>
                        </div>
                        <div class="popup-content-users">
                            <h4>Utenti</h4>
                            <ul>
                                <li v-for="user in users" class="users">
                                    <label>
                                        <input type="checkbox" :value="user" v-model="forwardUsersList" />
                                        {{ user.name }}
                                    </label>
                                </li>
                            </ul>

                        </div>

                    </div>
                    <div class="popup-footer">
                        <button @click="forwardMessage">
                            invia
                        </button>
                    </div>


                </div>
            </div>

            <!-- popup per invio di foto -->
            <div class="backpopup" v-if="photoPopup">
                <div class="popup">
                    <div class="popup-header">
                        <h3>Carica una foto</h3>
                        <button id="exit-button" @click="photoPopup = false">
                            <svg class="feather">
                                <use href="/feather-sprite-v4.29.0.svg#x" />
                            </svg>
                        </button>
                    </div>
                    <div class="popup-content">
                        <input type="file" @change="uploadImage" accept="image/*" />
                    </div>
                    <div class="popup-footer">
                        <button @click="sendImage">
                            invia
                        </button>
                    </div>
                </div>
            </div>

            <!-- popup per impostare la foto di un gruppo -->
            <div class="backpopup" v-if="groupPhotoPopup">
                <div class="popup">
                    <div class="popup-header">
                        <h3>Carica una foto</h3>
                        <button id="exit-button" @click="groupPhotoPopup = false">
                            <svg class="feather">
                                <use href="/feather-sprite-v4.29.0.svg#x" />
                            </svg>
                        </button>
                    </div>
                    <div class="popup-content">
                        <input type="file" @change="" accept="image/*" />
                    </div>
                    <div class="popup-footer">
                        <button @click="">
                            imposta immagine del gruppo
                        </button>
                    </div>
                </div>
            </div>

            <!-- popup per aggiungere una reazione -->
            <div class="backpopup" v-if="reactionPopup">
                <div class="popup">
                    <div class="popup-header">
                        <h3>aggiungi una reazione</h3>
                        <button id="exit-button" @click="reactionPopup = false">
                            <svg class="feather">
                                <use href="/feather-sprite-v4.29.0.svg#x" />
                            </svg>
                        </button>
                    </div>
                    <div class="popup-content">
                        <ul id="reaction">
                            <li v-for="reaction in ['👍', '❤️', '😂', '😮', '😢', '😡']" :key="reaction">
                                <button @click="this.reactionContent = reaction">
                                    {{ reaction }}
                                </button>
                            </li>
                        </ul>
                    </div>
                    <div class="popup-footer">
                        <button @click=" reactMessage(reactMessageID)">
                            imposta reazione al messaggio
                        </button>
                    </div>
                </div>
            </div>

        </div>

        <!-- barra del messaggio risposto -->
        <div v-if="replyedMessage.replied_message_id" class="reply-info">
            <p>
                risposta a: {{ (replyedMessage.photo == 0) ? replyedMessage.content : "photo" }}
            </p>
            <button @click="replyedMessage.replied_message_id = null; replyedMessage.content = ''">Annulla</button>
        </div>

        <!-- barra di inserimento -->
        <div class="input-area">
            <input v-model="newMessage" type="text" placeholder="Scrivi un messaggio..." @keyup.enter="sendMessage" />
            <button id="photo" @click="photoPopup = true">
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

        <!-- lista degli utenti  -->
        <MembersList v-if="conversation.is_group" :conversation="conversation" :conversationID="conversationID" />



    </div>
</template>

<script>
import MembersList from '@/components/MembersList.vue';
export default {
    name: 'Conversation',
    components: {
        MembersList,
    },
    data() {
        return {
            conversation: {},
            conversationImage: '',
            groups: [],
            users: [],
            messages: [],

            newMessage: '',
            replyedMessage: {},

            forwardPopup: false,
            photoPopup: false,
            groupPhotoPopup: false,
            reactionPopup: false,

            photoUrl: '',
            selectedFile: null,
            previewImage: null,
            forwardGroupsList: [],
            forwardUsersList: [],
            reactionContent: '',
            reactMessageID: null,

            forwardMessageId: null,

            conversationID: this.$route.params.conversationID,
            token: localStorage.getItem("token"),
            name: localStorage.getItem("name"),
            username: localStorage.getItem("username"),
            user_id: localStorage.getItem("user_id"),

            mesage: '',
            error: '',
        }
    },
    methods: {
        async getConversation() {
            try {
                const res = await this.$axios.get(`/conversation/${this.conversationID}`, {
                    headers: { Authorization: `Bearer ${this.token}` }
                })
                this.conversation = res.data
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
        async fetchUsers() {
            this.loading = true;
            try {
                const response = await this.$axios.get("/user", {
                    headers: {
                        Authorization: `Bearer ${this.token}`,
                    },
                });
                this.users = response.data;
            } catch (err) {
                this.error = "Errore nel recupero degli utenti";
            } finally {
                this.loading = false;
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
            for (const message of this.messages) {
                if (message.photo == true) {
                    try {
                        const res = await this.$axios.get(`/conversation/${message.conversation_id}/photo/${message.message_id}`, {
                            headers: { Authorization: `Bearer ${this.token}` },
                            responseType: 'blob'
                        })
                        message.photoUrl = URL.createObjectURL(res.data);
                    } catch (err) {
                        console.error('Errore nel caricamento della foto del messaggio:', err)
                    }
                }
                this.$nextTick(() => {
                    this.scrollToBottom()
                })


            }
        },
        uploadImage(e) {
            const file = e.target.files[0];
            this.selectedFile = file;   // salva il file originale
            const reader = new FileReader();
            reader.readAsDataURL(file);
            reader.onload = e => {
                this.previewImage = e.target.result;  // questa è la base64 per l'anteprima
                console.log(this.previewImage);
            };
        },
        async sendImage() {
            this.message = '';
            this.error = '';

            if (!this.selectedFile) {
                this.error = 'Seleziona un file prima di caricare.';
                return;
            }

            const formData = new FormData();
            formData.append('photo', this.selectedFile);  // metti il file originale qui

            if (this.replyedMessage.replied_message_id) {
                formData.append('reply_message_id', this.replyedMessage.message_id);
            }


            try {
                const response = await this.$axios.post(`/conversation/${this.conversationID}/photo`, formData,
                    {

                        headers: {
                            Authorization: `Bearer ${this.token}`,
                            "Content-Type": "multipart/form-data",
                        },
                    });
                this.message = 'Foto profilo caricata con successo!';
            } catch (err) {
                this.error = 'Errore durante il caricamento della foto profilo.';
            }
            this.photoPopup = false;
            this.selectedFile = null;  // resetta il file selezionato
            this.scrollToBottom()
            this.fetchMessages()
        },
        async sendMessage() {

            const messageToSend = {
                replied_message_id: this.replyedMessage.replied_message_id,
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

                this.fetchMessages()


            } catch (err) {
                console.error("Errore durante l'invio del messaggio:", err)
            }

            this.replyedMessage.replied_message_id = null
            this.replyedMessage.content = ''

        },
        replyMessage(message) {
            this.replyedMessage.replied_message_id = message.message_id
            this.replyedMessage.content = message.content
        },
        async fetchMyConversations() {
            this.loading = true;
            try {
                const response = await this.$axios.get("/conversation", {
                    headers: {
                        Authorization: `Bearer ${this.token}`,
                    },
                });
                this.groups = response.data;
            } catch (err) {
                this.error = "Errore nel recupero degli utenti";
            } finally {
                this.loading = false;
            }
            for (const conv of this.groups) {
                if (conv.name === this.name) {
                    this.users.splice(this.users.indexOf(user), 1);
                }
            }
            this.groups = this.groups.filter(conv => conv.is_group === true)
        },
        async forwardMessage(messageId) {
            if (this.forwardPopup === false) {
                this.forwardPopup = true
                this.forwardMessageId = messageId
                await this.fetchMyConversations()
                await this.fetchUsers()
                return
            } else {
                for (const conv of this.forwardGroupsList) {
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

                for (const user of this.forwardUsersList) {
                    try {
                        const response = await this.$axios.post("/conversation",
                            { name: user.name },
                            {
                                headers: {
                                    Authorization: `Bearer ${this.token}`,
                                },
                            }
                        );

                        const conversation = response.data;
                        try {
                            await this.$axios.post(
                                `/message/${this.forwardMessageId}/forwarded`,
                                {
                                    receiver_id: conversation.conversation_id,
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
                    } catch (err) {
                        if (err.response?.status === 409 && err.response?.data?.conversation_id) {
                            try {
                                await this.$axios.post(
                                    `/message/${this.forwardMessageId}/forwarded`,
                                    {
                                        receiver_id: err.response.data.conversation_id,
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
                        } else {
                            console.error("Errore nella creazione della conversazione", err);
                            this.error = "Impossibile avviare la conversazione.";
                        }
                    }

                }
                this.forwardGroupsList = []
                this.forwardUsersList = []
                this.forwardMessageId = null
                this.forwardPopup = false

            }

        },
        async reactMessage(messageId) {
            if (this.reactionPopup === false) {
                this.reactionPopup = true
                this.reactMessageID = messageId
                return
            }
            else {
                try {
                    await this.$axios.post(`/message/${messageId}/reaction`, 
                    {
                        user_id: this.user_id,
                        content: this.reactionContent,
                        message_id: messageId,
                    },
                    { headers: { Authorization: `Bearer ${this.token}` } }
                    )
                    this.reactionContent = ''
                } catch (err) {
                    console.error("Errore durante l'invio della reazione:", err)
                }
            }


        },
        closeForwardMessage() {
            this.forwardGroupsList = []
            this.forwardMessageId = null
            this.forwardPopup = false
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
    beforeUnmount() {
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

    height: 400px;
    margin-bottom: 10px;
    background-color: #f4f6f8;
    display: flex;
    flex-direction: column;
    overflow-y: auto;
    border-radius: 23px;
}


.message {
    margin: 0px 15px;
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

.photo-message {
    width: 200px;
    height: 200px;
    display: flex;
    flex-direction: column;
    align-items: center;
    border-radius: 10px;
}

.photo-message img {
    max-height: 600px;

    border-radius: 10px;
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
    border: #888 1px solid;
    border-radius: 10px;
    padding: 5px;
}

.input-area {
    display: flex;
    gap: 10px;
}



.backpopup {
    position: relative;
    width: 100%;
    height: 100%;
    background-color: rgba(138, 137, 137, 0.87);
    display: flex;
    justify-content: center;
    align-items: center;

}

.users {
    text-align: center;
    display: flex;
    margin-bottom: 10px;
    border: 1px solid #888;
    border-radius: 10px;
    padding: 10px;
    background-color: #d1e7dd;

}

.popup {
    border: 1px solid #888;
    padding: 10px;
    margin: 10px 0px;
    background-color: #f4f6f8;
    display: flex;
    flex-direction: column;
    overflow-y: auto;
    border-radius: 23px;
    position: relative;
    width: 80%;
    height: 80%;
}

.popup-header {
    display: flex;
    justify-content: space-between;
    margin-bottom: 10px;
}

.popup-content {
    display: flex;
    flex-direction: column;
    overflow-y: auto;
}

.popup-content ul {
    list-style: none;
    padding: 0;
    margin: 0;
}

#reaction {
    font-size: 40px;
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    gap: 10px;
}

.popup-footer {
    position: bottom fixed;
    display: flex;
    justify-content: center;
    margin-top: 10px;
}

input {
    flex: 1;
    padding: 5px;
    border: 1px solid #888;
    border-radius: 15px;
    background-color: #d1e7dd;
    color: black;
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
