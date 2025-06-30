<template>
    <div class="backpopup" v-if="show">
        <div class="popup">
            <div class="popup-header">
                <h3>Aggiungi una reazione</h3>
                <button id="exit-button" @click="$emit('close')">
                    <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x" /></svg>
                </button>
            </div>
            <div class="popup-content">
                <ul id="reaction-list">
                    <li v-for="reaction in reactions" :key="reaction">
                        <button @click="$emit('react', reaction)">
                            {{ reaction }}
                        </button>
                    </li>
                </ul>
            </div>
        </div>
    </div>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue';

defineProps({
    show: { type: Boolean, required: true }
});

defineEmits(['close', 'react']);

const reactions = ['👍', '❤️', '😂', '😮', '😢', '😡'];
</script>

<style scoped>
.backpopup { position: fixed; top: 0; left: 0; width: 100%; height: 100%; background-color: rgba(0, 0, 0, 0.6); display: flex; justify-content: center; align-items: center; z-index: 1000; }
.popup { background-color: #f4f6f8; border-radius: 23px; padding: 20px; width: 80%; max-width: 600px; }
.popup-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; }
#reaction-list { font-size: 40px; display: flex; flex-wrap: wrap; justify-content: center; gap: 10px; list-style: none; padding: 0; margin: 0; }
#reaction-list button { background: none; border: none; cursor: pointer; font-size: inherit; }
</style>