<script>
import ReplacingButton from "./ReplacingButton.vue";
import { GROUP_DEFAULT_PICTURE, isDefaultPicture } from "../services/axios.js";

export default {
	name: 'GroupSettings',
	components: {ReplacingButton},
	props: {
		showModal: Boolean,
		name: String,
		picture: String,
		path: String,

	},
	emits: ['close', 'save'],
	data() {
		return {
			members: null,
			newMember: null
		}
	},
	methods: {
		async refresh(){
				this.loading = true;
				this.errormsg = null;
				try {
					let response = await this.$axios.get(this.path+"/groupMembers");
					this.members = response.data
				} catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
		},
		close() {
			this.$emit('close');
		},

		save() {
			this.$emit('save');
		},
		async addNewMember(newItem){
			await this.$axios.post(this.path+"/groupMembers", {name: newItem});
			await this.refresh();
			this.save()
		},
		async leave(){
			await this.$axios.delete(this.path+"/groupMembers")
			let result = this.path.split('/').slice(0, 3).join('/');
			this.close()
			this.$router.push(result)
		},

		async changeName(newItem) {
			if (!newItem || newItem.trim() === '') {
				this.errormsg = 'Group name cannot be empty';
				return;
			}

			try {
				await this.$axios.put(this.path + "/conversationsettings/groupname", {name: newItem});
				await this.refresh();
				this.save()
			} catch (e) {
				this.errormsg = e.toString();
			}
		},
		async onFileChanged (event) {
			this.selectedFile = event.target.files[0]
			let base64 = await this.fileToBase64(this.selectedFile)
			this.selectedFile=base64
			try {
				await this.$axios.put(this.path + "/conversationsettings/grouppicture", {picture: this.selectedFile});
				await this.refresh();
				this.save()
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
	},
	computed: {
		displayPicture() {
			return isDefaultPicture(this.picture) ? GROUP_DEFAULT_PICTURE : this.picture;
		},
	},
}
</script>

<template>
	<div v-if="showModal" class="modal">
		<div class="modal-content">
			<span class="close" @click="close">&times;</span>
			<h1>Group Settings</h1>
			<ReplacingButton
				item-name="Group Name"
				:item="name"
				@save="changeName"
			/>
			<img :src="displayPicture" class="img" alt="chatPicture"> <br>
			<input type="file" @change="onFileChanged">


			<h2>Members:</h2>
			{{members}}
			<ReplacingButton
				item-name="member"
				item="Add Member"
				@save="addNewMember"
			/>
			<button class="redButton" @click="leave">LEAVE GROUP</button>
		</div>
	</div>

</template>

<style scoped>
.modal {
	position: fixed;
	z-index: 1000;
	left: 0;
	top: 0;
	width: 100%;
	height: 100%;
	background-color: rgba(0, 0, 0, 0.5);
	display: flex;
	justify-content: center;
	align-items: center;
}

.modal-content {
	background-color: white;
	padding: 20px;
	border-radius: 8px;
	width: 300px;
	position: relative;
	box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.close {
	position: absolute;
	top: 10px;
	right: 15px;
	font-size: 24px;
	cursor: pointer;
	color: #aaa;
}

.close:hover {
	color: #000;
}

.modal-content p {
	margin-top: 0;
	font-weight: bold;
}

.modal-content input {
	width: 100%;
	padding: 8px;
	margin: 10px 0;
	border: 1px solid #ddd;
	border-radius: 4px;
}

.modal-content button {
	padding: 8px 16px;
	margin-right: 10px;
	background-color: #007bff;
	color: white;
	border: none;
	border-radius: 4px;
	cursor: pointer;
}

.modal-content button:hover {
	background-color: #0056b3;
}

.error {
	color: red;
	padding: 10px;
	margin: 10px 0;
	background-color: #ffeeee;
	border-radius: 4px;
}

.modal-content .redButton {
	background-color: red !important;
	color: white !important;
	border: none;
	padding: 8px 16px;
	cursor: pointer;
}

.modal-content .redButton:hover {
	background-color: darkred !important;
}
</style>
