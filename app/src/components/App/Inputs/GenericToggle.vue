<script lang="ts" setup>

const props = defineProps<{
	required?: boolean;
	disabled?: boolean;
	label?: string;
}>();

const model = defineModel<boolean>();

</script>

<template>
	<div class="toggle-label" :class="{ disabled }" @click="model = !model">

		<div class="toggle" :class="{ active: model }"></div>

		<span v-if="label || $slots.default" class="label-text">
			<slot>
				{{ label }}
			</slot>
		</span>

	</div>
</template>

<style lang="scss" scoped>

.toggle-label {
	display: flex;
	flex-flow: row nowrap;
	gap: 0.75rem;
	align-items: center;
	padding: 0.25rem;
	cursor: pointer;

	.toggle {
		width: 2.85rem;
		height: 1.75rem;
		background-color: var(--app-theme-powder-trail);
		border-radius: 0.875rem;
		border: none;
		position: relative;
		padding: 0.25rem;
		transition: all 200ms ease;

		&::after {
			content: "";
			display: block;
			width: 1.3rem;
			height: 1.3rem;
			background-color: var(--app-theme-snow-white);
			border-radius: 50%;
			transition: transform 0.25s cubic-bezier(0.4, 0, 0.2, 1);
		}

		&.active {
			background-color: var(--app-theme-sky-blue);

			&::after {
				transform: translateX(1.05rem);
			}
		}
	}

	.label-text {
		font-size: 0.8rem;
		line-height: 1.5em;
		user-select: none;
		font-weight: 500;
	}

	&.disabled {
		opacity: 0.8;
		filter: saturate(0);
		pointer-events: none;
	}
}

</style>
