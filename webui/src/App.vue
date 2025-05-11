<script setup>
import { RouterLink, RouterView } from 'vue-router'
import { ref, onMounted } from 'vue'
import axios from 'axios'

// Stato reattivo per le conversazioni
const conversations = ref([])
const error = ref(null)

onMounted(async () => {
	try {
		const token = localStorage.getItem('token')
		const res = await axios.get('http://localhost:3000/conversation', {
			headers: { Authorization: `Bearer ${token}` }
	})
		conversations.value = res.data
	} catch (err) {
		error.value = 'Errore nel caricamento delle conversazioni'
		console.error(err)
	}
})
</script>

<script>
export default {}
</script>


<template>

	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WasApp</a>
		<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse"
			data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false"
			aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
	</header>

	<div class="container-fluid">
		<div class="row">
			<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div class="position-sticky pt-3 sidebar-sticky">
					<h6
						class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>General</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink to="/" class="nav-link">
								<svg class="feather" id="Login">
									<use href="/feather-sprite-v4.29.0.svg#user" />
								</svg>
								Login
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/users" class="nav-link">
								<svg class="feather">
									<use href="/feather-sprite-v4.29.0.svg#book" />
								</svg>
								nuova conversazione
							</RouterLink>
						</li>

					</ul>

					<h6
						class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>Conversazioni</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item" v-for="conv in conversations" :key="conv.id">
							<RouterLink :to="`/conversation/${conv.conversation_id}`" class="nav-link">
								<svg class="feather">
									<use href="/feather-sprite-v4.29.0.svg#message-square" />
								</svg>
								{{ conv.name }}
							</RouterLink>
						</li>
					</ul>

				</div>
			</nav>

			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<RouterView />
			</main>
		</div>
	</div>
</template>

<style></style>
