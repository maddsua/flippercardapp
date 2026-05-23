<script setup lang="ts">
import { ref } from 'vue';
import type { CardContentElementTheme, CardPollElementOptionNode } from '../../content';

const props = defineProps<{
	entry: CardPollElementOptionNode;
	is_quiz?: boolean;
	theme?: CardContentElementTheme;
}>();

const emit = defineEmits<{
	(e: 'select'): void;
}>();

const wrong = ref(false);
const right = ref(false);

const handleSelect = () => {

	if (props.is_quiz) {
		if (!props.entry.is_answer) {
			wrong.value = true;
		} else {
			right.value = true;
		}
	}

	emit('select');
};

</script>

<template>
	<button type="button" :class="{ wrong, right }" :style="{ backgroundColor: theme?.fill_color, color: theme?.mask_color }" @click.stop="handleSelect">
		{{ entry.value }}
	</button>
</template>

<style lang="scss" scoped>
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

		&:hover {
			cursor: pointer;
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
	}

	@keyframes horizontal-shaking {
		0% { transform: translateX(0) }
		25% { transform: translateX(5px) }
		50% { transform: translateX(-5px) }
		75% { transform: translateX(5px) }
		100% { transform: translateX(0) }
	}
</style>
