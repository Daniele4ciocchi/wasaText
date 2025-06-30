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
		<div v-if="isLoggedIn">
			<h2>Profilo</h2>
			<div class="user-profile">
				<img :src="photoUrl" alt="Foto profilo"
					style="width: 150px; height: 150px; border-radius: 15px; border: #888 1px solid;" />
				<div class="user-info">
					<p>nome : <strong>{{ name }}</strong></p>
					<p>username : <strong>{{ username }}</strong></p>
					<button v-if="isLoggedIn" @click="logout" class="logout-button">
						Logout
					</button>
				</div>
			</div>

			<div class="change">
				<div class="change-username">
					<div v-if="!changeUsername">
						<h4 class="username">Cambia username</h4>
						<button class="username-button" @click="changeUsername = true">Cambia</button>
					</div>
					<div>
						<form v-if="changeUsername" @submit.prevent="submitUsernameChange">
							<label for="newusername">Nuovo username:</label>
							<input class="form-control" type="text" v-model="newusername" id="newusername" required />
							<button class="btn btn-outline-success btn-lg" type="submit">Invia</button>
						</form>
						<button v-if="changeUsername" class="username-button"
							@click="changeUsername = false">Annulla</button>
					</div>


				</div>
				<div class="change-photo">
					<h4 class="username">Cambia foto profilo</h4>
					<input type="file" @change="uploadImage" accept="image/*" />
					<button @click="uploadProfilePicture">Carica Foto</button>

				</div>
			</div>
		</div>




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
			if (this.token) this.photoUrl = this.loadPhoto();
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
				}
			} catch (err) {
				if (err.response && err.response.status === 409) {
					this.error = 'Username già in uso.';
				} else {
					this.error = 'Errore durante il cambio dell\'username.';
				}
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

.change {
	display: flex;
	gap: 20px;
	margin-top: 20px;
	border: #888 1px solid;
	border-radius: 17px;
	padding: 10px;
	width: auto;
	display: flex;

}

.change-username,
.change-photo {
	padding: 10px;
	border: #888 1px solid;
	border-radius: 15px;
	background-color: #f4f6f8;
	width: 100%;
}

.logout-button {
	background-color: #f44336;
	color: white;
	border: none;
	padding: 10px 20px;
	cursor: pointer;
	width: 100%;
}

.user-profile {
	margin-top: 20px;
	display: flex;
	gap: 20px;
}

.user-info {
	padding: 10px;
	width: 400px;
	border: #888 1px solid;
	border-radius: 15px;
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
