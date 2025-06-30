<template>
    <ul class="nav flex-column">
        <h2>Conversazioni</h2>
        <li class="nav-item" v-for="conv in conversations" :key="conv.conversation_id">
            <RouterLink :to="'/conversation/' + conv.conversation_id" class="nav-link">
                <div class="conversation-list">
                    <div id="conversation-photo">
                        <img :src="conv.image" alt="Photo" class="conv-photo" />
                    </div>
                    <div id="conversation-info">
                        <div id="conversation-name">
                            {{ conv.name }}
                        </div>

                        <div id="conversation-message" v-if="conv.message">
                            <svg class="feather">
                                <use href="/feather-sprite-v4.29.0.svg#message-square" />
                            </svg>
                            {{ conv.message.content.length > 20 ? conv.message.content.slice(0, 20) + '...' :
                            conv.message.content }}
                            <span id="timestamp">
                                <svg v-if="conv.message.status === 0" class="feather status-icon"><use href="/feather-sprite-v4.29.0.svg#check" /></svg>
                                <svg v-if="conv.message.status === 1" class="feather status-icon"><use href="/feather-sprite-v4.29.0.svg#check-circle" /></svg>
                                <svg v-if="conv.message.status === 2" class="feather status-icon"><use href="/feather-sprite-v4.29.0.svg#check-square" /></svg>
                                {{ new Date(conv.message.timestamp).toLocaleString() }}
                            </span>
                        </div>
                    </div>

                </div>
            </RouterLink>
        </li>
    </ul>
</template>

<script>
export default {
    name: 'MyConversations',
    data() {
        return {
            token: localStorage.getItem('token'),
            name: localStorage.getItem('name'),
            username: localStorage.getItem('username'),
            user_id: localStorage.getItem('user_id'),
            conversations: [],
            error: null,
            pollingInterval: null, // Per memorizzare l'ID dell'intervallo
        }
    },
    mounted() {
        this.getConversations();
        // Avvia il polling ogni 3 secondi, assicurando il contesto corretto
        this.pollingInterval = setInterval(() => {
            this.getConversations();
        }, 3000);
    },
    beforeUnmount() {
        // Pulisce l'intervallo quando il componente viene distrutto
        clearInterval(this.pollingInterval);
    },
    methods: {
        async getConversations() {
            try {
                const res = await this.$axios.get('/conversation', {
                    headers: { Authorization: `Bearer ${this.token}` },
                });
                const newData = res.data;
                
                // Evita aggiornamenti inutili se i dati non sono cambiati
                if (JSON.stringify(newData) === JSON.stringify(this.conversations)) return;

                this.conversations = newData;

                for (const conv of this.conversations) {
                    try {
                        const res = await this.$axios.get(
                            `/conversation/${conv.conversation_id}/lastMessage`,
                            {
                                headers: { Authorization: `Bearer ${this.token}` },
                            }
                        );
                        conv.message = res.data;
                    } catch (err) {
                        this.error = `Errore nel caricamento dell'ultimo messaggio per la conversazione ${conv.conversation_id}`;
                        console.error(err);
                    }
                    conv.members = await this.getMembers(conv);
                    conv.image = await this.fetchPhoto(conv);
                }
            } catch (err) {
                this.error = 'Errore nel caricamento delle conversazioni';
                console.error(err);
                return;
            }
            await this.sortConversations();
        },
        async sortConversations() {
            this.conversations.sort((a, b) => {
                // Assicurati che message e timestamp esistano prima di accedere
                const timestampA = a.message ? new Date(a.message.timestamp).getTime() : 0;
                const timestampB = b.message ? new Date(b.message.timestamp).getTime() : 0;
                return timestampB - timestampA;
            });
        },
        async getMembers(conversation) {
            try {
                const res = await this.$axios.get(`/conversation/${conversation.conversation_id}/member`, {
                    headers: { Authorization: `Bearer ${this.token}` },
                });
                conversation.members = res.data;
            } catch (err) {
                this.error = 'Errore nel caricamento dei membri della conversazione';
                console.error(err);
            }
            return conversation.members;
        },
        async fetchPhoto(conversation) {
            if (conversation.is_group) {
                try {
                    const res = await this.$axios.get(`/group/${conversation.conversation_id}/photo`, {
                        headers: { Authorization: `Bearer ${this.token}` },
                        responseType: 'blob'
                    });
                    conversation.image = URL.createObjectURL(res.data);
                } catch (err) {
                    console.error('Errore nel caricamento della foto del gruppo:', err);
                }
            } else {
                try {
                    // Per le conversazioni individuali, la foto è quella dell'altro utente
                    // Assicurati che conversation.members sia già popolato
                    const otherUser = conversation.members.find(member => member.user_id !== this.user_id);
                    if (otherUser) {
                        const res = await this.$axios.get(`/user/${otherUser.user_id}/photo`, {
                            headers: { Authorization: `Bearer ${this.token}` },
                            responseType: 'blob'
                        });
                        conversation.image = URL.createObjectURL(res.data);
                    } else {
                        console.warn('Altro utente non trovato per la conversazione:', conversation.name);
                    }
                } catch (err) {
                    console.error('Errore nel caricamento della foto dell\'utente:', err);
                }
            }
            return conversation.image;
        },
        async checkNewMessagesAndUpdates() {
            try {
                const res = await this.$axios.get('/me/newmessage', {
                    headers: { Authorization: `Bearer ${this.token}` }
                });
                const { new_messages, status_updates } = res.data;

                if (new_messages && new_messages.length > 0) {
                    // Se ci sono nuovi messaggi, ricarica tutte le conversazioni
                    // per aggiornare gli ultimi messaggi e l'ordinamento.
                    await this.getConversations();
                }

                if (status_updates && status_updates.length > 0) {
                    status_updates.forEach(update => {
                        // Trova la conversazione a cui appartiene il messaggio
                        const conv = this.conversations.find(c => c.conversation_id === update.conversation_id);
                        if (conv && conv.message && conv.message.message_id === update.message_id) {
                            // Se l'aggiornamento riguarda l'ultimo messaggio della conversazione,
                            // aggiorna il suo stato.
                            conv.message.status = update.new_status;
                        }
                    });
                }
            } catch (err) {
                // Non loggare errori per il polling per evitare spam nella console
                // a meno che non sia un errore grave (es. 401 Unauthorized)
                if (err.response && err.response.status === 401) {
                    console.error("Sessione scaduta o non autorizzata.", err);
                    // Potresti voler reindirizzare al login qui
                }
            }
        },
    },
}

</script>


<style>
.conversation-list {
    color: black;
    display: flex;
    align-items: center;
    flex-direction: row;
    gap: 15px;
    padding: 5px 0px;
}

.status-icon {
    width: 12px; /* Rendi l'icona più piccola */
    height: 12px;
    margin-right: 3px; /* Spazio tra icona e timestamp */
    vertical-align: middle; /* Allinea verticalmente con il testo */
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

.conv-photo{
    width: 40px;
    height: 40px; 
    border-radius: 10px;
}
</style>