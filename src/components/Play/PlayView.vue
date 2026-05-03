<script setup lang="ts">
import { reactive } from 'vue';
import CardView from '../Cards/CardView.vue';
import EndscreenView from '../Endscreen/EndscreenView.vue';
import { exampleDeck } from '../../sample.data';

//	todo: pull from data
const deckLabels = ['Test collection', 'All decks'];

const state = reactive({
	isEndscreen: false,
	score: 0,
	questions: exampleDeck.length,
	started: new Date(),
	time: 0
});

const finishDeck = () => {
	state.time = Math.floor((new Date().getTime() - state.started.getTime()) / 1000);
	state.isEndscreen = true
};

const handleNavigate = (target: 'home' | 'try_again') => {
	switch (target) {
		case 'home':
			//	todo: implement
			alert('not implemented');
			break;
		case 'try_again':
			state.isEndscreen = false;
			state.score = 0;
			state.started = new Date();
			state.time = 0;
			break;
	}
};

</script>

<template>
	<CardView v-if="!state.isEndscreen" :labels="deckLabels" :entries="exampleDeck" @score="(score) => state.score += score " @end="finishDeck" />
	<EndscreenView v-else :stats="state" @navigate="handleNavigate" />
</template>
