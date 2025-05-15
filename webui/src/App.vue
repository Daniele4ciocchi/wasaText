<script setup>
import { RouterLink, RouterView, useRoute } from 'vue-router'
import { ref, onMounted, onUnmounted, watch } from 'vue'
import axios from 'axios'

// Stato reattivo per le conversazioni
const route = useRoute()

const token = localStorage.getItem('token')
const username = localStorage.getItem('name')


watch(
	() => route.params.conversation_id,

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
								{{ username || "Login" }}
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
						<li class="nav-item">
							<RouterLink to="/myconversations" class="nav-link">
								<svg class="feather">
									<use href="/feather-sprite-v4.29.0.svg#message-square" />
								</svg>
								conversazioni
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

.conversation {
	display: flex;
	flex-direction: column;
	font-size: smaller;
}

#conversation-name {
	font-size: 15px;
}

#conversation-message {
	font-size: 12px;
}
</style>
