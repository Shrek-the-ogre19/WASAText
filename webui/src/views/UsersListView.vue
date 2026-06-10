<script>
import LoadingSpinner from "@/components/LoadingSpinner.vue";
import ErrorMsg from "@/components/ErrorMsg.vue";
import { startAutoRefresh } from "@/services/refresh.js";

export default {
	components: {ErrorMsg, LoadingSpinner},
	data: function() {
		return {
			errormsg: null,
			loading: false,
			users: [],
			path: this.$route.path,
			stopAutoRefresh: null,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				const token = localStorage.getItem("identifier");
				const userId = localStorage.getItem("id");
				const url = token && userId ? `/mainpage/${userId}/users` : "/users";
				let response = await this.$axios.get(url);
				this.users = response.data ?? [];
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		}
	},
	mounted() {
		this.refresh();
		this.stopAutoRefresh = startAutoRefresh(() => this.refresh());
	},
	beforeUnmount() {
		if (this.stopAutoRefresh) {
			this.stopAutoRefresh();
		}
	},
}
</script>

<template>
	<div v-if="loading">
		<LoadingSpinner
			:loading = "loading"
		/>
	</div>
	<div v-if="errormsg">
		<ErrorMsg
			:msg = "errormsg"
		/>
	</div>
	<div v-for="user in users" :key="user">
		<div class="user">
			{{ user }}
		</div>
	</div>
</template>

<style scoped>
.user {
	font-size: 16px;
	color: #e9ecef;
	background-color: #2c3034;
	padding: 12px 20px;
	margin: 8px 0;
	border-radius: 8px;
	border-left: 4px solid #0d6efd;
	transition: all 0.3s ease;
	list-style: none;
	font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
}

</style>
