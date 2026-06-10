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
			if (parts[0]){this.text = parts[0]}

			if (parts[1]){
				this.image = 'data:image' + parts[1]}
			this.loading = false;
		},
		async openConversation(conversationId){
			this.$router.push(`${this.path}/${conversationId}`)
		},
	},
	watch: {
		conversation: {
			deep: true,
			handler() {
				this.refresh();
			},
		},
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<button class = "conversation" @click="openConversation(conversation.Id.Id)">
		<div v-if="conversation.Picture == 'default'">
			<img src="/9572728.png" class="img" alt="conversationPicture"/>
		</div>
		<div v-else>
		<img :src="conversation.Picture" class="img" alt="conversationPicture">
		</div>
		<div class="name">{{conversation.Name}}</div>
		<div class="snippet">{{text}}
		<div v-if="image">
			<img :src="image"  class="snippetPicture" alt="snippetPicture">
		</div>

		</div>
		<div class="date">
		{{conversation.Date}}
		</div>
	</button>
</template>

<style scoped>
.conversation {
	width: 100%;
	padding: 16px 20px;
	background-color: lightseagreen;
	border: none;
	border-bottom: 1px solid #e9ecef;
	cursor: pointer;
	transition: all 0.2s ease;
	text-align: left;
	display: flex;
	align-items: center;
	gap: 12px;
	position: relative;
	flex-wrap: wrap;
}

.conversation:hover {
	background-color: #f8f9fa;
}

.img {
	width: 52px;
	height: 52px;
	border-radius: 50%;
	object-fit: cover;
	border: 2px solid #ffffff;
	box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
	flex-shrink: 0;
}

.name {
	font-size: 20px;
	font-weight: 600;
	color: #212529;
	margin-bottom: 4px;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
	flex: 1;
}

.snippet {
	font-size: 25px;
	color: black;
	display: flex;
	align-items: center;
	gap: 8px;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
	max-width: 100%;
	margin-top: 4px;
	padding-right: 60px;
}

.snippetPicture {
	width: 24px;
	height: 24px;
	border-radius: 50%;
	object-fit: cover;
	border: 2px solid #ffffff;
	box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
	flex-shrink: 0;
}



.date {
	font-size: 11px;
	color: black;
	position: absolute;
	top: 16px;
	right: 20px;
	flex-shrink: 0;
}
</style>
