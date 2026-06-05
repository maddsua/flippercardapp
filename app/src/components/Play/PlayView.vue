<script setup lang="ts">
import { computed, onMounted, reactive } from 'vue';
import type { CardNode } from '../../content';
import CardWidget from '../Cards/CardWidget.vue';
import Endscreen from '../Endscreen/Endscreen.vue';
import FullscreenMessage from '../App/Messages/FullscreenMessage.vue';
import { shuffleArray } from '../../arrays';
import { useRoute, useRouter } from 'vue-router';
import LoadingMessage from '../App/Messages/LoadingMessage.vue';
import ErrorMessage from '../App/Messages/ErrorMessage.vue';
import GenericButton from '../App/Inputs/GenericButton.vue';
import { useClient } from '../../api';
import { useStorage } from '../../storage/storage';

const router = useRouter();
const route = useRoute();
const client = useClient();
const store = useStorage();

interface RoundState {
	startTime: Date;
	isFinished: boolean;
	questions: number;
	totalScore: number;
	playTime: number;
};

const state = reactive({
	cards: null as CardNode[] | null,
	labels: [] as string[],
	collectionID: null as string | null,
	deckID: null as string | null,
	isMarked: false,
	error: null as string | null,
	round: null as RoundState | null,
	showNavigation: false,
});

const statsScreen = computed(() => state.round?.isFinished ? ({
	questions: state.round.questions,
	score: state.round.totalScore,
	time: state.round.playTime,
}) : null);

const cards = computed(() => {

	if (!state.cards) {
		return null;
	}

	const entries = [...state.cards];
	shuffleArray(entries);
	return entries;
});

const initRound = () => {

	if (!state.cards) {
		return null;
	}

	state.round = {
		questions: state.cards.map(item => [item.front, item.back])
			.flat()
			.map(item => item.content)
			.flat()
			.filter(item => item.type === 'poll')
			.filter(item => item.is_quiz || item.content.some(item => item.is_answer))
			.length,
		totalScore: 0,
		playTime: 0,
		isFinished: false,
		startTime: new Date(),
	};
};

const updateRoundScore = (delta: number) => {
	if (!state.round) {
		return;
	}
	state.round.totalScore += delta;
};

onMounted(async () => {

	const id = route.params['deck_id'];
	if (!id || typeof id !== 'string') {
		state.error = 'Deck ID required';
		return;
	}

	const { data, error } = await client.decks.load(id);
	if (!data || error) {
		state.error = error?.message || 'Unable to load deck';
		return;
	}

	state.deckID = id;
	state.isMarked = await store.decks.starred.has(id).catch(() => false);
	state.collectionID = data.collection_id;
	state.labels = data.labels;
	state.cards = data.cards;

	state.showNavigation = store.preferences.playModeShowNavigation.load();

	initRound();
});

const toggleMarked = async () => {

	if (!state.deckID) {
		return;
	}

	const { deckID, isMarked } = state;

	if (isMarked) {
		state.isMarked = await store.decks.starred.del(deckID).then(() => false).catch(() => false);
	} else {
		state.isMarked = await store.decks.starred.add(deckID).then(() => true).catch(() => false);
	}
};

const finishDeck = () => {

	if (!state.round) {
		return;
	}

	state.round.playTime = Math.floor((new Date().getTime() - state.round.startTime.getTime()) / 1000);
	state.round.isFinished = true;

	updateStats();
};

const updateStats = async () => {

	if (!state.deckID || !statsScreen.value) {
		return;
	}

	const latestScore = Math.round((statsScreen.value.score / statsScreen.value.questions) * 100);
	const stats = await store.decks.stats.load(state.deckID).catch(() => null);

	await store.decks.stats.store({
		deck_id: state.deckID,
		collection_id: state.collectionID,
		score: stats && stats.score >= latestScore ? stats.score : latestScore,
	}).catch(() => null);
};

const exitView = () => {

	if (state.collectionID) {
		router.push(`/collection/${state.collectionID}`);
		return;
	}

	router.push('/');
};

</script>

<template>

	<div class="play-view">

		<template v-if="cards?.length">

			<CardWidget v-if="!statsScreen"
				:labels="state.labels"
				:entries="cards"
				:isMarked="state.isMarked"
				:showNavigation="state.showNavigation"
				@score="updateRoundScore"
				@finish="finishDeck"
				@exit="exitView"
				@toggleMarked="toggleMarked" />

			<Endscreen v-else :stats="statsScreen" @reset="initRound" @finish="exitView" />

		</template>

		<FullscreenMessage v-else>

			<template v-if="state.error">

				<ErrorMessage v-if="state.error">

					<template v-slot:message>
						Unable to load deck
					</template>

					<template v-slot:details>
						{{ state.error }}
					</template>

				</ErrorMessage>

				<GenericButton @click="exitView">
					Go back
				</GenericButton>

			</template>

			<LoadingMessage v-else>
				Loading cards...
			</LoadingMessage>

		</FullscreenMessage>

	</div>

</template>

<style lang="scss" scoped>
	.play-view {
		position: relative;
		display: flex;
		flex-direction: column;
		align-items: center;
		width: 100%;
		height: 100%;
		height: 100vh;
		height: 100svh;
	}
</style>
