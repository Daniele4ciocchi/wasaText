<template>
    <div>
        <h1>Users List</h1>
        <ul v-if="users.length">
            <li v-for="user in users" :key="user.id">{{ user.name }}</li>
        </ul>
        <p v-else>Loading users...</p>
    </div>
</template>

<script>
import axios from 'axios';

export default {
    name: 'UsersList',
    data() {
        return {
            users: [],
        };
    },
    created() {
        this.fetchUsers();
    },
    methods: {
        async fetchUsers() {
            try {
                const response = await axios.get('http://localhost:3000/users');
                this.users = response.data;
            } catch (error) {
                console.error('Error fetching users:', error);
            }
        },
    },
};
</script>

<style scoped>
h1 {
    font-size: 1.5em;
    margin-bottom: 1em;
}
ul {
    list-style-type: none;
    padding: 0;
}
li {
    margin: 0.5em 0;
}
</style>