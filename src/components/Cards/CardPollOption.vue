<script setup lang="ts">
import { ref } from 'vue';
import type { PollOption } from './content';

const props = defineProps<{
	entry: PollOption
}>();

const emit = defineEmits<{
	(e: 'select'): void;
}>();

const wrong = ref(false);

const handleSelect = () => {

	const { action } = props.entry;

	if (action === 'fail' || action === 'fail-show') {
		wrong.value = true;
	}

	emit('select');
};


//	todo: add additional styling

</script>

<template>
	<button type="button" class="card-poll-option" :class="{ wrong }" @click.stop="handleSelect">
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
		color: white;
		background-color: blue;
		transition: all 150ms ease;

		&:hover {
			cursor: pointer;
			transform: scale(1.025);
		}

		&.wrong {
			pointer-events: none;
			animation: horizontal-shaking 200ms 2;
			background-color: red;
			cursor: not-allowed;
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
