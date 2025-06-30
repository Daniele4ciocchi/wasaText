<template>
    <div class="messages" id="messagesContainer">
        <MessageItem 
            v-for="message in messages" 
            :key="message.message_id"
            :message="message"
            :current-user-name="currentUserName"
            :current-user-id="currentUserId"
            :replied-message-text="getRepliedMessageText(message)"
            @reply="$emit('reply', $event)"
            @delete="$emit('delete', $event)"
            @forward="$emit('forward', $event)"
            @react="$emit('react', $event)"
            @delete-reaction="$emit('delete-reaction', $event)"
        />
    </div>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue';
import MessageItem from './MessageItem.vue';

const props = defineProps({
    messages: { type: Array, required: true },
    currentUserName: { type: String, required: true },
    currentUserId: { type: [String, Number], required: true },
});

defineEmits(['reply', 'delete', 'forward', 'react', 'delete-reaction']);

const getMessageById = (id) => {
    return props.messages.find((m) => m.message_id === id) || {};
};

const getRepliedMessageText = (message) => {
    if (!message.replied_message_id) return '';
    const repliedMsg = getMessageById(message.replied_message_id);
    if (!repliedMsg) return '';
    return `${repliedMsg.sender}, ${repliedMsg.photo ? 'photo' : repliedMsg.content}`;
};
</script>

<style scoped>
.messages {
    display: flex;
    flex-direction: column;
    gap: 10px;
    padding: 10px;
    overflow-y: auto;
    height: 100%;
}
</style>