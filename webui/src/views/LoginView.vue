<script>
import ErrorMsg from "@/components/ErrorMsg.vue";

export default {
	components: { ErrorMsg },
	data: function() {
		return {
			errormsg: null,
			username: null,
			userId: null
		}
	},
	methods: {
		async login(username) {
			if (!username || username.trim() === '') {
				this.errormsg = 'Username is required'
				return
			}

			try {
				let response = await this.$axios.post("/session", {name: username});
				const token = response.data.identifier;
				const payload = JSON.parse(atob(token.split(".")[1]));
				this.userId = payload.uid ?? payload.sub;
				localStorage.setItem("identifier", token);
				localStorage.setItem("id", this.userId);
				this.$router.push(`/mainpage/${this.userId}/conversations`)
			} catch (e) {
				this.errormsg = e.response?.data?.error ?? e.message ?? "Login failed";
			}
		}
	}
}
</script>

<template>
	<h1>Log in / Sign up</h1>
	<div v-if="errormsg">
		<ErrorMsg :msg="errormsg" />
	</div>
	<input v-model="username" placeholder="username" class="input"/>
	<br>
	<p>
		<button @click="login(username)" class="button">Log in</button>
	</p>
</template>

<style scoped>
</style>
