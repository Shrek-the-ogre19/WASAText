<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			path: this.$route.path,
			user: null,
			newUsername: null,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get(this.path);
				this.user = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},

		async changeName(newUsername) {
			if (!newUsername || newUsername.trim() === '') {
				this.errormsg = 'Username cannot be empty';
				return;
			}

			try{
				await this.$axios.put(this.path + "/setting/name", {name: newUsername});
				await this.refresh();
			} catch (e){
				this.errormsg = e.toString();
			}
		},

		async changePicture(newPicture) {
			if (!newPicture) {
				this.errormsg = 'Picture URL cannot be empty';
				return;
			}

			try{
				await this.$axios.put(this.path + "/setting/picture", {photo: newPicture});
				await this.refresh();
			} catch (e){
				this.errormsg = e.toString();
			}
		},
		async goBack(){
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
