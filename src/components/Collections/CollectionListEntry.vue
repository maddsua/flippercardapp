<script setup lang="ts">
const props = defineProps<{
	title: string;
	//	todo: add starring
	//	todo: add stats
}>();
</script>

<template>
	<button type="button">
		{{ title }}
	</button>
</template>

<style lang="scss" scoped>
	button {
		display: flex;
		flex-flow: row nowrap;
		gap: 1rem;
		align-items: center;
		font-size: 1rem;
		font-weight: 600;
		padding: 1rem 2rem;
		border-radius: 1rem;
		border: none;
		outline: none;
		transition: all 150ms ease;

		background-color: var(--app-theme-irish-green);
		color: var(--app-theme-snow-white);

		&:hover {
			cursor: pointer;
			background-color: var(--app-theme-spooky-orange);
		}
	}
</style>
