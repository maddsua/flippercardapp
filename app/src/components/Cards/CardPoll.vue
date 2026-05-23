<script setup lang="ts">
import { computed } from 'vue';
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

let wrongQuizTakes = 0;

const handleOptionSelect = (opt: CardPollElementOptionNode) => {

	if (!props.entry.is_quiz) {
		emit('next');
		return;
	}

	if (!opt.is_answer) {
		wrongQuizTakes++;
		emit('score', 0);
	}

	if (opt.is_answer) {

		emit('score', 1);

		if (wrongQuizTakes === 0) {
			setTimeout(() => emit('next'), 300);
		} else {
			setTimeout(() => emit('flip'), 500);
		}

		return
	}
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
			:is_quiz="props.entry.is_quiz"
			:theme="theme"
			@select="handleOptionSelect(option)" />
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
