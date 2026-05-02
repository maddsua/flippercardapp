<script setup lang="ts">
import CardPollOption from './CardPollOption.vue';
import type { PollNode, PollOption } from './content';

const props = defineProps<{
	entry: PollNode;
}>();

const emit = defineEmits<{
	(e: 'flip'): void;
	(e: 'next'): void;
}>();

let wrongQuizTakes = 0;

const handleOptionSelect = (opt: PollOption) => {

	if (!props.entry.is_quiz) {
		emit('next');
		return;
	}

	if (!opt.is_answer) {
		wrongQuizTakes++;
	}

	if (opt.is_answer) {

		if (wrongQuizTakes === 0) {
			setTimeout(() => emit('next'), 300);
		} else {
			setTimeout(() => emit('flip'), 500);
		}

		return
	}
};

</script>

<template>
	<div class="card-poll">
		<CardPollOption v-for="option of entry.content" :entry="option" :is_quiz="props.entry.is_quiz" @select="handleOptionSelect(option)" />
	</div>
</template>

<style lang="scss" scoped>
	.card-poll {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.5rem;
		width: 100%;

		&:last-child {
			margin-top: auto;
		}
	}
</style>
