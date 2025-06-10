<template>
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
        token: {
            type: String,
            default: localStorage.getItem("token"),
        },
        name: {
            type: String,
            default: localStorage.getItem("name"),
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
        };
    },
    mounted() {

    },
    methods: {
    },
}
</script>