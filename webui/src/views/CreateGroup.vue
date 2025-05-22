<template>
    <div class="create-group">
        <h1>Crea un Nuovo Gruppo</h1>

        <div class="form-group">
            <input type="text" id="group-name" required v-model="groupName"
                placeholder="Inserisci il nome del gruppo" />
        </div>

        <h2>Lista Utenti</h2>
        <div class="users-list">
            
            <ul>
                <li class="user" v-for="user in users" :key="user.user_id">
                    <label>
                        <input type="checkbox" :value="user.user_id" v-model="selectedUsers" />
                        {{ user.name }}
                    </label>
                </li>
            </ul>
        </div>

        <button type="submit" @click="createGroup">Crea Gruppo</button>
        <p id="message" v-if="message" class="error">{{ message }}</p>
    </div>
</template>

<script>
export default {
    data() {
        return {
            name: localStorage.getItem("name"),
            groupName: '',
            users: [],
            selectedUsers: [],
            token: localStorage.getItem("token"),
            message: '',
        };
    },
    methods: {
        async createGroup() {
            if (!this.groupName.trim()) {
                alert("Il nome del gruppo è obbligatorio.");
                return;
            }

            try {
                const response = await this.$axios.post("/group", {
                    name: this.groupName,
                    members: this.selectedUsers,
                }, {
                    headers: {
                        Authorization: `Bearer ${this.token}`,
                    },
                });

                // Se siamo qui, la richiesta ha avuto successo
                this.message = "✅ Gruppo creato con successo!";
                console.log("Status:", response.status); // Es: 201
            } catch (error) {
                // Qui entri solo se lo status è 4xx o 5xx
                if (error.response) {
                    this.message = `❌ Errore ${error.response.status}: ${error.response.data.message || "Errore nella creazione del gruppo."}`;
                } else {
                    this.message = "❌ Errore di rete o server non raggiungibile.";
                }
            }
        }
        ,
        async fetchUsers() {
            this.loading = true;
            try {
                const response = await this.$axios.get("/user", {
                    headers: {
                        Authorization: `Bearer ${this.token}`,
                    },
                });
                this.users = response.data;
            } catch (err) {
                this.error = "Errore nel recupero degli utenti";
            } finally {
                this.loading = false;
            }
            for (let i = 0; i < this.users.length; i++) {
                if (this.users[i].name === this.name) {
                    this.users.splice(i, 1);
                    break;
                }
            }
        },

    },
    mounted() {
        this.fetchUsers();
    },
};

</script>

<style scoped>
.create-group {
    border: 1px solid #888;
    padding: 10px;
    height: 600px;
    overflow-y: auto;
    margin: 10px 0px;
    background-color: #f4f6f8;
    display: flex;
    flex-direction: column;
    border-radius: 23px;

}

.form-group {
    margin-bottom: 20px;
}

.users-list {
    border: 1px solid #888;
    padding: 10px;
    height: 400px;
    overflow-y: auto;
    margin-bottom: 10px;
    background-color: #f4f6f8;
    display: flex;
    flex-direction: column;

    border-radius: 23px;
}


.user {
    text-align: center;
    width: 60%;
    display: flex;
    margin-bottom: 10px;
    border: 1px solid #888;
    border-radius: 10px;
    padding: 10px;
    background-color: #d1e7dd;
}

button {
    margin: 0px 10px;
    padding: 10px 20px;
    background-color: #bde4a8;
    color: black;
    border: 1px solid #888;
    border-radius: 20px;
    cursor: pointer;
}

button:hover {
    background-color: #9cbd8a;
}

#group-name {
    width: 60%;
    padding: 10px;
    border: 1px solid #888;
    border-radius: 20px;
    background-color: #d1e7dd;
    color: black;
}

#message {
    align-self: center;
    margin-top: 10px;
    color: green;
}
</style>