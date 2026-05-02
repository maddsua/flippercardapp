<script setup lang="ts">
import { ref } from 'vue';
import type { PollOption } from './content';

const props = defineProps<{
	entry: PollOption;
	is_quiz?: boolean;
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
	<button type="button" class="card-poll-option" :class="{ wrong, right }" @click.stop="handleSelect">
		{{ entry.value }}
	</button>
</template>

<style lang="scss" scoped>
	.card-poll-option {
		display: block;
		width: 100%;
		max-width: 25rem;
		font-weight: 600;
		font-size: 1.125em;
		color: var(--app-color-white);
		background-color: var(--app-accent-blue);
		transition: all 150ms ease;

		&:hover {
			cursor: pointer;
			transform: scale(1.025);
		}

		//	todo: add style overrides

		&.wrong {
			pointer-events: none;
			animation: horizontal-shaking 200ms 2;
			color: var(--app-color-white);
			background-color: var(--app-accent-red);
			cursor: not-allowed;
		}

		&.right {
			color: var(--app-color-white);
			background-color: var(--app-accent-green);
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
