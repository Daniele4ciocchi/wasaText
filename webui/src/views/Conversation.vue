<template>
    <div class="conversation">
        <div class="header">
            <h1 id="convName">{{ conversation.name }}</h1>
            <button class="leave" v-if="conversation.is_group == 1" @click="leaveGroup()">
                <svg class="feather">
                    <use href="/feather-sprite-v4.29.0.svg#x" />
                </svg>
                Esci dal gruppo
            </button>
        </div>

        <div class="messages" ref="messagesContainer">
            <div v-for="(message, index) in messages" :key="index"
                :class="['message', message.sender === currentUser ? 'user' : 'receiver']">

                <!-- Mostra il contenuto del messaggio a cui si sta rispondendo, se presente -->
                <div v-if="message.replied_message_id != 0">
                    <p class="replied-message">
                        {{ getMessageById(message.replied_message_id).content }}
                    </p>
                </div>

                <p class="sender">{{ message.sender }}</p>
                <p class="content">{{ message.content }}</p>
                <p id="timestamp">{{ new Date(message.timestamp).toLocaleString() }}</p>
                <button id="replies" @click="replyMessage(message.message_id, message.content)">
                    <svg class="feather">
                        <use href="/feather-sprite-v4.29.0.svg#corner-down-right" />
                    </svg>
                </button>
            </div>
        </div>


        <div>
            <label for="replies" v-if="replyID">
                risposta a : {{ replyContent, replyID }}
            </label>
        </div>
        <div class="input-area">

            <input v-model="newMessage" type="text" placeholder="Scrivi un messaggio..." @keyup.enter="sendMessage" />
            <button id ="send-button"@click="sendMessage">
                <svg class="feather">
                    <use href="/feather-sprite-v4.29.0.svg#send" />
                </svg>
            </button>
        </div>

        <div class="user-list" v-if="conversation.is_group">
            <h2>Utenti nella conversazione</h2>
            <ul>
                <li v-for="user in users" :key="user.name">
                    {{ user.name }}
                </li>
            </ul>

        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'


const route = useRoute()
const token = localStorage.getItem('token')
const currentUser = localStorage.getItem('name')

const conversationID = ref(route.params.conversationID)
const conversation = ref({})
const users = ref([])
const messages = ref([])
const newMessage = ref('')
const replyID = ref(null)
const replyContent = ref('')
const lastMessageId = ref(0)
let messagesContainer = document.getElementById('messagesContainer')
let intervalID



const getConversation = async () => {
    try {
        const res = await axios.get(`http://100.87.168.104:3000/conversation/${conversationID.value}`, {
            headers: { Authorization: `Bearer ${token}` }
        })
        conversation.value = res.data
    } catch (err) {
        console.error('Errore nel caricamento della conversazione:', err)
    }
    if (conversation.value.is_group === true) fetchUsers()
}

const scrollToBottom = async () => {
    await nextTick()
    if (messagesContainer.value) {
        requestAnimationFrame(() => {
            messagesContainer.scrollTop = messagesContainer.scrollHeight
            messagesContainer.scrollIntoView({ behavior: 'smooth' })
        })
    }

}


const fetchMessages = async () => {
    try {
        const res = await axios.get(`http://100.87.168.104:3000/conversation/${conversationID.value}/message`, {
            headers: { Authorization: `Bearer ${token}` }
        })
        messages.value = res.data
    } catch (err) {
        console.error('Errore nel caricamento dei messaggi:', err)
    }

    if (messages.value.at(-1).message_id != lastMessageId.value.message_id) {
        lastMessageId.value.message_id = messages.value.at(-1).message_id
        scrollToBottom()
        alert('Nuovo messaggio ricevuto')
    }
    scrollToBottom()

}
const fetchUsers = async () => {
    try {
        const res = await axios.get(`http://100.87.168.104:3000/group/${conversationID.value}`, {
            headers: { Authorization: `Bearer ${token}` }
        })
        users.value = res.data.members

    } catch (err) {
        console.error('Errore nel caricamento degli utenti:', err)
    }

}

