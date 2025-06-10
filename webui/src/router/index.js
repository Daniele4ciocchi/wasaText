import {createRouter, createWebHashHistory} from 'vue-router'

import Login from '../views/Login.vue'
import Users from '../views/Users.vue'
import CreateGroup from '../views/CreateGroup.vue'
import Conversation from '../views/Conversation.vue'
import MyConversations from '../views/MyConversations.vue'





const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: Login},
		//{path: '/link1', component: HomeView},
		//{path: '/link2', component: HomeView},

		{path: '/users', component: Users},
		{path: '/group', component: CreateGroup},
		{path: '/conversation/:conversationID', component: Conversation},
		{path: '/myconversations', component: MyConversations},

	]
})

export default router
