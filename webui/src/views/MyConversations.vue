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
        }
    },
    mounted() {
        this.getConversations()
    },
    methods: {
        async getConversations() {
            try {
                const res = await this.$axios.get('/conversation', {
                    headers: { Authorization: `Bearer ${this.token}` },
                })
                const newData = res.data
                if (JSON.stringify(newData.value) === JSON.stringify(this.conversations)) return
                this.conversations = newData

                for (const conv of this.conversations) {
                    try {
                        const res = await this.$axios.get(
                            `/conversation/${conv.conversation_id}/lastmessage`,
                            {
                                headers: { Authorization: `Bearer ${this.token}` },
                            }
                        )
                        conv.message = res.data
                    } catch (err) {
                        this.error = `Errore nel caricamento dell'ultimo messaggio per la conversazione ${conv.conversation_id}`
                        console.error(err)
                    }
                    conv.members = await this.getMembers(conv)
                    conv.image = await this.fetchPhoto(conv)
                }
            } catch (err) {
                this.error = 'Errore nel caricamento delle conversazioni'
                console.error(err)
                return
            }
            this.conversations.sort((a, b) => {
                return new Date(b.message.timestamp) - new Date(a.message.timestamp)
            })
        },
        async getMembers(conversation) {
            try {
                const res = await this.$axios.get(`/group/${conversation.conversation_id}/members`, {
                    headers: { Authorization: `Bearer ${this.token}` },
                })
                conversation.members = res.data
            } catch (err) {
                this.error = 'Errore nel caricamento dei membri della conversazione'
                console.error(err)
            }
            return conversation.members
        },
        async fetchPhoto(conversation) {

            if (conversation.is_group) {
                try {
                    const res = await this.$axios.get(`/group/${conversation.conversation_id}/photo`, {
                        headers: { Authorization: `Bearer ${this.token}` },
                        responseType: 'blob'
                    })
                    conversation.image = URL.createObjectURL(res.data);
                } catch (err) {
                    console.error('Errore nel caricamento della foto del gruppo:', err)
                }
            } else {
                try {
                    console.log("Utenti disponibili:", this.users)

                    const user = conversation.members.find(user => user.name === conversation.name);
                    const res = await this.$axios.get(`/user/${user.user_id}/photo`, {
                        headers: { Authorization: `Bearer ${this.token}` },
                        responseType: 'blob'
                    })
                    conversation.image = URL.createObjectURL(res.data);
                } catch (err) {
                    console.error('Errore nel caricamento della foto dell\'utente:', err)
                }
            }
            return conversation.image
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