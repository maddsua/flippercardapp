<script setup lang="ts">
const props = defineProps<{
	theme?: 'blue' | 'orange' | 'green' | 'red';
	variant?: 'thin' | 'thick' | 'wide';
	spinner?: boolean;
	disabled?: boolean;
}>();
</script>

<template>
	<button type="button" :class="{ [`variant-${variant}`]: !!variant, [`theme-${theme}`]: !!theme, spinner }" :disabled="disabled">
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
		width: fit-content;
		flex-shrink: 0;

		&:hover {
			cursor: pointer;
			background-color: var(--app-theme-deep-lavender);
		}

		&.theme-orange {
			background-color: var(--app-theme-spooky-orange);
		}

		&.theme-green {
			background-color: var(--app-theme-irish-green);
		}

		&.theme-red {
			background-color: var(--app-theme-blood-red);
		}

		&.variant-thick {
			padding: 1rem 2rem;
			font-size: 0.85rem;
		}

		&.variant-thin {
			padding: 0.5rem 0.75rem;
			font-size: 0.65rem;
			border-radius: 0.5rem;
		}

		&.variant-wide {
			width: 100%;
		}

		&:disabled {
			pointer-events: none;
			filter: saturate(0);
			opacity: 0.75;
		}

		&.spinner {
			color: transparent;
			background-image: url(/src/assets/animated/ring-resize-white.svg);
			background-position: center;
			background-repeat: no-repeat;
			background-size: 1.5rem;
		}
	}
</style>
