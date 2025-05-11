<template>
  <div class="container mt-4">
    <h1>Lista Utenti</h1>

    <div v-if="loading">Caricamento utenti...</div>
    <div v-if="error" class="alert alert-danger">{{ error }}</div>

    <div v-if="users.length">
      <div v-for="user in users" :key="user.id" class="mb-2">
        <button
          class="btn btn-outline-primary w-100 text-start"
          @click="startConversation(user.name)"
        >
          {{ user.name }}
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
  data() {
    return {
      users: [],
      loading: false,
      error: null,
      token: localStorage.getItem("token"),
    };
  },
  methods: {
    async fetchUsers() {
      this.loading = true;
      try {
        const response = await this.$axios.get("http://localhost:3000/user", {
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
    },

    async startConversation(userName) {
      try {
        const response = await this.$axios.post(
          "http://localhost:3000/conversation",
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
