<script setup lang="ts">
import type { CardFaceTheme } from './content';

const props = defineProps<{
	theme?: CardFaceTheme;
}>();

</script>

<template>
	<div class="card-canvas" :style="{ color: theme?.mask_color }">
		<div class="card-decoration"></div>
		<div class="card-content" :style="{ backgroundColor: theme?.fill_color, borderColor: theme?.fill_color || theme?.outline_color, color: theme?.mask_color }">
			<slot>
				[Card content]
			</slot>
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.card-canvas {
		position: absolute;
		top: 0;
		bottom: 0;
		overflow: hidden;
		display: flex;
		flex-direction: column;
		width: 100%;
		height: 100%;
		color: var(--app-theme-spooky-orange);
		background-color: var(--app-theme-snow-white);
		border-radius: 2rem;
		padding: 1rem;
		
		-webkit-backface-visibility: hidden;
		backface-visibility: hidden;

		// these two somehow prevent jagged edges
		outline: 1px solid transparent;
		box-shadow: 1rem 0.25rem 0.25rem rgba(0, 0, 0, 0.25);

		&:nth-child(2n) {
			transform: rotateY(180deg);
		}

		//	todo: connect
		&:not(:nth-child(2n))::before {
			content: "?";
			display: block;
			font-size: 3.5rem;
			line-height: 1rem;
			position: absolute;
			left: 2.65rem;
			top: 3.25rem;
			font-weight: 600;
			color: inherit;
		}
	}

	.card-content {
		display: flex;
		flex-direction: column;
		gap: 2rem;
		align-items: center;
		flex-grow: 1;
		border: 3px solid var(--app-theme-spooky-orange);
		border-radius: 1.5rem;
		padding: 2rem;
		color: var(--app-theme-midnight);
	}
</style>
