<script>
import ReplacingButton from "../components/ReplacingButton.vue";
import ErrorMsg from "@/components/ErrorMsg.vue";

export default {
	components:{
		ErrorMsg,
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
				this.errormsg = "Username already exists";
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
		},
	},
	mounted() {
		this.refresh()
	}

}
</script>

<template>
	<div v-if="errormsg">
		<ErrorMsg
			:msg = "errormsg"
		/>
	</div>
	<div class="profile">
		<div class="name">
		<ReplacingButton
			item-name="Username"
			:item="user?.Name"
			@save="changeName"
		/>
		</div>

<div class="picture">
	<div v-if="user?.Picture=='default'">
		<img src="/default-avatar-icon-of-social-media-user-vector.jpg" class="img" alt="userPicture"/>
	</div>
		<div v-else>
	<img :src="user?.Picture" class="img" alt="userPicture"> <br>
	</div>
	<input type="file" @change="onFileChanged">
</div>
	</div>



</template>

<style scoped>
.img {
	height: 200px;
	width: 200px;
	border-radius: 50%;
}
.profile {
	display: flex;
	justify-content: space-between;
	align-items: center;
	gap: 20px;
}

.name:deep(.mainButton) {
	flex: 1;
	font-size: 30px;
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
.name:deep(.mainButton:hover) {
	background-color: #3a3f44;
	transform: translateX(5px);
	border-left-color: #0dcaf0;
}
.picture {
	display: flex;
	align-items: center;
	gap: 15px;
}
</style>
