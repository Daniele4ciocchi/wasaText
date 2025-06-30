<template>
    <div class="header">
        <div class="header-left">
            <img :src="conversationImage" alt="Foto profilo" style="width: 40px; height:40px; border-radius: 10px;" />
            <button v-if="conversation.is_group" @click="$emit('change-photo')">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#edit" /></svg>
            </button>
        </div>

        <div class="header-center">
            <h1 id="convName">{{ conversation.name }}</h1>
            <button v-if="!isEditingName && conversation.is_group" @click="isEditingName = true">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#edit-3" /></svg>
            </button>
            <div v-if="isEditingName">
                <input v-model="groupName" type="text" placeholder="nuovo nome..." @keyup.enter="updateGroupName" />
                <button @click="updateGroupName">invia</button>
            </div>
        </div>

        <div class="header-right">
            <button class="leave" v-if="conversation.is_group" @click="$emit('leave-group')">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x" /></svg>
                Esci dal gruppo
            </button>
        </div>
    </div>
</template>

<script setup>
import { ref, defineProps, defineEmits, watch } from 'vue';

const props = defineProps({
    conversation: { type: Object, required: true },
    conversationImage: { type: String, default: '' }
});

const emit = defineEmits(['change-name', 'leave-group', 'change-photo']);

const isEditingName = ref(false);
const groupName = ref('');

watch(() => props.conversation, (newVal) => {
    groupName.value = newVal.name;
}, { immediate: true });

const updateGroupName = () => {
    if (groupName.value.trim()) {
        emit('change-name', groupName.value);
        isEditingName.value = false;
    }
};
</script>

<style scoped>
.header { display: flex; justify-content: space-between; align-items: center; margin-top: 10px; }
.header-left { display: flex; justify-content: left; gap: 10px; }
.header-center { display: flex; justify-content: center; align-items: center; gap: 10px; height: 40px; }
#convName { font-size: 2em; text-align: left; }
.leave { background-color: #fa716c; color: black; border-radius: 10px; cursor: pointer; border: #888 1px solid; padding: 5px 10px; }
.leave:hover { background-color: #b02a2a; }
input { padding: 5px; border: 1px solid #888; border-radius: 15px; background-color: #d1e7dd; color: black; }
</style>