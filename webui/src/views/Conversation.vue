<template>

    <div class="conversation">

        <ConversationHeader 
        :conversation="conversation"
        :conversation-image="conversationImage"
        @change-name="changeGroupName"
        @leave-group="leaveGroup"
        @change-photo="groupPhotoPopup = true"
    />

        <!-- conversazione -->
        <div class="container">
            <MessageList v-if="!forwardPopup && !photoPopup && !groupPhotoPopup && !reactionPopup"
                :messages="messages"
                :current-user-name="name"
                :current-user-id="user_id"
                @reply="replyMessage"
                @delete="deleteMessage"
                @forward="openForwardModal"
                @react="openReactionModal"
                @delete-reaction="deletereaction"
            />

            <!-- popup per inoltro messaggi -->
            <ForwardModal 
                :show="forwardPopup" 
                :groups="groups" 
                :users="users" 
                @close="forwardPopup = false"
                @forward="handleForward"
            />

            <!-- popup per invio di foto -->
            <PhotoUploadModal 
                :show="photoPopup" 
                title="Invia una foto"
                button-text="Invia"
                @close="photoPopup = false"
                @submit="sendImage"
            />

            <!-- popup per impostare la foto di un gruppo -->
            <PhotoUploadModal 
                :show="groupPhotoPopup"
                title="Cambia foto del gruppo"
                button-text="Imposta immagine"
                @close="groupPhotoPopup = false"
                @submit="changeGroupPhoto"
            />

            <!-- popup per aggiungere una reazione -->
            <ReactionModal 
                :show="reactionPopup"
                @close="reactionPopup = false"
                @react="handleReaction"
            />

        </div>

        <MessageInput 
            :replyed-message="replyedMessage"
            @send-message="sendMessage"
            @cancel-reply="replyedMessage = { replied_message_id: null, content: '' }"
            @open-photo-popup="photoPopup = true"
        />

        <!-- lista degli utenti  -->
        <MembersList v-if="conversation.is_group" :conversation="conversation" :conversationID="conversationID" />



    </div>
</template>

<script>
import MembersList from '@/components/MembersList.vue';
import ForwardModal from '@/components/ForwardModal.vue';
import MessageList from '@/components/MessageList.vue';
import ConversationHeader from '@/components/ConversationHeader.vue';
import MessageInput from '@/components/MessageInput.vue';
import PhotoUploadModal from '@/components/PhotoUploadModal.vue';
import ReactionModal from '@/components/ReactionModal.vue';

