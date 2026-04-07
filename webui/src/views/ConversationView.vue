<script>
import GroupSettings from "../components/GroupSettings.vue";
import Message from "../components/Message.vue";
import ConversationHead from "../components/ConversationHead.vue";

export default {
	components: {
		ConversationHead, GroupSettings,
		Message,
	},
	data: function() {
		return {
			path: this.$route.path,
			showModal: false,
			errormsg: null,
			loading: false,
			text: null,
			groupchat: false,
			name: null,
			picture: null,
			messages: [],
			members: [],
			showSettingsModal: false
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
				this.messages = this.conversation.Content;
				this.picture = this.conversation.Picture;
				this.members = this.conversation.Members;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async sendMessage(text){
			try{
				await this.$axios.post(this.path, {content: text});
				this.showModal = false
			}catch(e){
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
		}
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
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
		{{name}}

	</header>



	<li v-for="message in messages">
		<Message
			:path=path
			:messageId=message.Id
			@save="refresh()"
		/>
	</li>

	<div class = "bottom">
			<button @click="showModal = true">+</button>
	</div>
	<div v-if="showModal" class="modal">
		<div class="modal-content">
			<span class="close" @click="showModal = false">&times;</span>
			<p>Send a new message:</p>
			<input v-model="text" placeholder="text" />
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
</style>
