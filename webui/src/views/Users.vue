<template>
  <div class="container mt-4">
    <h1>Lista Utenti</h1>

    <div v-if="loading">Caricamento utenti...</div>
    <div v-if="error" class="alert alert-danger">{{ error }}</div>

    <div v-if="users.length">
      <div v-for="user in users" :key="user.id" class="mb-2">
        <button
          class="user"
          @click="startConversation(user.name)"
        >
          {{ user.username }}
        </button>
      </div>
    </div>

    <div v-else-if="!loading">
      Nessun utente trovato.
    </div>
  </div>
</template>

<script>
export default {
  name: "Users",
  data() {
    return {
      users: [],
      loading: false,
      error: null,
      token: localStorage.getItem("token"),
      name: localStorage.getItem("name"),
      username: localStorage.getItem("username"),
      user_id: localStorage.getItem("user_id"),
    };
  },
  methods: {
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
      for (const user of this.users) {
        if (user.name === this.name) {
          this.users.splice(this.users.indexOf(user), 1);
        }
      }
    },

    async startConversation(userName) {
      try {
        const response = await this.$axios.post("/conversation",
          { name: userName },
          {
            headers: {
              Authorization: `Bearer ${this.token}`,
            },
          }
        );

        const conversation = response.data;
        this.$router.push(`/conversation/${conversation.conversation_id}`);
      } catch (err) {
        if (err.response?.status === 409 && err.response?.data?.conversation_id) {
          // Se la conversazione esiste già
          this.$router.push(`/conversation/${err.response.data.conversation_id}`);
        } else {
          console.error("Errore nella creazione della conversazione", err);
          this.error = "Impossibile avviare la conversazione.";
        }
      }
    },
  },
  mounted() {
    this.fetchUsers();
  },
};
</script>

<style>
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
.user:hover {
    background-color: #c3e6cb;
    cursor: pointer;
}
</style>