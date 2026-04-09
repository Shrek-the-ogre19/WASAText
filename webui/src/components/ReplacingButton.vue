<script>
export default {
	name: 'ReplacingButton',
	props: {
		itemName: String,
		item: String
	},
	computed: {
		text() {
			return `Enter ${this.itemName}`;
		}},
	data() {
		return {
			showModal: false,
			newItem: "",
		}
	},
	methods:{
		close() {
			this.$emit('close');
		},

		save() {
			if (this.newItem && this.newItem.trim()) {
				this.$emit('save', this.newItem);
				this.showModal = false;
				this.newItem = '';
			}},
	},
	emits: ['close', 'save'],
}
</script>

<template>
	<p>
		<button class="mainButton" @click="showModal = true">{{item}}</button>
	</p>

	<div v-if="showModal" class="modal">
		<div class="modal-content">
			<span class="close" @click="showModal = false">&times;</span>
			<input v-model="newItem" :placeholder="text" />
			<button class="button" @click="save">CONFIRM</button>
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
</style>
