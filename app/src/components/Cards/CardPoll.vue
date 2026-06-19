<script setup lang="ts">
import { computed, reactive } from 'vue';
import CardPollOption from './CardPollOption.vue';
import type { CardContentElementTheme, CardPollNode, CardPollNodeOption } from '../../content';
import { shuffleArray } from '../../arrays';

const props = defineProps<{
	entry: CardPollNode;
	theme?: CardContentElementTheme | null;
}>();

const emit = defineEmits<{
	(e: 'score', score: number, final?: boolean): void;
}>();

const state = reactive({
	answered: false,
	givenAnswers: new Set<CardPollNodeOption>(),
	wrongAnswers: 0,
});

const selectAnswer = (answer: CardPollNodeOption) => {

	if (state.givenAnswers.has(answer) || state.answered) {
		return;
	}

	state.givenAnswers.add(answer);

	// not really used for anything tbh
	if (!props.entry.is_quiz) {
		state.answered = true;
		return;
	}

	if (answer.is_answer) {
		state.answered = true;
		emit('score', state.wrongAnswers === 0 ? 1 : 0, true);
		return;
	}

	state.wrongAnswers++;
	emit('score', 0, false);
};

const options = computed(() => {
	const entries = [...props.entry.content];
	shuffleArray(entries);
	return entries;
});

</script>

<template>
	<div class="card-poll" data-interactive="">

		<CardPollOption v-if="options.length" v-for="option of options"
			:entry="option"
			:is_quiz="!!props.entry.is_quiz"
			:theme="theme"
			:disabled="state.answered"
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
	}
</style>