export default {
    name: 'Conversation',
    components: {
        MembersList,
        ForwardModal,
        MessageList,
        ConversationHeader,
        MessageInput,
        PhotoUploadModal,
        ReactionModal,
    },
    data() {
        return {
            // Core data
            conversation: {},
            messages: [],
            users: [],
            groups: [],
            
            // UI State
            conversationImage: '',
            replyedMessage: {},
            forwardPopup: false,
            photoPopup: false,
            groupPhotoPopup: false,
            reactionPopup: false,

            // IDs for modals
            forwardMessageId: null,
            reactMessageID: null,

            // Session/User data
            conversationID: this.$route.params.conversationID,
            token: localStorage.getItem("token"),
            name: localStorage.getItem("name"),
            user_id: localStorage.getItem("user_id"),

            // Error handling
            error: '',
        }
    },
    methods: {
        // --- DATA FETCHING ---
        async getConversation() {
            try {
                const res = await this.$axios.get(`/conversation/${this.conversationID}`, {
                    headers: { Authorization: `Bearer ${this.token}` }
                });
                this.conversation = res.data;

                if (!this.conversation.is_group) {
                    await this.fetchUsers();
                }
                await this.fetchPhoto();
            } catch (err) {
                console.error('Errore nel caricamento della conversazione:', err);
            }
        },
        async fetchPhoto() {
            if (this.conversation.is_group) {
                try {
                    const res = await this.$axios.get(`/group/${this.conversationID}/photo`, {
                        headers: { Authorization: `Bearer ${this.token}` },
                        responseType: 'blob'
                    });
                    this.conversationImage = URL.createObjectURL(res.data);
                } catch (err) {
                    console.error('Errore nel caricamento della foto del gruppo:', err);
                }
            } else {
                try {
                    const user = this.users.find(user => user.name === this.conversation.name);
                    if (user) {
                        const res = await this.$axios.get(`/user/${user.user_id}/photo`, {
                            headers: { Authorization: `Bearer ${this.token}` },
                            responseType: 'blob'
                        });
                        this.conversationImage = URL.createObjectURL(res.data);
                    }
                } catch (err) {
                    console.error('Errore nel caricamento della foto dell\'utente:', err);
                }
            }
        },
        async fetchUsers() {
            try {
                const response = await this.$axios.get("/user", {
                    headers: { Authorization: `Bearer ${this.token}` },
                });
                this.users = response.data;
            } catch (err) {
                this.error = "Errore nel recupero degli utenti";
            }
        },
        async fetchMessages() {
            try {
                const res = await this.$axios.get(`/conversation/${this.conversationID}/message`, {
                    headers: { Authorization: `Bearer ${this.token}` }
                });
                const messages = res.data;

                // Fetch photos and reactions in parallel for performance
                await Promise.all(messages.map(async (message) => {
                    if (message.photo) {
                        try {
                            const photoRes = await this.$axios.get(`/conversation/${message.conversation_id}/photo/${message.message_id}`, {
                                headers: { Authorization: `Bearer ${this.token}` },
                                responseType: 'blob'
                            });
                            message.photoUrl = URL.createObjectURL(photoRes.data);
                        } catch (err) {
                            console.error('Errore nel caricamento della foto del messaggio:', err);
                            message.photoUrl = ''; // or a placeholder
                        }
                    }
                    message.reactions = await this.getMessageReaction(message.message_id);
                }));

                this.messages = messages;

                this.$nextTick(() => {
                    this.scrollToBottom();
                });
            } catch (err) {
                console.error('Errore nel caricamento dei messaggi:', err);
            }
        },
        async getMessageReaction(messageId) {
            try {
                const res = await this.$axios.get(`/message/${messageId}/reaction`, {
                    headers: { Authorization: `Bearer ${this.token}` }
                });
                const reactions = res.data; // res.data dovrebbe essere un array

                if (!Array.isArray(reactions)) {
                    console.warn(`API /message/${messageId}/reaction non ha restituito un array:`, reactions);
                    return []; // Restituisce un array vuoto per evitare errori
                }

                return Object.values(reactions.reduce((acc, { content, user_id, reaction_id }) => {
                    if (!acc[content]) {
                        acc[content] = { content, users: [] };
                    }
                    acc[content].users.push({ user_id, reaction_id });
                    return acc;
                }, {}));
            } catch (err) {
                console.error("Errore nel recupero delle reazioni:", err);
                return [];
            }
        },

        // --- MESSAGE ACTIONS ---
        async sendMessage(content) {
            try {
                await this.$axios.post(
                    `/conversation/${this.conversationID}/message`,
                    {
                        content: content,
                        replied_message_id: this.replyedMessage.message_id
                    },
                    { headers: { Authorization: `Bearer ${this.token}` } }
                );
                this.fetchMessages();
            } catch (err) {
                console.error("Errore durante l'invio del messaggio:", err);
            }
            this.replyedMessage = { replied_message_id: null, content: '' };
        },
        async sendImage(file) {
            const formData = new FormData();
            formData.append('photo', file);
            if (this.replyedMessage.message_id) {
                formData.append('reply_message_id', this.replyedMessage.message_id);
            }

            try {
                await this.$axios.post(`/conversation/${this.conversationID}/photo`, formData, {
                    headers: {
                        Authorization: `Bearer ${this.token}`,
                        "Content-Type": "multipart/form-data",
                    },
                });
            } catch (err) {
                this.error = 'Errore durante il caricamento della foto.';
            }
            this.photoPopup = false;
            this.fetchMessages();
        },
        async deleteMessage(messageId) {
            try {
                await this.$axios.delete(`/message/${messageId}`, {
                    headers: { Authorization: `Bearer ${this.token}` }
                });
                this.messages = this.messages.filter((m) => m.message_id !== messageId);
            } catch (err) {
                console.error("Errore durante l'eliminazione del messaggio:", err);
            }
        },
        replyMessage(message) {
            this.replyedMessage = message;
        },

        // --- FORWARDING ---
        async openForwardModal(messageId) {
            this.forwardMessageId = messageId;
            // We need to fetch conversations that are groups
            try {
                const response = await this.$axios.get("/conversation", {
                    headers: { Authorization: `Bearer ${this.token}` },
                });
                this.groups = response.data.filter(conv => conv.is_group);
            } catch (err) {
                this.error = "Errore nel recupero dei gruppi";
            }
            await this.fetchUsers();
            this.forwardPopup = true;
        },
        async handleForward({ groups, users }) {
            this.forwardPopup = false;
            const forwardPromises = [
                ...groups.map(conv => this.$axios.post(`/message/${this.forwardMessageId}/forwarded`, { receiver_id: conv.conversation_id }, { headers: { Authorization: `Bearer ${this.token}` } })),
                ...users.map(user => this.forwardToUser(user))
            ];

            try {
                await Promise.all(forwardPromises);
            } catch (error) {
                console.error("Uno o più inoltri sono falliti:", error);
            } finally {
                this.forwardMessageId = null;
                this.fetchMessages(); // To update the original message status
            }
        },
        async forwardToUser(user) {
            let conversationId;
            try {
                const response = await this.$axios.post("/conversation", { name: user.name }, { headers: { Authorization: `Bearer ${this.token}` } });
                conversationId = response.data.conversation_id;
            } catch (err) {
                if (err.response?.status === 409 && err.response?.data?.conversation_id) {
                    conversationId = err.response.data.conversation_id;
                } else {
                    throw err;
                }
            }
            return this.$axios.post(`/message/${this.forwardMessageId}/forwarded`, { receiver_id: conversationId }, { headers: { Authorization: `Bearer ${this.token}` } });
        },

        // --- REACTIONS ---
        openReactionModal(messageId) {
            this.reactMessageID = messageId;
            this.reactionPopup = true;
        },
        async handleReaction(reaction) {
            this.reactionPopup = false;
            try {
                await this.$axios.post(`/message/${this.reactMessageID}/reaction`,
                    { content: reaction },
                    { headers: { Authorization: `Bearer ${this.token}` } }
                );
                this.fetchMessages();
            } catch (err) {
                console.error("Errore durante l'invio della reazione:", err);
            }
        },
        async deletereaction({ messageId, reactionId }) {
            try {
                await this.$axios.delete(`message/${messageId}/reaction/${reactionId}`, {
                    headers: { Authorization: `Bearer ${this.token}` }
                });
                this.fetchMessages();
            } catch (err) {
                console.error("Errore durante l'eliminazione della reazione:", err);
            }
        },

        // --- GROUP ACTIONS ---
        async changeGroupName(newName) {
            try {
                await this.$axios.put(`/group/${this.conversationID}/name`, { name: newName }, {
                    headers: { Authorization: `Bearer ${this.token}` }
                });
                this.getConversation();
            } catch (err) {
                if (err.response && err.response.status === 409) {
                    this.error = "Nome utente già in uso. Scegli un altro nome.";
                } else {
                    console.error("Errore durante l'aggiornamento del nome del gruppo:", err);
                    this.error = "Errore durante il cambio nome. Riprova.";
                }
            }
        },
        async changeGroupPhoto(file) {
            const formData = new FormData();
            formData.append('photo', file);
            try {
                await this.$axios.post(`/group/${this.conversationID}/photo`, formData, {
                    headers: {
                        Authorization: `Bearer ${this.token}`,
                        "Content-Type": "multipart/form-data",
                    },
                });
                this.fetchPhoto();
            } catch (err) {
                this.error = 'Errore durante il caricamento della foto del gruppo.';
            }
            this.groupPhotoPopup = false;
        },
        async leaveGroup() {
            try {
                await this.$axios.delete(`/group/${this.conversationID}`, {
                    headers: { Authorization: `Bearer ${this.token}` }
                });
                alert('Sei uscito dal gruppo');
                this.$router.push('/'); // Use router for navigation
            } catch (err) {
                console.error("Errore durante l'uscita dal gruppo:", err);
            }
        },

        // --- UTILS ---
        scrollToBottom() {
            const container = document.getElementById('messagesContainer');
            if (container) {
                container.scrollTop = container.scrollHeight;
            }
        },
        
        // --- POLLING ---
        async getNewMessages() {
            // This is a simple polling mechanism. For a real-time app, consider WebSockets.
            try {
                const res = await this.$axios.get('/me/newmessage', {
                    headers: { Authorization: `Bearer ${this.token}` }
                });
                if (res.data && res.data.length > 0) { // Check if there are new messages
                    await this.fetchMessages();
                }
            } catch (err) {
                // Don't log error for polling to avoid console spam
            }
        }
    },
    mounted() {
        this.getConversation();
        this.fetchMessages();
        this.intervalID = setInterval(this.getNewMessages, 3000); // Poll every 3 seconds
    },
    beforeUnmount() {
        clearInterval(this.intervalID);
    }
}
</script>

<style scoped>
.conversation {
    width: 100%;
    margin: 0;
    font-family: Arial, sans-serif;
}

.container {
    border: 1px solid #888;
    height: 400px;
    margin: 10px 0px;
    padding: 0px;
    background-color: #f4f6f8;
    overflow-y: auto;
    border-radius: 23px;
}
</style>
