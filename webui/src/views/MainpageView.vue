<script>

import ConversationHead from "../components/ConversationHead.vue";
import ErrorMsg from "@/components/ErrorMsg.vue";
import { startAutoRefresh } from "../services/axios.js";

export default {
	components:{
		ConversationHead,
		ErrorMsg,
	},
	data: function() {
		return {
			errormsg: null,
			loading: false,
			path: this.$route.path,
			conversations: [],
			showModal: false,
			id: null,
			receivers: null,
			conversationId: null,
			stopAutoRefresh: null,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get(`/mainpage/${this.$route.params.Id}/conversations`);
				this.conversations = [...(response.data ?? [])].reverse()
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async startConversation(receivers) {
			try {
				const data = (await this.$axios.post(`/mainpage/${this.$route.params.Id}/conversations`, {receivers: receivers})).data;
				this.conversationId = data.Id?.Id ?? data.Id;
				this.showModal = false;
				this.$router.push(`${this.path}/${this.conversationId}`)
			} catch (e) {
				this.errormsg = e.toString();
			}
		}
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
	<div v-for="conversation in conversations" :key="conversation.Id?.Id ?? conversation.Id" class="conversation">
		<ConversationHead
			:path=path
			:conversation="conversation"
		/>
	</div>

</template>

<style scoped>
.conversation {
	width: 900px;
	height: 120px;
	overflow: auto;
}
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
.button {
	position: fixed;
	top: 70px;
	left: 220px;
	z-index: 1000;

	/* Button styling */
	padding: 10px 20px;
	background-color: forestgreen;
	color: white;
	border: none;
	border-radius: 8px;
	cursor: pointer;
	font-size: 16px;
	font-weight: 500;
	transition: all 0.3s ease;
}

.button:hover {
	background-color: darkgreen;
	transform: scale(1.05);
}

</style>
