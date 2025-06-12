<template>
    <div class="users-list" v-if="addUser">
        <div class="users-list-header">
            <h2>Utenti nella conversazione</h2>
            <button @click="addUser = !addUser">
                <svg class="feather">
                    <use href="/feather-sprite-v4.29.0.svg#x" />
                </svg>
            </button>
        </div>
        <div class="users-list-content">
            <h2>{{ this.users.length != 0? "seleziona gli utenti da aggiungere" : "nessun utente disponibile"  }}</h2>

            <ul>
                <li v-for="user in users" :key="user.name">

                    <label>
                        <input type="checkbox" :value="user" v-model="selectedUsers" />
                        {{ user.name }}
                    </label>
                </li>
            </ul>
            

            <button v-if="this.users.length != 0" @click="addUser = !addUser, addMembers()">
                aggiungi
            </button>

        </div>



    </div>
    <div class="users-list" v-if="!addUser">
        <div class="users-list-header">
            <h2>Utenti nella conversazione</h2>
            <button @click="addUser = !addUser; getUsers()">
                <svg class="feather">
                    <use href="/feather-sprite-v4.29.0.svg#user-plus" />
                </svg>
            </button>
        </div>
        <div class="users-list-content">
            <ul>
                <li v-for="member in members" :key="member.name">
                    {{ member.name }}
                </li>
            </ul>
        </div>
    </div>
</template>

<script>
export default {
    name: 'MembersList',
    props: {
        conversation: {
            type: Object,
            required: true,
        },
        conversationID: {
            type: [String, Number],
            required: true,
        },

    },
    data() {
        return {
            members: [],
            addUser: false,
            users: [],
            selectedUsers: [],

            token: localStorage.getItem("token"),
        };
    },
    mounted() {
        this.getMembers();
    },
    methods: {
        async getMembers() {
            try {
                const response = await this.$axios.get(`/conversation/${this.conversationID}/member`, {
                    headers: {
                        Authorization: `Bearer ${this.token}`,
                    },
                });
                this.members = response.data;
            } catch (error) {
                console.error("Errore nel recupero degli utenti:", error);
            }

        },
        async getUsers() {
            try {
                const response = await this.$axios.get('/user', {
                    headers: {
                        Authorization: `Bearer ${this.token}`,
                    },
                });
                this.users = response.data;
            } catch (error) {
                console.error("Errore nel recupero degli utenti:", error);
            }

            this.users = this.users.filter(user =>
                !this.members.some(member => member.user_id === user.user_id)
            );
        },

        async addMembers() {
            try {
                const response = await this.$axios.post(
                    `/conversation/${this.conversationID}/member`,
                    this.selectedUsers, // dati da inviare
                    {
                        headers: {
                            Authorization: `Bearer ${this.token}`,
                        },
                    }
                );
                // Aggiorna la lista membri dopo l'aggiunta
                await this.getMembers();
                this.selectedUsers = [];
            } catch (error) {
                console.error("Errore nell'aggiunta di un utente:", error);
            }
        }
    },
};
</script>

<style scoped>
.users-list {
    border: 1px solid #888;
    padding: 10px;
    overflow-y: auto;
    margin-top: 10px;
    background-color: #f4f6f8;
    display: flex;
    flex-direction: column;
    border-radius: 23px;
}

.users-list-content {
    margin-top: 20px;
    padding: 10px;
    border: 1px solid #888;
    border-radius: 15px;
    background-color: #f4f6f8;
    justify-content: center;
}

.users-list-content ul {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.users-list-content li {
    padding: 8px 12px;
    border-radius: 8px;
    background-color: #e9ecef;
    border: 1px solid #bdbdbd;
    font-size: 1em;
    color: #333;

}


.users-list-header {
    display: flex;
    justify-content:space-between;
    align-items: left;
    margin-bottom: 10px;
    gap: 10px;
}
</style>
