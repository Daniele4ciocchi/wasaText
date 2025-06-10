<template>
    <div class="users-list">
        <div class="users-list-header">
            <h2>Utenti nella conversazione</h2>
            <button>
                <svg class="feather">
                    <use href="/feather-sprite-v4.29.0.svg#user-plus" />
                </svg>
            </button>
        </div>
        <div class="users-list-content">
            <ul>
                <li v-for="member in members" :key="member.name">
                    {{ member.name }}
                </li>
            </ul>
        </div>
    </div>
</template>

<script>
export default {
    name: 'MembersList',
    props: {
        conversation: {
            type: Object,
            required: true,
        },
        conversationID: {
            type: [String, Number],
            required: true,
        },
        token: {
            type: String,
            default: localStorage.getItem("token"),
        },
    },
    data() {
        return {
            members: [],
        };
    },
    mounted() {
        this.getMembers();
    },
    methods: {
        async getMembers() {
            try {
                const response = await this.$axios.get(`/conversation/${this.conversationID}/member`, {
                    headers: {
                        Authorization: `Bearer ${this.token}`,
                    },
                });
                this.members = response.data;
            } catch (error) {
                console.error("Errore nel recupero degli utenti:", error);
            }
        },
    },
};
</script>

<style scoped>
.users-list {
    border: 1px solid #888;
    padding: 10px;
    overflow-y: auto;
    margin-top: 10px;
    background-color: #f4f6f8;
    display: flex;
    flex-direction: column;
    border-radius: 23px;
}

.users-list-content {
    margin-top: 20px;
    padding: 10px;
    border: 1px solid #888;
    border-radius: 15px;
    background-color: #f4f6f8;
    justify-content: center;
}

.users-list-content ul {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.users-list-content li {
    padding: 8px 12px;
    border-radius: 8px;
    background-color: #e9ecef;
    border: 1px solid #bdbdbd;
    font-size: 1em;
    color: #333;

}


.users-list-header {
    display: flex;
    align-items: left;
    margin-bottom: 10px;
    gap: 10px;
}
</style>
