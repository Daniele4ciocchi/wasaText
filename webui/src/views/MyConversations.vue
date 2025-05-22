<template>
    <ul class="nav flex-column">
        <h2>Conversazioni</h2>
        <li class="nav-item" v-for="conv in conversations" :key="conv.conversation_id">
            <RouterLink :to="'/conversation/' + conv.conversation_id" class="nav-link">
                <div class="conversation">
                    <div id="conversation-name">
                        <svg class="feather">
                            <use href="/feather-sprite-v4.29.0.svg#book" />
                        </svg>
                        {{ conv.name }}
                    </div>

                    <div id="conversation-message" v-if="conv.message">
                        <svg class="feather">
                            <use href="/feather-sprite-v4.29.0.svg#message-square" />
                        </svg>
                        {{ conv.message.content.length > 20 ? conv.message.content.slice(0, 20) + '...' : conv.message.content }}
                        <span id="timestamp">
                            {{ new Date(conv.message.timestamp).toLocaleString() }}
                        </span>
                    </div>
                </div>
            </RouterLink>
        </li>
    </ul>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'


// Token: prendi da props, store o composable
const token = localStorage.getItem('token') // o come preferisci


const conversations = ref([])
const error = ref(null)



const getConversations = async () => {
    try {
        const res = await axios.get('/conversation', {
            headers: { Authorization: `Bearer ${token}` },
        })
        const newData = res.data
        if (JSON.stringify(newData.value) === JSON.stringify(conversations.value)) return
        conversations.value = newData

        for (const conv of conversations.value) {
            try {
                const res = await axios.get(
                    `/conversation/${conv.conversation_id}/lastmessage`,
                    {
                        headers: { Authorization: `Bearer ${token}` },
                    }
                )
                conv.message = res.data
            } catch (err) {
                error.value = `Errore nel caricamento dell'ultimo messaggio per la conversazione ${conv.conversation_id}`
                console.error(err)
            }
        }
    } catch (err) {
        error.value = 'Errore nel caricamento delle conversazioni'
        console.error(err)
        return
    }


}


onMounted(() => {
    getConversations()
    
})
</script>

<style>
.conversation {
    color: black;

}
#conversation-message {
    font-size: 10px;
}

#timestamp {
    font-size: 8px;
    color: gray;
    margin-left: 5px;
    text-align: right;
}
</style>