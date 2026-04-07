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
			emojis:[]
		}
	},methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			this.emojis=[]
			try {
				let response = await this.$axios.get(`${this.path}//${this.messageId}`);
				this.message = response.data;
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
			console.log(response.data.Content)
			return response.data.Content
		}
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<button class = "message" @click="showMessageModal = true">
		{{message?.Sender}}<br>
		{{message?.Content}}<br>
		{{message?.Timestamp}}<br>
		{{message?.Status}}<br>
		<div v-for="emoji in emojis.slice(-3)">
			{{emoji.Content}}
		</div>
	</button>
	<MessageModal
		:showModal="showMessageModal"
		:path="`${this.path}//${this.messageId}`"
		:emojis="emojis"
		@close="showMessageModal = false"
		@save="save"
	/>
</template>

<style scoped>

</style>