const sendMessage = async () => {
    if (newMessage.value.trim() === '') return


    const messageToSend = {
        replied_message_id: replyID.value,
        content: newMessage.value,
        sender: currentUser,
        timestamp: new Date().toISOString()
    }

    try {
        await axios.post(
            `http://100.87.168.104:3000/conversation/${conversationID.value}/message`,
            {
                content: messageToSend.content,
                replied_message_id: messageToSend.replied_message_id,
            },
            { headers: { Authorization: `Bearer ${token}` } }
        )

        messages.value.push(messageToSend)
        newMessage.value = ''
    } catch (err) {
        console.error('Errore durante l\'invio del messaggio:', err)
    }
    replyID.value = null
    replyContent.value = ''
}
const replyMessage = (messageId, messageContent) => {
    replyID.value = messageId
    replyContent.value = messageContent
}

const getMessageById = (id) => {
    return messages.value.find(message => message.message_id === id) || {};
}

const leaveGroup = async () => {
    try {
        await axios.delete(`http://100.87.168.104:3000/group/${conversationID.value}`, {
            headers: { Authorization: `Bearer ${token}` }
        })
        alert('Sei uscito dal gruppo')
        window.location.href = '/'
    } catch (err) {
        console.error('Errore durante l\'uscita dal gruppo:', err)
    }
}

onMounted(() => {
    getConversation()
    fetchMessages()
    scrollToBottom()
    intervalID = setInterval(fetchMessages, 5000)
})

onUnmounted(() => {
    clearInterval(intervalID)
})

watch(() => route.params.conversationID, (newId) => {
    conversationID.value = newId
    getConversation()
    fetchMessages()
})
</script>

<style scoped>
#convName {
    font-size: 2em;
    text-align: left;
    margin: 15px 0px;
}

.conversation {
    width: 100%;
    margin: 0 auto;
    font-family: Arial, sans-serif;
}

.header {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.leave {
    background-color: #fa716c;
    color: black;
    border-radius: 10px;
    cursor: pointer;
    border: #888 1px solid;
    padding: 5px 10px;
}

.leave:hover {
    background-color: #b02a2a;
}

#send-button {
    border-radius: 50px;
}

.messages {
    border: 1px solid #888;
    padding: 10px;
    height: 400px;
    overflow-y: auto;
    margin-bottom: 10px;
    background-color: #f4f6f8;
    display: flex;
    flex-direction: column;
    align-items: flex-end;

    border-radius: 23px;
}

.message {
    margin-bottom: 10px;
    padding: 5px 15px;
    border-radius: 15px;
    width: auto;
    border: #888 1px solid;
    display: block;
}

.message.user {
    background-color: #d1e7dd;
    text-align: right;
    align-self: flex-end;
}

.message.receiver {
    background-color: #ffffff;
    text-align: left;
    align-self: flex-start;
}

.replied-message {
    background-color: #e2e3e5;
    border-radius: 10px;
    border: #888 1px solid;
    padding: 5px;
    margin: 5px 0;
}

.sender {
    font-weight: bold;
    margin: 0;
}

.content {
    margin: 5px 0;
}

.input-area {
    display: flex;
    gap: 10px;
}

.user-list {
    margin-top: 20px;
    padding: 10px;
    border: 1px solid #888;
    border-radius: 15px;
    background-color: #f4f6f8;
}

input {
    flex: 1;
    padding: 5px;
    border: 1px solid #888;
    border-radius: 15px;
    background-color: #d1e7dd;
    color: black;
}

button {
    padding: 5px 10px;
    border: none;
    background-color: #d7f2ba;
    color: black;
    border-radius: 10px;
    cursor: pointer;
    border: #888 1px solid;
}

button:hover {
    background-color: #bde4a8;
}

#replies {
    align-items: left;
    width: 100%;
    padding: 0px;
    margin: 0px;
}

#timestamp {
    font-size: 0.8em;
    color: #888;
}
</style>
