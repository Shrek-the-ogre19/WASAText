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
	<li v-for="item in users">
		{{ item }}
	</li>
</template>

<style scoped>

</style>
