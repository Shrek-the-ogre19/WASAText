<script >
export default{
	name: 'ConversationHead',
	props: {
		path: String,
		conversation: Object,
},data: function() {
		return {
			errormsg: null,
			loading: false,
			text: null,
			image: null,
		}
	},
	methods:{
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			let str = this.conversation.Snippet;
			let parts = str.split('data:image');
			this.text = parts[0]
			if (parts[1]){
				this.image = 'data:image' + parts[1]}
			this.loading = false;
		},
		async openConversation(conversationId){
			this.$router.push(`${this.path}/${conversationId}`)
		},
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<button class = "conversation" @click="openConversation(conversation.Id.Id)">
		<img :src="conversation.Picture" class="img" alt="conversationPicture"> <br>
		{{conversation.Name}}<br>
		{{text}}<br>
		<div v-if="image">
			img
		</div>
		{{conversation.Date}}<br>
	</button>
</template>

<style scoped>

</style>
