<template>
    <div class="conversation">
      <h1>Conversazione</h1>
  
      <div class="messages">
        <div
          v-for="(message, index) in messages"
          :key="index"
          :class="['message', message.sender_id === currentUserID ? 'user' : 'bot']"
        >
          <p class="sender">Utente #{{ message.sender_id }}</p>
          <p class="content">{{ message.content }}</p>
        </div>
      </div>
  
      <div class="input-area">
        <input
          v-model="newMessage"
          type="text"
          placeholder="Scrivi un messaggio..."
          @keyup.enter="sendMessage"
        />
        <button @click="sendMessage">Invia</button>
      </div>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted } from 'vue'
  import { useRoute } from 'vue-router'
  import axios from 'axios'
  
  const route = useRoute()
  const conversationID = route.params.conversationID
  const token = localStorage.getItem('token')
  const currentUserID = parseInt(localStorage.getItem('user_id')) // ID utente corrente
  
  const messages = ref([])
  const newMessage = ref('')
  
  const fetchMessages = async () => {
    try {
      const response = await axios.get(`http://localhost:3000/conversation/${conversationID}/message`, {
        headers: {
          Authorization: `Bearer ${token}`
        }
      })
      messages.value = response.data || []
    } catch (error) {
      console.error('Errore durante il caricamento dei messaggi:', error)
    }
  }
  
  const sendMessage = async () => {
    if (newMessage.value.trim() === '') return
  
    const messageToSend = {
      content: newMessage.value,
      sender_id: currentUserID
    }
  
    try {
      await axios.post(
        `http://localhost:3000/conversation/${conversationID}/message`,
        { content: messageToSend.content },
        {
          headers: {
            Authorization: `Bearer ${token}`
          }
        }
      )
  
      // Aggiunta locale
      messages.value.push(messageToSend)
      newMessage.value = ''
    } catch (err) {
      console.error("Errore durante l'invio del messaggio:", err)
    }
  }
  
  onMounted(fetchMessages)
  </script>
  
  <style scoped>
  .conversation {
    max-width: 600px;
    margin: 0 auto;
    font-family: Arial, sans-serif;
  }
  
  .messages {
    border: 1px solid #ccc;
    padding: 10px;
    height: 300px;
    overflow-y: auto;
    margin-bottom: 10px;
    background-color: #f9f9f9;
  }
  
  .message {
    margin-bottom: 10px;
    padding: 5px;
    border-radius: 5px;
  }
  
  .message.user {
    background-color: #d1e7dd;
    text-align: right;
  }
  
  .message.bot {
    background-color: #f8d7da;
    text-align: left;
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
  
  input {
    flex: 1;
    padding: 5px;
    border: 1px solid #ccc;
    border-radius: 4px;
  }
  
  button {
    padding: 5px 10px;
    border: none;
    background-color: #007bff;
    color: white;
    border-radius: 4px;
    cursor: pointer;
  }
  
  button:hover {
    background-color: #0056b3;
  }
  </style>
  