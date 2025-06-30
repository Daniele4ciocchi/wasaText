<template>
    <div :class="['message', message.sender === currentUserName ? 'user' : 'receiver']">
        <!-- Messaggio Risposto -->
        <div v-if="message.replied_message_id && repliedMessageText" class="replied-message">
            <p>Risposta a {{ repliedMessageText }}</p>
        </div>

        <!-- Tag Forwarded -->
        <div v-if="message.forwarded" class="forwarded-tag">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#corner-down-right" /></svg>
            <span>Inoltrato</span>
        </div>

        <!-- Contenuto del Messaggio -->
        <div v-if="!message.photo" class="normal-message">
            <p class="sender">{{ message.sender }}</p>
            <p class="content">{{ message.content }}</p>
            <p class="timestamp">{{ new Date(message.timestamp).toLocaleString() }}</p>
        </div>
        <div v-else class="photo-message">
            <p class="sender">{{ message.sender }}</p>
            <img :src="message.photoUrl" alt="Foto del messaggio" style="max-width: 100%; border-radius: 10px;" />
            <p class="timestamp">{{ new Date(message.timestamp).toLocaleString() }}</p>
        </div>

        <!-- Footer del Messaggio (Azioni) -->
        <div class="message-footer">
            <div v-if="message.sender === currentUserName" class="view">
                <svg v-if="message.status === 0" class="feather"><use href="/feather-sprite-v4.29.0.svg#check" /></svg>
                <svg v-if="message.status === 1" class="feather"><use href="/feather-sprite-v4.29.0.svg#check-circle" /></svg>
                <svg v-if="message.status === 2" class="feather"><use href="/feather-sprite-v4.29.0.svg#check-square" /></svg>
            </div>
            <button class="reply-btn" @click="$emit('reply', message)">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#repeat" /></svg>
            </button>
            <button v-if="message.sender === currentUserName" class="delete-btn" @click="$emit('delete', message.message_id)">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#trash" /></svg>
            </button>
            <button class="forward-btn" @click="$emit('forward', message.message_id)">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#corner-down-right" /></svg>
            </button>
            <button class="react-btn" @click="$emit('react', message.message_id)">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#plus-circle" /></svg>
            </button>
        </div>

        <!-- Reazioni -->
        <div class="message-reactions" v-if="message.reactions && message.reactions.length > 0">
            <span v-for="reaction in message.reactions" :key="reaction.reaction_id" class="reaction">
                {{ reaction.content }} {{ reaction.users.length }}
                <button v-if="reaction.users.some(user => user.user_id == currentUserId)" 
                        @click="$emit('delete-reaction', { messageId: message.message_id, reactionId: reaction.users.find(user => user.user_id == currentUserId).reaction_id })"
                        class="reaction-delete-button">
                    <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x" /></svg>
                </button>
            </span>
        </div>
    </div>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue';

defineProps({
    message: { type: Object, required: true },
    currentUserName: { type: String, required: true },
    currentUserId: { type: [String, Number], required: true },
    repliedMessageText: { type: String, default: '' }
});

defineEmits(['reply', 'delete', 'forward', 'react', 'delete-reaction']);
</script>

<style scoped>
.message {
    margin: 0px 15px 10px;
    padding: 5px 15px;
    border-radius: 15px;
    width: auto;
    border: #888 1px solid;
    display: block;
}
.message.user { background-color: #d1e7dd; text-align: right; align-self: flex-end; }
.message.receiver { background-color: #ffffff; text-align: left; align-self: flex-start; }

.replied-message { background-color: #e2e3e5; border-radius: 10px; border: #888 1px solid; padding: 5px; margin: 5px 0; }
.sender { font-weight: bold; margin: 0; }
.content { margin: 5px 0; border: #888 1px solid; border-radius: 10px; padding: 5px; }
.timestamp { font-size: 0.8em; color: #888; }

.forwarded-tag {
    font-size: 0.8em;
    color: #666;
    display: flex;
    align-items: center;
    gap: 5px;
}
.forwarded-tag .feather { width: 12px; height: 12px; }

.photo-message { 
    max-width: 250px; /* Limita la larghezza massima del box */
    max-height: 250px; /* Imposta un'altezza massima anche per il contenitore */
    display: flex;
    flex-direction: column;
    align-items: center;
    border-radius: 10px;
    overflow: hidden; /* Nasconde eventuali overflow */
}
.photo-message img { 
    width: 100%; /* Fa in modo che l'immagine riempia la larghezza del contenitore */
    height: 100%; /* Fa in modo che l'immagine riempia l'altezza del contenitore */
    object-fit: contain; /* Assicura che l'intera immagine sia visibile e mantenga le proporzioni */
    border-radius: 10px; 
}

.message-footer { display: flex; flex-direction: row; justify-content: flex-end; gap: 10px; margin-top: 5px; }
.message-footer button { 
    background-color: white; /* Sfondo leggero */
    border: #888 1px solid; /* Bordo sottile */
    border-radius: 10px; /* Angoli arrotondati */
    cursor: pointer; 
    padding: 5px 8px; /* Più padding */
    transition: background-color 0.2s ease; /* Transizione per hover */
}
.message-footer button:hover {
    background-color: #d0d0d0; /* Sfondo più scuro al hover */
}
.message-footer .feather { width: 16px; height: 16px; color: #555; }

.view { width: 15px; height: 15px; color: #6b6b6b; align-self: center; margin-right: auto; }

.message-reactions { display: flex; justify-content: flex-start; margin-top: 5px; gap: 5px; flex-wrap: wrap; }
.reaction { border: #888 1px solid; padding: 4px 8px; border-radius: 12px; background: #fff; font-size: 1em; display: flex; align-items: center; gap: 4px; }
.reaction-delete-button { width: 20px; cursor: pointer; padding: 0; margin: 0; background: none; border: none; }
.reaction-delete-button .feather { height: 10px; width: 10px; padding: 0; margin: 0; }
</style>