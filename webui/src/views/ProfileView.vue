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
			showPictureModal: false,
			selectedFile: null
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
		async onFileChanged (event) {
			this.selectedFile = event.target.files[0]
			let base64 = await this.fileToBase64(this.selectedFile)
			this.selectedFile=base64
			try {
				await this.$axios.put(this.path + "/settings/picture", {photo: this.selectedFile});
				await this.refresh();
			} catch (e) {
				this.errormsg = e.toString();
			}


		},
		fileToBase64(file) {
			return new Promise((resolve, reject) => {
				const reader = new FileReader();

				reader.onload = () => resolve(reader.result);
				reader.onerror = (error) => reject(error);

				reader.readAsDataURL(file);
			});
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


	</div>
	<img :src="user?.Picture" class="img" alt="userPicture"> <br>
	<input type="file" @change="onFileChanged">



</template>

<style scoped>
img {
	height: 200px;
	width: 200px;
	border-radius: 50%;
}
</style>
