<script>

import ReplacingButton from "./ReplacingButton.vue";
import 'emoji-picker-element';
export default {
	name: 'MessageModal',
	components: {ReplacingButton},
	props: {
		path: String,
		showModal: Boolean,
		emojis: Array
	},
	emits: ['close', 'save'],
	data: function() {
		return {
			showEmojiModal : false,
			showEmojiPicker: false
		}},
	methods:{
		async emojiClick(event) {
			await this.$axios.post(`${this.path}/comments`,{comment: event.detail.unicode})
			this.showEmojiPicker = false
			this.save()
		},
		async uncomment(id){
			await this.$axios.delete(`${this.path}/comments/${id}`)
			this.showEmojiModal = false
			this.save()
		},
		async forwardMessage(receiver){
			await this.$axios.post(this.path,{conversation: receiver})
			this.save()
		},
		async deleteMessage(){
			await this.$axios.delete(this.path)
			this.save()
		},
		close() {
			this.$emit('close');
		},

		save() {
			this.$emit('save');
		},
	}
}

</script>

<template>
	<div v-if="showModal" class="modal">
		<div class="modal-content">
			<span class="close" @click="close">&times;</span>
			<ReplacingButton
				item-name="recipient"
				item="Forward"
				@save="forwardMessage"
			/>
			<button class="deleteButton" @click="deleteMessage">DELETE MESSAGE</button>

			<button @click="showEmojiPicker=true">+</button>
			<div v-if="showEmojiPicker">
				<emoji-picker @emoji-click="emojiClick"></emoji-picker>
			</div>
			<button @click="showEmojiModal=true">&#128512</button>
			<div v-if="showEmojiModal">
				<div v-for="emoji in emojis">
					<button @click="uncomment(emoji.Id.Id)">{{emoji.Content}}</button>
				</div>
			</div>
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

.modal-content .deleteButton {
	background-color: red !important;
	color: white !important;
	border: none;
	padding: 8px 16px;
	cursor: pointer;
}

.modal-content .deleteButton:hover {
	background-color: darkred !important;
}
</style>
