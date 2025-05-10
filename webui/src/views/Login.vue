<template>
	<div class="home">
		<h1>Login</h1>

		<!-- Form di login -->
		<form v-if="!isLoggedIn" @submit.prevent="login">
			<div>
				<label for="username">Nome:</label>
				<input class="form-control"type="text" v-model="name" id="username" required />
			</div>
			<button class="btn btn-outline-success btn-lg" type="submit">Accedi</button>
		</form>

		<!-- Info utente -->
		<form v-if="isLoggedIn" @submit.prevent="logout">
			<div>
				<p>Sei loggato come <strong>{{ name }}</strong></p>
			</div>
			<button class="btn btn-outline-danger btn-lg">Logout</button>
		</form>
		

		<!-- Messaggi -->
		<div v-if="message" class="message">{{ message }}</div>
		<div v-if="error" class="error">{{ error }}</div>
	</div>
</template>

<script>
export default {
	data() {
		return {
			token: '',
			name: '',
			message: '',
			error: ''
		};
	},
	computed: {
		isLoggedIn() {
			return this.token && this.name;
		}
	},
	mounted() {
		this.loadSession();
	},
	methods: {
		loadSession() {
			this.token = localStorage['token'];
			this.name = localStorage['name'];
		},
		async login() {
			this.message = '';
			this.error = '';

			try {
				const response = await this.$axios.post('http://localhost:3000/session', {
					name: this.name
				});

				const token = response.data.Authorization;

				if (token) {
					localStorage.setItem('token', token);
					localStorage.setItem('name', this.name);

					this.token = token;
					this.message = 'Login avvenuto con successo!';
				} else {
					this.error = 'Token non ricevuto dal server.';
				}
			} catch (err) {
				this.error = 'Errore durante il login.';
			}
		},
		logout() {
			localStorage.removeItem('token');
			localStorage.removeItem('name');
			this.token = '';
			this.name = '';
			this.message = 'Logout effettuato con successo.';
			this.error = '';
		}
	}
};
</script>

<style scoped>
form {
	display: flex;
	flex-direction: column;
	gap: 12px;
	max-width: 300px;
}

button {
	padding: 8px;
	font-size: 16px;
	cursor: pointer;
}

.message {
	color: green;
	margin-top: 10px;
}

.error {
	color: red;
	margin-top: 10px;
}
</style>
