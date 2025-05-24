
<template>
	<div class="home">

		<!-- Form di login -->
		<form v-if="!isLoggedIn" @submit.prevent="login">
			<h1>Login</h1>
			<div>
				<label for="name">Nome:</label>
				<input class="form-control" type="text" v-model="name" id="username" required />
			</div>
			<button class="btn btn-outline-success btn-lg" type="submit">Accedi</button>
		</form>

		<!-- Info utente -->
		<form v-if="isLoggedIn">
			<h2>
				{{ username }}
				<button v-if="changeUsername == false" class="btn btn-link" @click="changeUsername = true">
					<svg class="feather">
						<use href="/feather-sprite-v4.29.0.svg#edit-3" />
					</svg>
				</button>
			</h2>
			<div>
				<img :src="photoUrl" alt="Foto profilo" style="width: 150px; height: 150px; border-radius: 50%;" />
				<p>nome :  <strong>{{ name }}</strong></p>
			</div>


		</form>

		<form v-if="isLoggedIn" @submit.prevent="uploadProfilePicture">
			<div>
				<label for="profilePicture">Carica una foto profilo:</label>
				<input type="file" accept="image/jpeg" @change=uploadImage>

			</div>
			<button class="btn btn-primary btn-lg" type="submit">Carica</button>
		</form>

		<form v-if="isLoggedIn && changeUsername">

			<div class="username">
				<div>
					<label for="username">username:</label>
					<input class="form-control" id="username" v-model="newusername" required />

				</div>
				<button class="btn btn-secondary btn-lg username-button" @click.prevent="submitUsernameChange">
					cambia il tuo username
				</button>

			</div>
		</form>

		<button v-if="isLoggedIn" @click="logout" class="btn btn-outline-danger btn-lg ">
			Logout
		</button>

		<!-- Messaggi -->
		<div v-if="message" class="message">{{ message }}</div>
		<div v-if="error" class="error">{{ error }}</div>
	</div>
</template>

<script>
export default {
	data() {
		return {
			changeUsername: false,
			newusername: '',
			username: '',
			token: '',
			name: '',
			user_id: '',
			photoUrl: '',
			selectedFile: null,
			previewImage: null,
			message: '',
			error: ''
		};
	},
	computed: {
		isLoggedIn() {
			return this.token;
		}
	},
	mounted() {
		this.loadSession();
	},
	methods: {
		uploadImage(e) {
			const file = e.target.files[0];
			this.selectedFile = file;   // salva il file originale
			const reader = new FileReader();
			reader.readAsDataURL(file);
			reader.onload = e => {
				this.previewImage = e.target.result;  // questa è la base64 per l'anteprima
				console.log(this.previewImage);
			};
		},
		loadSession() {
			this.token = localStorage.getItem('token');
			this.name = localStorage.getItem('name');
			this.username = localStorage.getItem('username');
			this.user_id = localStorage.getItem('user_id');
			this.photoUrl = this.loadPhoto();
		},
		async login() {
			this.message = '';
			this.error = '';

			try {
				const response = await this.$axios.post(`/session`, {
					name: this.name
				});

				const token = response.data.Authorization;

				if (token) {
					localStorage.setItem('token', token);
					this.token = token;
					this.message = 'Login avvenuto con successo!';
					this.getMe();

				} else {
					this.error = 'Token non ricevuto dal server.';
				}

				try {
					const response = await this.$axios.get('/me/newmessage', {
						headers: {
							Authorization: `Bearer ${this.token}`
						}
					});
				} catch (err) {
					this.error = 'Errore durante il recupero dei messaggi.';
				}

			} catch (err) {
				this.error = 'Errore durante il login.';
			}

		},
		logout() {
			localStorage.removeItem('token');
			localStorage.removeItem('name');
			localStorage.removeItem('username');
			localStorage.removeItem('user_id');
			this.token = '';
			this.name = '';
			this.username = '';
			this.user_id = '';
			this.message = 'Logout effettuato con successo.';
			this.error = '';
		},
		async getMe() {

			try {
				const userResponse = await this.$axios.get(`/me`, {
					headers: {
						Authorization: `Bearer ${this.token}`
					}
				});
				this.username = userResponse.data.username;
				this.name = userResponse.data.name;
				this.user_id = userResponse.data.user_id;
				this.loadPhoto();

			} catch (err) {
				this.error = 'Errore durante il recupero dell\'username.';
			}
			localStorage.setItem('name', this.name);
			localStorage.setItem('username', this.username);
			localStorage.setItem('user_id', this.user_id);
		},
		async submitUsernameChange() {
			this.message = '';
			this.error = '';

			if (!this.newusername.trim()) {
				this.error = 'Il nome utente è obbligatorio.';
				return;
			}

			try {
				const response = await this.$axios.put(`/me/username`, {
					username: this.newusername
				}, {
					headers: {
						Authorization: `Bearer ${this.token}`
					}
				});
				this.getMe();
				if (response.status === 200) {
					this.message = 'Username cambiato con successo!';
					this.changeUsername = false;
					this.newusername = '';
				} else {
					this.error = 'Errore durante il cambio username.';
				}
			} catch (err) {
				this.error = 'Errore durante il cambio username.';
			}
		},
		async loadPhoto() {
			try {
				const response = await this.$axios.get(`/user/${this.user_id}/photo`, {
					headers: {
						Authorization: `Bearer ${this.token}`
					},
					responseType: 'blob' // per gestire immagine binaria
				});

				this.photoUrl = URL.createObjectURL(response.data);
			} catch (err) {
				this.error = 'Errore durante il caricamento della foto profilo.';
			}
		},
		async uploadProfilePicture() {
			this.message = '';
			this.error = '';

			if (!this.selectedFile) {
				this.error = 'Seleziona un file prima di caricare.';
				return;
			}

			const formData = new FormData();
			formData.append('photo', this.selectedFile);  // metti il file originale qui

			try {
				const response = await this.$axios.put(`/me/photo`, formData, {
					headers: {
						Authorization: `Bearer ${this.token}`,
						"Content-Type": "multipart/form-data",
					},
				});
				this.message = 'Foto profilo caricata con successo!';
				this.loadPhoto();
			} catch (err) {
				this.error = 'Errore durante il caricamento della foto profilo.';
			}
		},
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

.username {
	margin-top: 10px;
}

.username-button {
	margin-top: 10px;
}
</style>
