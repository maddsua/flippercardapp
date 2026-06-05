<script setup lang="ts">
import { computed, ref, watch } from 'vue';

type Value = string | null;

interface Option {
	value: Value;
	label?: string;
}

const props = defineProps<{
	placeholder?: string;
	options: Option[];
	disabled?: boolean;
}>();

const model = defineModel<Value>();

const activeOption = computed(() => props.options.find((item) => item.value === model.value));
const isOpen = ref(false);

const toggle = () => {
	if (isOpen.value) {
		isOpen.value = false;
		return;
	}

	if (!props.options.length) {
		return;
	}
	isOpen.value = true;
};

const select = (val: Value) => {
	isOpen.value = false;
	model.value = val;
};

watch(model, () => (isOpen.value = false));
</script>

<template>
	<div class="dropdown" :disabled="disabled">
		<div class="header" :class="{ disabled: !$props.options.length, open: isOpen }" @click="toggle">
			<span>
				{{ activeOption?.label || activeOption?.value || $props.placeholder || '-' }}
			</span>
		</div>
		<div v-if="isOpen" class="options">
			<div class="list">
				<button type="button" v-for="opt of $props.options" @click="select(opt.value)">
					{{ opt.label || opt.value }}
				</button>
			</div>
		</div>
	</div>
</template>

<style lang="scss" scoped>
.dropdown {
	position: relative;
	z-index: 1;
	width: 100%;

	.header {
		display: flex;
		flex-flow: row nowrap;
		align-items: center;
		gap: 1rem;
		justify-content: space-between;
		padding: 0.65rem 1rem;
		border-radius: 0.5rem;
		background-color: var(--app-theme-ghostly-glow);

		span {
			display: block;
			color: var(--app-theme-snow-white);
			font-weight: 600;
			font-size: 0.75rem;
			min-width: 0;
			overflow: hidden;
			text-overflow: ellipsis;
			white-space: nowrap;
		}

		&::after {
			content: '';
			display: block;
			width: 0.35rem;
			height: 0.35rem;
			border-left: 2px solid var(--app-theme-snow-white);
			border-bottom: 2px solid var(--app-theme-snow-white);
			transform: rotate(-45deg);
			transition: transform 150ms ease;
		}

		&.open::after {
			transform: rotate(135deg);
		}

		&.disabled {
			pointer-events: none;
			opacity: 0.5rem;
		}

		&:hover {
			cursor: pointer;
			background-color: var(--app-theme-midnight-glow);
		}
	}

	.options {
		position: absolute;
		bottom: -0.5rem;
		left: 0;
		width: 100%;

		.list {
			position: absolute;
			top: 0;
			left: 0;
			width: 100%;
			max-height: 15rem;
			overflow: hidden auto;
			scrollbar-width: thin;
			background-color: var(--app-theme-carbon);
			border: 1px solid var(--app-theme-powder-trail);
			border-radius: 0.5rem;
		}

		button {
			display: block;
			width: 100%;
			border: none;
			outline: none;
			background: none;
			padding: 0.75rem 1rem;
			text-align: left;
			font-size: 0.75rem;
			font-weight: 400;
			color: var(--app-theme-snow-white);

			&:hover {
				background-color: var(--app-theme-midnight-glow);
				cursor: pointer;
			}
		}
	}

	&:disabled {
		pointer-events: none;
		filter: saturate(0);
		opacity: 0.75;
	}
}
</style>
