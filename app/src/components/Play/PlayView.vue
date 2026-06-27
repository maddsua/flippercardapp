<script setup lang="ts">
import { intl, useLanguage } from '@/intl';
import { computed, onMounted, reactive } from 'vue';
import { onBeforeRouteLeave, useRoute, useRouter } from 'vue-router';
import { useClient } from '../../api';
import { shuffleArray } from '../../arrays';
import type { CardNode } from '../../content';
import { useStorage } from '../../storage/storage';
import GenericButton from '../App/Inputs/GenericButton.vue';
import FullscreenMessage from '../App/Messages/FullscreenMessage.vue';
import LoadingMessage from '../App/Messages/LoadingMessage.vue';
import OverlayErrorMessage from '../App/Messages/OverlayErrorMessage.vue';
import Endscreen from '../Endscreen/Endscreen.vue';
import PlayableDeckScreen from './PlayableDeckScreen.vue';
import UIPrompt from '../App/Prompts/UIPrompt.vue';

const router = useRouter();
const route = useRoute();
const client = useClient();
const store = useStorage();
const lang = useLanguage();

interface RoundState {
	readonly startTime: Date;
	readonly attainableScore: number;
	hasInteracted: boolean;
	isFinished: boolean;
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
	options: {
		showNavigation: false,
		disableRotation: false,
	},
	exitPrompt: false,
	exitGuardDisabled: false,
});

const statsScreen = computed(() => state.round?.isFinished ? ({
	questions: state.round.attainableScore,
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

	const questionCount = state.cards.map(item => [item.front, item.back])
		.flat()
		.map(item => item.content)
		.flat()
		.filter(item => item.type === 'poll')
		.filter(item => item.is_quiz || item.content.some(item => item.is_answer))
		.length;

	state.round = {
		attainableScore: questionCount,
		totalScore: 0,
		playTime: 0,
		isFinished: false,
		hasInteracted: false,
		startTime: new Date(),
	};
};

const updateRoundScore = (delta: number) => {
	if (!state.round) {
		return;
	}
	state.round.hasInteracted = true;
	state.round.totalScore += delta;
};

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

const backHref = computed(() => state.collectionID ? `/collection/${state.collectionID}` : '/');

const exitView = () => {
	state.exitGuardDisabled = true;
	router.push(backHref.value);
};

const promptViewExit = (): boolean => {

	if (state.round?.hasInteracted) {
		state.exitPrompt = true;
		return false;
	}

	exitView();

	return true;
};

const handleExitPrompt = (confirmed?: boolean) => {
	if (confirmed) {
		exitView();
	}
	state.exitPrompt = false;
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

	state.options = {
		showNavigation: store.preferences.playMode.showNavigation.load(),
		disableRotation: store.preferences.playMode.disableCardRotation.load(),
	};

	initRound();
});

onBeforeRouteLeave(() => state.exitGuardDisabled ? true : promptViewExit());

</script>

<template>

	<div class="play-view">

		<OverlayErrorMessage v-if="state.error" :backHref="backHref">

			{{ intl(lang, {
				en: 'Unable to load deck',
				de: 'Datei konnten nicht geladen werden',
				uk: 'Не вдалося завантажити картки'
			}) }}

			<template v-slot:details>
				{{ state.error }}
			</template>

			<template v-slot:after>
				<GenericButton variant="thin" @click="exitView">
					{{ intl(lang, {
						en: 'Go back',
						de: 'Zurück',
						uk: 'Назад'
					}) }}
				</GenericButton>
			</template>

		</OverlayErrorMessage>

		<FullscreenMessage v-else-if="!cards?.length">
			<LoadingMessage>
				{{ intl(lang, {
					en: 'Loading cards...',
					de: 'Datei lädt...',
					uk: 'Завантаження...'
				}) }}
			</LoadingMessage>
		</FullscreenMessage>

		<template v-else>

			<PlayableDeckScreen v-if="!statsScreen"
				:labels="state.labels"
				:entries="cards"
				:isMarked="state.isMarked"
				:showNavigation="state.options.showNavigation"
				:disableRotation="state.options.disableRotation"
				@score="updateRoundScore"
				@finish="finishDeck"
				@exit="promptViewExit"
				@toggleMarked="toggleMarked" />

			<Endscreen v-else :stats="statsScreen" @reset="initRound" @finish="exitView" />

		</template>

		<UIPrompt v-if="state.exitPrompt" @done="handleExitPrompt">
			Exit game?
		</UIPrompt>

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
