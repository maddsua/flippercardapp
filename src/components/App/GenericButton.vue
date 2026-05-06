<script setup lang="ts">
const props = defineProps<{
	theme?: 'blue' | 'orange';
	variant?: 'thick' | 'expanded';
}>();
</script>

<template>
	<button type="button" :class="[theme, variant]">
		<slot>
			[Button label]
		</slot>
	</button>
</template>

<style lang="scss" scoped>
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

		&.orange {
			background-color: var(--app-theme-spooky-orange);
		}

		&.thick {
			padding: 1rem 2rem;
			font-size: 0.85rem;
		}

		&.wide {
			flex-grow: 1;
		}
	}
</style>
