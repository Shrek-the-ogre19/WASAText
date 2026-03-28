<script>
export default {
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
				let response = await this.$axios.get(this.path);
				this.users = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async goBack(){
			let pathParts = path.split('/');
			pathParts.pop();
			let newPath = pathParts.join('/');
			await this.$router.push(newPath+ "/conversations");
		}
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>

</template>

<style scoped>

</style>
