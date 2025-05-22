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
                }
            } catch (err) {
                this.error = 'Errore nel caricamento delle conversazioni'
                console.error(err)
                return
            }
        },
    },
}

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