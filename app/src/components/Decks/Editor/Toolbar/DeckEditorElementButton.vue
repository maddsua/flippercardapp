<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
	label: string;
	disabled?: boolean;
	icon?: 'add-title' | 'add-text' | 'add-image' | 'add-poll';
}>();

const applyClass = computed(() => ({
	disabled: props.disabled,
	[`icon-${props.icon}`]: !!props.icon,
}));

</script>

<template>
	<div class="element-button" :class="applyClass" :title="label"></div>
</template>

<style lang="scss" scoped>
	.element-button {
		position: relative;
		width: 1.65rem;
		height: 1.65rem;
		border-radius: 0.25rem;
		background-color: var(--app-theme-ghostly-glow);
		display: flex;
		flex-direction: column;
		justify-content: end;
		background-size: 1rem;
		background-repeat: no-repeat;
		background-position: center;

		&.icon-add-title {
			background-image: url(/src/assets/icons/add-title-color.svg);
		}

		&.icon-add-text {
			background-image: url(/src/assets/icons/add-text-color.svg);
		}

		&.icon-add-image {
			background-image: url(/src/assets/icons/add-image-color.svg);
		}

		&.icon-add-poll {
			background-image: url(/src/assets/icons/add-poll-color.svg);
		}

		&.disabled {
			pointer-events: none;
			opacity: 0.5;
			filter: saturate(0);
		}

		.color-fill {
			width: 100%;
			height: 0.25rem;
			border-radius: 0.25rem;
		}

		&:hover, &.active {
			cursor: pointer;
			background-color: var(--app-theme-powder-trail);
		}
	}
</style>
