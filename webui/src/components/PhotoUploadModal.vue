<template>
    <div class="backpopup" v-if="show">
        <div class="popup">
            <div class="popup-header">
                <h3>{{ title }}</h3>
                <button id="exit-button" @click="$emit('close')">
                    <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x" /></svg>
                </button>
            </div>
            <div class="popup-content">
                <input type="file" @change="handleFileChange" accept="image/*" />
                <img v-if="previewImage" :src="previewImage" alt="Preview" class="preview-image"/>
            </div>
            <div class="popup-footer">
                <button @click="submit" :disabled="!selectedFile">
                    {{ buttonText }}
                </button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, defineProps, defineEmits } from 'vue';

defineProps({
    show: { type: Boolean, required: true },
    title: { type: String, default: 'Carica una foto' },
    buttonText: { type: String, default: 'Invia' }
});

const emit = defineEmits(['close', 'submit']);

const selectedFile = ref(null);
const previewImage = ref(null);

const handleFileChange = (e) => {
    const file = e.target.files[0];
    if (!file) return;

    selectedFile.value = file;
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = (event) => {
        previewImage.value = event.target.result;
    };
};

const submit = () => {
    if (selectedFile.value) {
        emit('submit', selectedFile.value);
        // Reset state after submit
        selectedFile.value = null;
        previewImage.value = null;
    }
};
</script>

<style scoped>
.backpopup { position: fixed; top: 0; left: 0; width: 100%; height: 100%; background-color: rgba(0, 0, 0, 0.6); display: flex; justify-content: center; align-items: center; z-index: 1000; }
.popup { background-color: #f4f6f8; border-radius: 23px; padding: 20px; width: 80%; max-width: 500px; display: flex; flex-direction: column; }
.popup-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; }
.popup-content { text-align: center; }
.preview-image { max-width: 100%; max-height: 300px; margin-top: 15px; border-radius: 10px; }
.popup-footer { display: flex; justify-content: center; margin-top: 10px; }
</style>