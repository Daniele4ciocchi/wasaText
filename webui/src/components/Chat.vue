<template>
    <!-- conversazione -->
    <div class="messages" id="messagesContainer">
        <div v-if="!forwardPopup && !photoPopup && !groupPhotoPopup && !reactionPopup"
            v-for="(message, index) in messages" v-bind:key="index"
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
                <img :src="message.photoUrl" alt="Foto del messaggio" style="max-width: 100%; border-radius: 10px;" />
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
                <button v-if="message.sender === name" class="delete-btn" @click="deleteMessage(message.message_id)">
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
                    <svg class="feather" @click="reactionPopup = true">
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
                            <li v-for="group in groups" class="users" v-bind:key="group.id"> >
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
                            <li v-for="user in users" v-bind:key="user.id" class="users">
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
                            <button @click="">
                                {{ reaction }}
                            </button>
                        </li>
                    </ul>
                </div>
                <div class="popup-footer">
                    <button @click="">
                        imposta reazione al messaggio
                    </button>
                </div>
            </div>
        </div>

    </div>

</template>
<script>
export default {
    name: 'Chat',
    props: {
        conversation: {
            type: Object,
            required: true,
        },
        conversationID: {
            type: [String, Number],
            required: true,
        },

    },
    data() {
        return {
            messages: [],

            forwardPopup: false,
            photoPopup: false,
            groupPhotoPopup: false,
            reactionPopup: false,

            forwardGroupsList: [],
            forwardUsersList: [],
            groups: [],
            users: [],

            token: localStorage.getItem("token"),
            name: localStorage.getItem("name"),
        };
    },
    mounted() {

    },
    methods: {
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
        replyMessage(message) {
            this.replyedMessage.replied_message_id = message.message_id
            this.replyedMessage.content = message.content
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

    },
}
</script>