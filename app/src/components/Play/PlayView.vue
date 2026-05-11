<script setup lang="ts">
import { computed, onMounted, reactive } from 'vue';
import type { CardNode } from '../../content';
import CardWidget from '../Cards/CardWidget.vue';
import Endscreen from '../Endscreen/Endscreen.vue';
import FullscreenMessage from '../App/FullscreenMessage.vue';
import { shuffleArray } from '../../shuffle';
import { useRoute, useRouter } from 'vue-router';
import LoadingMessage from '../App/LoadingMessage.vue';
import ErrorMessage from '../App/ErrorMessage.vue';
import GenericButton from '../App/GenericButton.vue';
import { useClient } from '../../api';

const router = useRouter();
const route = useRoute();
const client = useClient();

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
	collectionID: null as string | null,
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

	let questionCount = 0;

	for (const item of state.cards) {
		for (const side of [item.back, item.front]) {
			if (side.content.some(item => item.type === 'poll')) {
				questionCount++;
				break;
			}
		}
	}

	state.round = {
		questions: questionCount,
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

	const { data, error } = await client.loadDeck(id);
	if (!data || error) {
		state.error = error?.message || 'Unable to load deck';
		return;
	}

	state.cards = data.cards.map(({ id, content }) => ({ ... content, id }));
	state.collectionID = data.collection_id;
	state.labels = data.labels;
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

	if (state.collectionID) {
		router.push(`/app/collection/${state.collectionID}`);
		return;
	}

	router.push('/app/collections');
};

</script>

<template>

	<div class="play-view">

		<template v-if="cards?.length">
			<CardWidget v-if="!statsScreen" :labels="state.labels" :entries="cards" @score="updateRoundScore" @finish="finishDeck" @exit="exitView" />
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
