<script setup lang="ts">
import { computed, reactive, type CSSProperties } from 'vue';
import type { CardContentElementTheme, CardPollElementOptionNode } from '../../content';

const props = defineProps<{
	entry: CardPollElementOptionNode;
	is_quiz?: boolean;
	theme?: CardContentElementTheme;
}>();

const emit = defineEmits<{
	(e: 'select'): void;
}>();

const state = reactive({
	wrong: false,
	right: false,
	selected: false,
});

const handleSelect = () => {

	if (state.selected) {
		return;
	}
	state.selected = true;

	if (props.is_quiz) {
		if (!props.entry.is_answer) {
			state.wrong = true;
		} else {
			state.right = true;
		}
	}

	emit('select');
};

const applyStyles = computed((): CSSProperties => ({
	backgroundColor: props.theme?.fill_color,
	color: props.theme?.mask_color,
}));

</script>

<template>
	<button type="button" :class="state" :style="applyStyles" @click.stop="handleSelect">
		{{ entry.value }}
	</button>
</template>

<style lang="scss" scoped>

	@use '../../media.scss';

	button {
		display: block;
		width: 100%;
		font-weight: 600;
		font-size: 1em;
		font-weight: 500;
		border-radius: 0.5em;
		padding: 0.6em 1.2em;
		color: var(--app-theme-snow-white);
		background-color: var(--app-theme-sky-blue);
		outline: none;
		border: none;
		transition: color, background-color, transform 150ms ease;

		@include media.non-sticky-hover {
			transform: scale(1.025);
			background-color: var(--app-theme-deep-lavender);
		}

		&.wrong {
			pointer-events: none;
			cursor: not-allowed;
			animation: horizontal-shaking 200ms 2;
			color: var(--app-theme-snow-white) !important;
			background-color: var(--app-theme-blood-red) !important;
		}

		&.right {
			color: var(--app-theme-snow-white) !important;
			background-color: var(--app-theme-irish-green) !important;
		}

		&.selected {
			color: var(--app-theme-snow-white) !important;
			background-color: var(--app-theme-spooky-orange) !important;
		}
	}

	@keyframes horizontal-shaking {
		0% { transform: translateX(0) }
		25% { transform: translateX(5px) }
		50% { transform: translateX(-5px) }
		75% { transform: translateX(5px) }
		100% { transform: translateX(0) }
	}
</style>
