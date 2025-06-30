<template>
    <div class="backpopup" v-if="show">
        <div class="popup">
            <div class="popup-header">
                <h3>Invia a</h3>
                <button id="exit-button" @click="$emit('close')">
                    <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x" /></svg>
                </button>
            </div>
            <div class="popup-content">
                <div v-if="groups.length > 0" class="popup-content-conversations">
                    <h4>Gruppi</h4>
                    <ul>
                        <li v-for="group in groups" :key="group.conversation_id" class="users">
                            <label>
                                <input type="checkbox" :value="group" v-model="selectedGroups" />
                                {{ group.name }}
                            </label>
                        </li>
                    </ul>
                </div>
                <div v-if="users.length > 0" class="popup-content-users">
                    <h4>Utenti</h4>
                    <ul>
                        <li v-for="user in users" :key="user.user_id" class="users">
                            <label>
                                <input type="checkbox" :value="user" v-model="selectedUsers" />
                                {{ user.name }}
                            </label>
                        </li>
                    </ul>
                </div>
            </div>
            <div class="popup-footer">
                <button @click="submitForward">Invia</button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, defineProps, defineEmits } from 'vue';

defineProps({
    show: { type: Boolean, required: true },
    groups: { type: Array, default: () => [] },
    users: { type: Array, default: () => [] },
});

const emit = defineEmits(['close', 'forward']);

const selectedGroups = ref([]);
const selectedUsers = ref([]);

const submitForward = () => {
    emit('forward', { groups: selectedGroups.value, users: selectedUsers.value });
    // Resetta le selezioni dopo l'invio
    selectedGroups.value = [];
    selectedUsers.value = [];
};
</script>

<style scoped>
.backpopup {
    position: fixed; /* Usa fixed per coprire l'intera viewport */
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.6);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000; /* Assicura che sia sopra tutto il resto */
}

.popup {
    background-color: #f4f6f8;
    border-radius: 23px;
    padding: 20px;
    width: 80%;
    max-width: 500px;
    height: 80%;
    max-height: 600px;
    display: flex;
    flex-direction: column;
    overflow: hidden; /* Controlla lo scroll dal content */
}

.popup-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
}

.popup-content {
    flex-grow: 1; /* Fa in modo che il contenuto occupi lo spazio disponibile */
    overflow-y: auto; /* Abilita lo scroll solo per il contenuto */
}

.popup-content ul {
    list-style: none;
    padding: 0;
    margin: 0;
}

.users {
    text-align: left; /* Allineamento più leggibile */
    display: flex;
    margin-bottom: 10px;
    border: 1px solid #888;
    border-radius: 10px;
    padding: 10px;
    background-color: #d1e7dd;
}

.popup-footer {
    display: flex;
    justify-content: center;
    margin-top: 10px;
}

button {
    cursor: pointer;
}
</style>