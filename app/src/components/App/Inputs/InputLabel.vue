<script setup lang="ts">
const props = defineProps<{
	variant?: 'slick';
}>();
</script>

<template>
	<div class="input-label" :class="{ [`variant-${variant}`]: !!variant }">
		<div class="label">
			<slot name="label">
				[Label]
			</slot>
		</div>
		<div class="input-container">
			<slot>
				[Input]
			</slot>
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.input-label {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;

		.label {
			font-size: 0.85rem;
			color: var(--app-theme-mysterious-white);
			text-transform: uppercase;
		}

		.input-container {
			display: flex;
			flex-flow: row nowrap;
			align-items: center;
		}

		&.variant-slick {
			gap: 0.35rem;

			.label {
				font-size: 0.65rem;
			}
		}
	}
</style>
