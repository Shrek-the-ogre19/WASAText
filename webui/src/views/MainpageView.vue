<script>

import ConversationHead from "../components/ConversationHead.vue";
import ReplacingButton from "../components/ReplacingButton.vue";

export default {
	components:{
		ReplacingButton,
		ConversationHead
	},
	data: function() {
		return {
			errormsg: null,
			loading: false,
			path: this.$route.path,
			conversations: null,
			showModal: false,
			id: null,
			receivers: null,
			conversationId: null
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get(this.path);
				this.conversations = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async openConversation(conversationId){
			this.$router.push(`${this.path}/${conversationId}`)
		},
		async startConversation(receivers) {
			try {
				this.conversationId = (await this.$axios.post(this.path, {receivers: receivers})).data.Id;
				this.showModal = false;
				this.$router.push(`${this.path}/${this.conversationId}`)
			} catch (e) {
				this.errormsg = e.toString();
			}
		}
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<p>
		<button class="button" @click="showModal = true">+</button>
	</p>
	<div v-if="showModal" class="modal">
		<div class="modal-content">
			<span class="close" @click="showModal = false">&times;</span>
			<p>Start new conversation:</p>
			<input v-model="receivers" placeholder="receivers" />
			<button @click="startConversation(receivers)">CONFIRM</button>
		</div>
	</div>
	<li v-for="conversation in conversations">
		<ConversationHead
			:path=path
			:conversation="conversation"
		/>
	</li>

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
</style>
