<script setup lang="ts">
const props = defineProps<{
	turned?: boolean;
}>();
</script>

<template>
	<div class="card-preview-slot" :class="{ turned }">
		<div class="card-face-container">
			<slot>
				[Card face]
			</slot>
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.card-preview-slot {
		display: flex;
		flex-grow: 1;
		align-items: center;
		justify-content: center;
		transition: all 150ms ease;

		&.turned {
			transform: rotateY(90deg);
		}

		.card-face-container {
			position: relative;
			width: 70%;
			container-type: inline-size;
			aspect-ratio: 0.6667;
		}
	}
</style>
