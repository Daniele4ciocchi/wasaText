<template>
  <div>
    <h1>Contatti</h1>

    <div v-if="loading">Caricamento...</div>
    <div v-if="error" class="error">{{ error }}</div>

    <!-- Input ricerca -->
    <input
      v-model="searchQuery"
      type="text"
      placeholder="Cerca utente..."
      class="form-control mb-3"
    />

    <!-- Lista contatti -->
    <div v-if="filteredData.length && !loading">
      <ul class="contatti">
        <li class="contatto" v-for="item in filteredData" :key="item.id">
          <router-link :to="`/conversation/${item.name}`">
            <button id="user" class="btn btn-bg btn-outline-secondary">
              {{ item.name }}
            </button>
          </router-link>
        </li>
      </ul>
    </div>

    <div v-else-if="!loading && !filteredData.length">
      Nessun utente trovato.
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      data: null,
      loading: false,
      error: null,
      token: localStorage.getItem("token"),
      searchQuery: "" // ← Nuovo dato per la ricerca
    };
  },
  computed: {
    // Restituisce i risultati filtrati in base alla ricerca
    filteredData() {
      if (!this.searchQuery.trim()) {
        return this.data || [];
      }
      return this.data.filter(item =>
        item.name.toLowerCase().includes(this.searchQuery.toLowerCase())
      );
    }
  },
  methods: {
    async fetchData() {
      this.loading = true;
      this.error = null;

      if (!this.token) {
        this.error = "Token non trovato. Assicurati di essere loggato.";
        this.loading = false;
        return;
      }

      try {
        const response = await this.$axios.get("http://localhost:3000/user", {
          headers: {
            Authorization: `Bearer ${this.token}`
          }
        });
        this.data = response.data;
      } catch (e) {
        this.error = "Errore durante il recupero dei dati";
      } finally {
        this.loading = false;
      }
    }
  },
  mounted() {
    this.fetchData();
  }
};
</script>

<style>
.error {
  color: red;
}
.contatti {
  list-style-type: none;
  padding-left: 0;
}
.contatto {
  margin: 10px 0;
}
#user {
  width: 100%;
}
</style>
