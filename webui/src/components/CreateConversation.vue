<script>
export default {
	data: function(){
		return{
			receivers : null,
	}

	},
	props: {
		isVisible: {
			type: Boolean,
			default: false
		},
		title: {
			type: CreateConversation,
			default: 'Modal'
		}
	},
	methods: {
		close() {
			this.$emit('close');
		},
		confirm() {
			this.$emit('confirm');
		}
	}
}
</script>

<template>

</template>

<style scoped>

</style>
