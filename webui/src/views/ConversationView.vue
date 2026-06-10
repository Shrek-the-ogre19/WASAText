<script>
import GroupSettings from "../components/GroupSettings.vue";
import Message from "../components/Message.vue";
import ErrorMsg from "@/components/ErrorMsg.vue";
import { startAutoRefresh } from "../services/refresh.js";

export default {
	components: {
		GroupSettings,
		Message,
		ErrorMsg,
	},
	data: function() {
		return {
			path: this.$route.path,
			showModal: false,
			errormsg: null,
			loading: false,
			text: "",
			image: null,
			groupchat: false,
			name: null,
			messages: [],
			members: [],
			picture: null,
			showSettingsModal: false,
			selectedFile : null,
			stopAutoRefresh: null,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get(this.path);
				this.conversation = response.data;
				this.name = this.conversation.Name;
				this.groupchat = this.conversation.Groupchat;
				this.messages = [...this.conversation.Content].reverse();
				this.members = this.conversation.Members;
				this.picture = this.conversation.Picture;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},

		async sendMessage(text){
			if (this.image != null) {
				text = text + this.image;
			}
			try {
				await this.$axios.post(this.path, {content: text});
				this.showModal = false;
				this.text = "";
				this.image = null;
				this.selectedFile = null;
				await this.refresh();
			} catch (e) {
				this.errormsg = e.toString();
			}
		},
		async openSettings(){
			if (this.groupchat == true){
				this.showSettingsModal = true
			}
		},
		async saveSettings(){
			this.showSettingsModal=false;
			await this.refresh()
		},

		async onFileChanged (event) {
			this.selectedFile = event.target.files[0]
			let base64 = await this.fileToBase64(this.selectedFile)
			this.selectedFile=base64
			this.image=this.selectedFile
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
		this.refresh();
		this.stopAutoRefresh = startAutoRefresh(() => this.refresh());
	},
	beforeUnmount() {
		if (this.stopAutoRefresh) {
			this.stopAutoRefresh();
		}
	},
}
</script>

<template>
	<div v-if="errormsg">
		<ErrorMsg :msg="errormsg" />
	</div>
	<div>
		<div class="settings">
			<button @click = "openSettings">
				settings
			</button>
		</div>

		<div v-if="groupchat">
		<GroupSettings
			:showModal="showSettingsModal"
			:name="name"
			:picture="picture"
			:path="path"
			@close="showSettingsModal = false"
			@save="saveSettings"
		/>
		</div>
	</div>

	<header class="chatname">
		<div v-if="picture == 'default'">
			<img src="/default-avatar-icon-of-social-media-user-vector.jpg" class="img" alt="chatPicture"/>
		</div>
		<div v-else>
		<img :src="picture" class="img" alt="chatPicture">
		</div>
		{{name}}
	</header>


	<br>
	<br>
	<br>
	<br>
	<div v-for="message in messages" :key="message.Id.Id">
		<Message
			:path="path"
			:messageId="message.Id.Id"
			@save="refresh()"
		/>
	</div>

	<div class = "bottom">
			<button @click="showModal = true">+</button>
	</div>
	<div v-if="showModal" class="modal">
		<div class="modal-content">
			<span class="close" @click="showModal = false">&times;</span>
			<p>Send a new message:</p>
			<input v-model="text" placeholder="text" />
			<input type="file" @change="onFileChanged">
			<button @click="sendMessage(text)">CONFIRM</button>
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

.settings{
	text-align: justify;
	text-align-last: right;
}

.chatname{
	background-color: darkgreen;
	padding: 10px;
}

.bottom {
	position: fixed;
	bottom: 0;
}
.chatname {
	position: fixed;
	top: 45px;
	left: 220px;
	right: 0;
	z-index: 1000;

	/* Layout */
	display: flex;
	align-items: center;
	gap: 12px;

	/* Styling */
	padding: 12px 20px;
	background-color: forestgreen;
	box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
	border-bottom: 1px solid #e9ecef;

	/* Text styling */
	font-size: 30px;
	font-weight: 600;
	color: #212529;
}

.img {
	width: 60px;
	height: 60px;
	border-radius: 50%;
	object-fit: cover;
	border: 2px solid #ffffff;
	box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	flex-shrink: 0;
}
body {
	padding-top: 70px;
}
.bottom {
	position: fixed;
	bottom: 20px;
	right: 20px;
	z-index: 1000;
}

.bottom button {
	width: 60px;
	height: 60px;
	border-radius: 50%;
	background-color: #28a745;
	color: white;
	border: none;
	cursor: pointer;
	font-size: 32px;
	font-weight: bold;
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
	transition: all 0.3s ease;
	display: flex;
	align-items: center;
	justify-content: center;
}

.bottom button:hover {
	background-color: #218838;
	transform: scale(1.1);
	box-shadow: 0 6px 16px rgba(0, 0, 0, 0.3);
}
.settings {
	position: fixed;
	top: 45px;
	right: 20px;
	z-index: 9999;

	/* Button styling */
	padding: 10px 16px;
	background-color: #0d6efd;
	color: white;
	border: none;
	border-radius: 8px;
	cursor: pointer;
	font-size: 14px;
	font-weight: 500;
	transition: all 0.3s ease;
}
</style>
