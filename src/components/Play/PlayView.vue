<script setup lang="ts">
import { computed, onMounted, reactive } from 'vue';
import type { CardDeck, CardNode } from '../Cards/content';
import CardView from '../Cards/CardView.vue';
import EndscreenView from '../Endscreen/EndscreenView.vue';
import { shuffleArray } from '../../shuffle';

const { deck } = defineProps<{
	deck: CardDeck;
}>();

interface gameState {
	startTime: Date;
	isFinished: boolean;
	questions: number;
	totalScore: number;
	playTime: number;
}

const state = reactive({
	data: {
		cards: [] as CardNode[],
		labels: [] as string[],
		busy: false,
		ready: false,
		error: null as string | null,
	},
	game: null as gameState | null,
});

const endgameStats = computed(() => state.game?.isFinished ? ({
	questions: state.game.questions,
	score: state.game.totalScore,
	time: state.game.playTime,
}) : null);

const cards = computed(() => {

	if (!state.data.ready) {
		return null;
	}

	const entries = [...state.data.cards];
	shuffleArray(entries);
	return entries;
});

const unwrapData = async () => {

	const { data, error } = await deck.cards();
	if (!data) {
		state.data.error = error?.message || 'unable to load cards';
		return
	} else if (data.length === 0) {
		state.data.error = 'Empty deck';
		return;
	}

	state.data.cards = data;
	state.data.ready = true;

	state.data.labels = [
		await deck.collection()
			.then(data => data.data?.name)
			.catch(() => null) || '',
		deck.name,
	].filter(item => item.length);
};

const loadData = async () => {

	if (state.data.busy) {
		return;
	}

	state.data.busy = true;

	state.data.error = null;
	state.data.ready = false;
	state.data.cards = [];

	await unwrapData()

	state.data.busy = false;
};

const initGame = () => {

	const questions = state.data.cards
		.map(item => [item.back, item.front])
		.flat()
		.map(item => item.content.some(item => item.type === 'poll')).length;

	state.game = {
		questions,
		totalScore: 0,
		playTime: 0,
		isFinished: false,
		startTime: new Date(),
	};
};

const updateGameScore = (delta: number) => {
	if (!state.game) {
		return;
	}
	state.game.totalScore += delta;
};

onMounted(async () => {
	await loadData();
	initGame();
});

const finishDeck = () => {

	if (!state.game) {
		return;
	}

	state.game.playTime = Math.floor((new Date().getTime() - state.game.startTime.getTime()) / 1000);
	state.game.isFinished = true;
};

const handleNavigate = (target: 'home' | 'try_again') => {
	switch (target) {
		case 'home':
			//	todo: implement
			alert('not implemented');
			break;
		case 'try_again':
			initGame();
			break;
	}
};

</script>

<template>

	<template v-if="cards?.length">
		<CardView v-if="!endgameStats" :labels="state.data.labels" :entries="cards" @score="updateGameScore" @end="finishDeck" />
		<EndscreenView v-else :stats="endgameStats" @navigate="handleNavigate" />
	</template>

	<div class="view-message" v-else>

		<template v-if="state.data.error">
			Error: {{ state.data.error }}
		</template>

		<template v-else>
			Loading cards...
		</template>

	</div>

</template>

<style lang="scss" scoped>
	.view-message {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		width: 100%;
		height: 100%;

		font-size: 0.85rem;
	}
</style>
