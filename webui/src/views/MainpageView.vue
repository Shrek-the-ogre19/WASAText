<script>
import CreateConversation from "../components/CreateConversation.vue";

export default {
	components: {
		CreateConversation
	},
	data: function() {
		return {
			errormsg: null,
			loading: false,
			path: this.$route.path,
			conversations: null,
			showModal: false,
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
		async openProfile(){
			let pathParts = this.path.split('/');
			pathParts.pop();
			let newPath = pathParts.join('/');
			await this.$router.push(newPath);
		},
		async openUsersList(){
			let pathParts = this.path.split('/');
			pathParts.pop();
			let newPath = pathParts.join('/');
			await this.$router.push(newPath + "/users");
		},
		async openConversation(conversationId){
			this.$router.push(`/${conversationId}`)
		},
		openModal() {
			this.showModal = true;
		},
		closeModal() {
			this.showModal = false;
		},
		handleModal() {
			console.log('Confirmed:', this.inputValue);
			this.closeModal();
		}
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>

</template>

<style scoped>

</style>
