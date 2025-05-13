<script setup>
import { RouterLink, RouterView, useRoute } from 'vue-router'
import { ref, onMounted, onUnmounted, watch } from 'vue'
import axios from 'axios'

// Stato reattivo per le conversazioni
const route = useRoute()
const conversations = ref([])
const messages = ref([])
const error = ref(null)
const token = localStorage.getItem('token')
const username = localStorage.getItem('name')
let intervalID


const getConversations = async () => {
	try {
		const res = await axios.get('http://100.87.168.104:3000/conversation', {
			headers: { Authorization: `Bearer ${token}` }
		})
		conversations.value = res.data
	} catch (err) {
		error.value = 'Errore nel caricamento delle conversazioni'
		console.error(err)
	}
}
const fetchMessages = async () => {
  const convId = route.params.conversationID
  try {
    const res = await axios.get(`http://100.87.168.104:3000/conversation/${convid}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    messages.value = res.data.messages
  } catch (err) {
    console.error('Errore nel caricamento dei messaggi:', err)
  }
}
onMounted(async () => {

	intervalID = setInterval(getConversations, 3000)
	fetchMessages()
})

onUnmounted(() => {
	clearInterval(intervalID)
})
watch(
  () => route.params.conversation_id,
  () => {
    fetchMessages()
  },
  { immediate: true }
)

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
								{{ username || "Login"}}
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
						<li class="nav-item">
							<RouterLink to="/group" class="nav-link">
								<svg class="feather">
									<use href="/feather-sprite-v4.29.0.svg#users" />
								</svg>
								nuovo gruppo
							</RouterLink>
						</li>

					</ul>

					<h6
						class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>Conversazioni</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item" v-for="conv in conversations" :key="conv.conversation_id">
							<RouterLink :to="'/conversation/' + conv.conversation_id" :key="route.params.conversation_id" class="nav-link">
								<svg class="feather">
									<use href="/feather-sprite-v4.29.0.svg#book" />
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

<style>
.nav-item {
	background-color: #d1e7dd;
	border-radius: 13px;
	border: #888 1px solid;
	margin: 5px 10px;
}
</style>
