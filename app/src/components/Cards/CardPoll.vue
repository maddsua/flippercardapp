<script setup lang="ts">
import { computed, reactive } from 'vue';
import CardPollOption from './CardPollOption.vue';
import type { CardContentElementTheme, CardPollElement, CardPollElementOptionNode } from '../../content';
import { shuffleArray } from '../../arrays';

const props = defineProps<{
	entry: CardPollElement;
	theme?: CardContentElementTheme;
}>();

const emit = defineEmits<{
	(e: 'flip'): void;
	(e: 'next'): void;
	(e: 'score', score: number): void;
}>();

const state = reactive({
	answered: false,
	givenAnswers: new Set<CardPollElementOptionNode>(),
	wrongAnswers: 0,
});

const selectAnswer = (answer: CardPollElementOptionNode) => {

	if (state.givenAnswers.has(answer)) {
		return;
	}

	state.givenAnswers.add(answer);

	if (!props.entry.is_quiz) {

		if (!state.answered) {
			setTimeout(() => emit('next'), 250);
		}
		state.answered = true;

		return;
	}

	if (answer.is_answer) {

		if (!state.answered) {

			if (state.wrongAnswers === 0) {
				emit('score', 1);
				setTimeout(() => emit('next'), 300);
			} else {
				setTimeout(() => emit('flip'), 500);
			}

			state.answered = true;
		}

		return
	}

	state.wrongAnswers++;
	emit('score', 0);
};

const options = computed(() => {
	const entries = [...props.entry.content];
	shuffleArray(entries);
	return entries;
});

</script>

<template>
	<div class="card-poll" data-interactive="" :disabled="state.answered">

		<CardPollOption v-if="options.length" v-for="option of options"
			:entry="option"
			:is_quiz="props.entry.is_quiz"
			:theme="theme"
			@select="selectAnswer(option)" />

		<template v-else>
			[Poll options]
		</template>

	</div>
</template>

<style lang="scss" scoped>
	.card-poll {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.5em;
		width: 100%;
		margin-top: auto;

		&:disabled {
			pointer-events: none;
			filter: saturate(0.25);
		}
	}
</style>
