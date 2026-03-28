<script>
export default{
	data: function(){
		return{
			errormsg:null,
			username:null,
			userId:null
		}
	},
	methods:{
		async login(username){
			if (!username || username.trim() === '') {
				this.errormsg = 'Username is required'
				return
			}

			try {
				let response = await this.$axios.post("/session", {username});
				this.userId = response.data;
				this.$router.push(`/mainpage/${this.userId}/conversations`)
			} catch (e) {
				this.errormsg = e.toString();
			}

		}
	}
}


</script>

<template>
	<h1>Log in / Sign up</h1>
	<input v-model="username" />
	<p>
		<button @click="login(username)">Log in</button>
	</p>

</template>

<style scoped>

</style>
