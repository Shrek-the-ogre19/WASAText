<script>
import LoadingSpinner from "@/components/LoadingSpinner.vue";
import ErrorMsg from "@/components/ErrorMsg.vue";

export default {
	components: {ErrorMsg, LoadingSpinner},
	data: function() {
		return {
			errormsg: null,
			loading: false,
			users: null,
			path: this.$route.path,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("mainpage/0/users");
				this.users = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		}
	},
	mounted() {
		this.refresh()
	}
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
	<div v-for="user in users">
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
