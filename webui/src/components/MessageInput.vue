<template>
    <div>
        <!-- Barra del messaggio risposto -->
        <div v-if="replyedMessage && replyedMessage.content" class="reply-info">
            <p>Risposta a: {{ replyedMessage.photo ? 'photo' : replyedMessage.content }}</p>
            <button @click="$emit('cancel-reply')">Annulla</button>
        </div>

        <!-- Barra di inserimento -->
        <div class="input-area">
            <input v-model="text" type="text" placeholder="Scrivi un messaggio..." @keyup.enter="submitMessage" />
            <button id="photo" @click="$emit('open-photo-popup')">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#image" /></svg>
            </button>
            <button id="send-button" @click="submitMessage" :disabled="!text.trim()">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#send" /></svg>
            </button>
        </div>
    </div>
</template>

<script setup>
import { ref, defineProps, defineEmits, watch } from 'vue';

defineProps({
    replyedMessage: { type: Object, default: null }
});

const emit = defineEmits(['send-message', 'cancel-reply', 'open-photo-popup']);

const text = ref('');

const submitMessage = () => {
    if (text.value.trim()) {
        emit('send-message', text.value);
        text.value = '';
    }
};
</script>

<style scoped>
.reply-info { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; padding: 5px; background-color: #e2e3e5; border-radius: 10px; border: #888 1px solid; }
.reply-info p { margin: 0; }
.input-area { display: flex; gap: 10px; }
input { flex: 1; padding: 5px; border: 1px solid #888; border-radius: 15px; background-color: #d1e7dd; color: black; }
#send-button, #photo { border-radius: 50px; background: none; border: 1px solid #ccc; cursor: pointer; padding: 5px; }
</style>