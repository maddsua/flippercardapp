<script setup lang="ts">
import { computed, onMounted, reactive } from 'vue';
import type { CardCollection, CardDeck, CardNode } from '../../content';
import CardView from '../Cards/CardView.vue';
import EndscreenView from '../Endscreen/EndscreenView.vue';
import FullscreenMessage from '../App/FullscreenMessage.vue';
import { shuffleArray } from '../../shuffle';
import { useRoute, useRouter } from 'vue-router';
import { useCollectionProvider } from '../../content.loaders';
import LoadingMessage from '../App/LoadingMessage.vue';
import ErrorMessage from '../App/ErrorMessage.vue';
import Button from '../App/Button.vue';

const router = useRouter();
const route = useRoute();

interface GameState {
	startTime: Date;
	isFinished: boolean;
	questions: number;
	totalScore: number;
	playTime: number;
};

const state = reactive({
	data: {
		cards: [] as CardNode[],
		labels: [] as string[],
		collection: null as CardCollection | null,
		busy: false,
		ready: false,
		error: null as string | null,
	},
	game: null as GameState | null,
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

//	todo: refactor all the shitty loaders

const unwrapData = async (deck: CardDeck) => {

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

	const collection = await deck.collection();

	if (!collection.data) {
		state.data.labels = [deck.name];
		return;
	}

	state.data.collection = collection.data;
	state.data.labels = [state.data.collection.name, deck.name];
};

const loadDeck = async (deck: CardDeck) => {

	if (state.data.busy) {
		return;
	}

	state.data.busy = true;

	state.data.error = null;
	state.data.ready = false;
	state.data.cards = [];

	await unwrapData(deck);

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

	const id = route.params['deck_id'];
	if (!id || typeof id !== 'string') {
		state.data.error = 'Deck ID required'
		return;
	}

	const { data, error } = await useCollectionProvider().decks(id);
	if (!data || error) {
		state.data.error = error?.message || 'Unable to load a deck'
		return;
	}

	if (data.length === 0) {
		state.data.error = 'Deck ID not found';
		return;
	}

	await loadDeck(data[0]);

	initGame();
});

const finishDeck = () => {

	if (!state.game) {
		return;
	}

	state.game.playTime = Math.floor((new Date().getTime() - state.game.startTime.getTime()) / 1000);
	state.game.isFinished = true;
};

const exitView = () => {

	if (state.data.collection) {
		router.push(`/app/collection/${state.data.collection.id}`);
		return;
	}

	router.push('/app/collections');
};	

</script>

<template>

	<template v-if="cards?.length">
		<CardView v-if="!endgameStats" :labels="state.data.labels" :entries="cards" @score="updateGameScore" @end="finishDeck" />
		<EndscreenView v-else :stats="endgameStats" @reset="initGame" @finish="exitView" />
	</template>

	<FullscreenMessage v-else>

		<template v-if="state.data.error">

			<ErrorMessage v-if="state.data.error">
	
				Unable to load deck
	
				<template v-slot:details>
					{{ state.data.error }}
				</template>
	
			</ErrorMessage>

			<Button @click="exitView">
				Go back
			</Button>

		</template>

		<LoadingMessage v-else>
			Loading cards...
		</LoadingMessage>

	</FullscreenMessage>

</template>
