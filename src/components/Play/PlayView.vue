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
import GenericButton from '../App/GenericButton.vue';

const router = useRouter();
const route = useRoute();

interface RoundState {
	startTime: Date;
	isFinished: boolean;
	questions: number;
	totalScore: number;
	playTime: number;
};

const state = reactive({
	cards: [] as CardNode[],
	labels: [] as string[],
	collection: null as CardCollection | null,
	ready: false,
	error: null as string | null,
	round: null as RoundState | null,
});

const statsScreen = computed(() => state.round?.isFinished ? ({
	questions: state.round.questions,
	score: state.round.totalScore,
	time: state.round.playTime,
}) : null);

const cards = computed(() => {

	if (!state.ready) {
		return null;
	}

	const entries = [...state.cards];
	shuffleArray(entries);
	return entries;
});

const initRound = () => {

	const questions = state.cards
		.map(item => [item.back, item.front])
		.flat()
		.map(item => item.content.some(item => item.type === 'poll')).length;

	state.round = {
		questions,
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

const loadDeck = async () => {

	const id = route.params['deck_id'];
	if (!id || typeof id !== 'string') {
		return { error: new Error('Deck ID required') };
	}

	const { data, error } = await useCollectionProvider().decks(id);
	if (!data || error) {
		return { error: error || new Error('Unable to load a deck') };
	}

	if (data.length === 0) {
		return { error: new Error('Deck ID not found') };
	}

	return { deck: data[0] };
};

const loadParentCollection = async (deck: CardDeck) => {

	const { data, error } = await deck.collection();
	if (!data || error) {
		state.error = error?.message || 'Unable to load parent collection';
		return new Error(state.error);
	}

	state.collection = data;
	state.labels = [data.name, deck.name];
	return null;
};

const loadDeckCards = async (deck: CardDeck) => {

	const { data, error } = await deck.cards();
	if (!data) {
		state.error = error?.message || 'unable to load cards';
		return new Error(state.error);
	}

	if (data.length === 0) {
		state.error = 'Empty deck';
		return new Error(state.error);
	}

	state.cards = data;
	return null;
};

onMounted(async () => {

	const { deck, error } = await loadDeck();
	if (!deck || error) {
		state.error = error.message;
		return;
	}

	if (!!await loadParentCollection(deck)) {
		return;
	}

	if (!!await loadDeckCards(deck)) {
		return;
	}

	state.ready = true;

	initRound();
});

const finishDeck = () => {

	if (!state.round) {
		return;
	}

	state.round.playTime = Math.floor((new Date().getTime() - state.round.startTime.getTime()) / 1000);
	state.round.isFinished = true;
};

const exitView = () => {

	if (state.collection) {
		router.push(`/app/collection/${state.collection.id}`);
		return;
	}

	router.push('/app/collections');
};	

</script>

<template>

	<template v-if="cards?.length">
		<CardView v-if="!statsScreen" :labels="state.labels" :entries="cards" @score="updateRoundScore" @finish="finishDeck" @exit="exitView" />
		<EndscreenView v-else :stats="statsScreen" @reset="initRound" @finish="exitView" />
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

</template>
