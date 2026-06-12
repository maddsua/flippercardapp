<script setup lang="ts">
const props = defineProps<{
	title: string;
}>();
</script>

<template>
	<div class="editor-group">
		<div class="title">
			{{ title }}
		</div>
		<div class="content">
			<slot>[Content]</slot>
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.editor-group {
		display: flex;
		flex-direction: column;
		gap: 1.75rem;

		.title {
			font-size: 0.95rem;
			font-weight: 300;
		}

		.content {
			display: flex;
			flex-direction: column;
			gap: 1.5rem;
		}
	}
</style>
