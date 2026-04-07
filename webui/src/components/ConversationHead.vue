<script >
export default{
	name: 'ConversationHead',
	props: {
		path: String,
		conversation: Object,
},
	methods:{
		async openConversation(conversationId){
			this.$router.push(`${this.path}/${conversationId}`)
		},
	}
}
</script>

<template>
	<button class = "conversation" @click="openConversation(conversation.Id.Id)">
		{{conversation.Name}}<br>
		{{conversation.Picture}}<br>
		{{conversation.Snippet}}<br>
		{{conversation.Date}}<br>
	</button>
</template>

<style scoped>

</style>
