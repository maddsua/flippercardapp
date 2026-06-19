<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue';
import CardColorSwatch from '../CardColorSwatch.vue';
import { isContainerOutsideClick } from '@/dom';

const props = defineProps<{
	label: string;
	disabled?: boolean;
	icon?: 'fill' | 'fill-interactive' | 'text' | 'text-interactive' | 'outline';
}>();

const model = defineModel<string | null>();

const trayOpen = ref(false);
const containerRef = ref<HTMLElement | null>(null);

const handleOutsideClicks = (event: Event) => {
	if (isContainerOutsideClick(containerRef.value, event.target)) {
		trayOpen.value = false;
	}
};

onMounted(() => window.addEventListener('click', handleOutsideClicks));
onUnmounted(() => window.removeEventListener('click', handleOutsideClicks));

const applyClass = computed(() => ({
	active: trayOpen.value,
	disabled: props.disabled,
	[`icon-${props.icon}`]: !!props.icon,
}));

</script>

<template>
	<div class="element-dropdown" :class="applyClass" :title="label" @click="trayOpen = !trayOpen" ref="containerRef">

		<div v-if="model" class="color-fill"  :style="{ backgroundColor: model }"></div>

		<div v-if="trayOpen" class="dropdown-anchor">
			<div class="dropdown-tray">
				<CardColorSwatch size="small" :label="label" v-model="model" />
			</div>
		</div>

	</div>
</template>

<style lang="scss" scoped>
	.element-dropdown {
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
		background-position: center 0.25rem;

		&.icon-text {
			background-image: url(/src/assets/icons/text-mask.svg);
		}

		&.icon-text-interactive {
			background-image: url(/src/assets/icons/text-interactive-mask.svg);
		}

		&.icon-fill {
			background-image: url(/src/assets/icons/fill-mask.svg);
		}

		&.icon-fill-interactive {
			background-image: url(/src/assets/icons/fill-interactive-mask.svg);
		}

		&.icon-outline {
			background-image: url(/src/assets/icons/border-mask.svg);
		}

		&.disabled {
			pointer-events: none;
			opacity: 0.5;
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

		.dropdown-anchor {
			position: absolute;
			bottom: 0;
			left: 0;
		}

		.dropdown-tray {
			position: absolute;
			top: 0;
			left: 0;
			padding: 0.5rem;
			border-radius: 0.5rem;
			overflow: hidden;
			background-color: var(--app-theme-carbon);
		}
	}
</style>
