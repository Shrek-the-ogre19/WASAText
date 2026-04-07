<script>
import ReplacingButton from "../components/ReplacingButton.vue";

export default {
	components:{
		ReplacingButton
	},
	data: function() {
		return {
			errormsg: null,
			loading: false,
			path: this.$route.path,
			user: null,
			newUsername: null,
			newPicture: null,
			showUsernameModal: false,
			showPictureModal: false
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

		async changeName(newItem) {
			if (!newItem || newItem.trim() === '') {
				this.errormsg = 'Username cannot be empty';
				return;
			}

			try {
				await this.$axios.put(this.path + "/settings/name", {name: newItem});
				await this.refresh();
			} catch (e) {
				this.errormsg = e.toString();
			}
		},

		async changePicture(newItem) {
			if (!newItem) {
				this.errormsg = 'Picture URL cannot be empty';
				return;
			}

			try {
				await this.$axios.put(this.path + "/settings/picture", {photo: newItem});
				await this.refresh();
			} catch (e) {
				this.errormsg = e.toString();
			}
		},},
		mounted() {
			this.refresh()
		}

}
</script>

<template>
	<div>
		<ReplacingButton
			item-name="Username"
			:item="user?.Name"
			@save="changeName"
		/>

		<ReplacingButton
			item-name="Picture"
			:item="user?.Picture"
			@save="changePicture"
		/>
	</div>



</template>

<style scoped>
</style>
