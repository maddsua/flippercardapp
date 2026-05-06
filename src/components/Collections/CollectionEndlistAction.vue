<script setup lang="ts">
const emit = defineEmits<{
	(e: 'action'): void;
}>();
</script>

<template>
	<div class="collection-action">
		<button type="button" @click="emit('action')">
			<slot>
				[Label]
			</slot>
		</button>
	</div>
</template>

<style lang="scss" scoped>

	.collection-action {
		display: flex;
		flex-direction: column;
		align-items: center;
	}

	button {
		display: block;
		text-align: center;
		padding: 0.75rem 1rem;
		font-size: 0.75rem;
		font-weight: 600;
		border-radius: 0.75rem;
		color: var(--app-theme-snow-white);
		background-color: var(--app-theme-sky-blue);
		border: none;
		outline: none;
		transition: all 200ms ease;

		&:hover {
			cursor: pointer;
			background-color: var(--app-theme-deep-lavender);
		}
	}
</style>
