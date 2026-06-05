<script setup lang="ts">
import { computed, ref, type CSSProperties } from 'vue';
import type { CardContentElementTheme, CardPollNodeOption } from '../../content';

const props = defineProps<{
	entry: CardPollNodeOption;
	is_quiz?: boolean;
	theme?: CardContentElementTheme;
}>();

const emit = defineEmits<{
	(e: 'select'): void;
}>();

const answered = ref(false);

const handleSelect = () => {

	if (answered.value) {
		return;
	}
	answered.value = true;

	emit('select');
};

const applyClasses = computed((): Record<string, boolean> => ({
	right: answered.value && props.is_quiz && !!props.entry.is_answer,
	wrong: answered.value && props.is_quiz && !props.entry.is_answer,
	selected: answered.value && !props.is_quiz,
	answered: answered.value,
}));

const applyStyles = computed((): CSSProperties => ({
	backgroundColor: props.theme?.fill_color,
	color: props.theme?.mask_color,
}));

</script>

<template>
	<button type="button" :class="applyClasses" :style="applyStyles" @click.stop="handleSelect">
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

		@include media.phone {
			padding: 0.75em 1.25em;
		}

		@include media.non-sticky-hover {
			transform: scale(1.025);
			background-color: var(--app-theme-deep-lavender);
		}

		&:disabled {
			pointer-events: none;
			filter: saturate(0);
		}

		&.answered {
			pointer-events: none;
			cursor: default;
			filter: unset;
		}

		&.wrong {
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
