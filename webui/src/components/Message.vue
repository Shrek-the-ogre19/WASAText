<script>
import GroupSettings from "./GroupSettings.vue";
import MessageModal from "./MessageModal.vue";

export default {
	name: 'Message',
	components: {MessageModal, GroupSettings},
	props: {
		path: String,
		messageId: String
	},
	emits: ['save'],
	data: function() {
		return {
			errormsg: null,
			loading: false,
			message: null,
			showMessageModal: false,
			comments : [],
			emojis:[],
			text: null,
			image: null,
			senderName: null,
			received: true
		}
	},methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			this.emojis=[]
			try {
				let response = await this.$axios.get(`${this.path}//${this.messageId}`);
				this.message = response.data;
				let str = this.message.Content;
				let parts = str.split('data:image');
				this.text = parts[0]
				if (parts[1]){
				this.image = 'data:image' + parts[1]}
				response = await this.$axios.get(`/mainpage/0/users/${this.message.Sender.Id}`)
				this.senderName = response.data.Name
				let start = this.path.indexOf('/mainpage/') + '/mainpage/'.length;
				let end = this.path.indexOf('/', start);
				let currentUser = this.path.substring(start, end);
				if (currentUser == this.message.Sender.Id){
					this.received=false
				}
				response = await this.$axios.get(`${this.path}//${this.messageId}/comments`);
				this.comments = response.data;
				for (let i = 0; i < this.comments.length; i++) {
					let id = this.comments[i].Id;
					response =await this.$axios.get(`${this.path}//${this.messageId}/comments/${id}`);
					this.emojis.push(response.data)
				}

			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async save(){
			this.showMessageModal = false
			await this.refresh()
			this.$emit('save')
		},
		async emoji(id){
			let response =await this.$axios.get(`${this.path}//${this.messageId}/comments/${id}`);
			return response.data.Content
		}
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<div>
		<div :class="['message-container', received ? 'received' : 'sent']">
			<button class="message-bubble" @click="showMessageModal = true">
				<!-- Top part: Name -->
				<div class="message-header">
					<span class="sender-name">{{ senderName }}</span>
				</div>

				<!-- Biggest part: Content (text and image) -->
				<div class="message-content">
					<div v-if="text" class="message-text">{{ text }}</div>
					<div v-if="image" class="message-image">
						<img :src="image" class="img" alt="messagePicture">
					</div>
				</div>

				<!-- Bottom part: Status and Timestamp -->
				<div class="message-footer">
					<div v-if="message?.Status=='sent'">
						<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#check"/></svg>
					</div>
					<div v-else>
						<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#check"/></svg>
						<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#check"/></svg>
					</div>
					<span class="message-status">{{ message?.Status }}</span>
					<span class="message-timestamp">{{ message?.Timestamp }}</span>
				</div>

				<!-- Emojis -->
				<div v-if="emojis.length" class="message-emojis">
					<div v-for="emoji in emojis.slice(-3)" :key="emoji.Id" class="emoji">
						{{ emoji.Content }}
					</div>
				</div>
			</button>
		</div>

		<MessageModal
			:showModal="showMessageModal"
			:path="`${this.path}//${this.messageId}`"
			:emojis="emojis"
			@close="showMessageModal = false"
			@save="save"
		/>
	</div>
</template>

<style scoped>
/* Container for positioning on left/right */
.message-container {
	display: flex;
	margin-bottom: 16px;
	width: 100%;
}

.message-container.received {
	justify-content: flex-start;
}

.message-container.sent {
	justify-content: flex-end;
}

/* Message bubble as button */
.message-bubble {
	max-width: 70%;
	min-width: 200px;
	padding: 12px 16px;
	border: none;
	border-radius: 12px;
	box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
	transition: all 0.2s ease;
	cursor: pointer;
	text-align: left;
	font-family: inherit;

	/* Same color for both sides */
	background-color: #f1f3f5;
	color: #212529;
}

/* Remove default button styles */
.message-bubble:focus {
	outline: none;
}

.message-bubble:active {
	transform: scale(0.98);
}

/* Hover effect */
.message-bubble:hover {
	transform: translateY(-2px);
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
	background-color: #e9ecef;
}

/* Message header (Name) */
.message-header {
	margin-bottom: 8px;
}

.sender-name {
	font-size: 16px;
	font-weight: 600;
	display: block;
	color: #495057;
}

/* Message content (biggest part) */
.message-content {
	margin-bottom: 8px;
	word-wrap: break-word;
}

.message-text {
	font-size: 18px; /* Increased font size */
	line-height: 1.5;
	margin-bottom: 8px;
	color: #212529;
}

.message-image {
	margin-top: 8px;
}

.message-image .img {
	max-width: 100%;
	max-height: 300px;
	border-radius: 8px;
	object-fit: cover;
}

/* Message footer (Status and Timestamp) */
.message-footer {
	display: flex;
	justify-content: flex-end;
	align-items: center;
	gap: 8px;
	font-size: 12px;
	margin-top: 4px;
	color: #6c757d;
}

.message-status {
	text-transform: capitalize;
}

.message-timestamp {
	font-size: 11px;
}

/* Emojis section */
.message-emojis {
	display: flex;
	gap: 6px;
	margin-top: 8px;
	padding-top: 8px;
	border-top: 1px solid rgba(0, 0, 0, 0.1);
}

.emoji {
	font-size: 16px;
	cursor: pointer;
	padding: 2px 4px;
	border-radius: 12px;
	transition: all 0.2s;
	background-color: rgba(0, 0, 0, 0.05);
}

.emoji:hover {
	transform: scale(1.1);
	background-color: rgba(0, 0, 0, 0.1);
}

/* Remove the old message-button styles since it's no longer needed */
.message-button {
	display: none;
}

/* Responsive design */
@media (max-width: 768px) {
	.message-bubble {
		max-width: 85%;
		min-width: 150px;
		padding: 10px 12px;
	}

	.sender-name {
		font-size: 14px;
	}

	.message-text {
		font-size: 16px; /* Slightly smaller on mobile but still larger */
	}

	.message-footer {
		font-size: 10px;
	}

	.message-timestamp {
		font-size: 9px;
	}
}
</style>
