<script setup lang="ts">
import { onMounted, ref } from 'vue';
import PlayView from './components/Play/PlayView.vue';
import type { CardDeck } from './components/Cards/content';
import { sampleProvider } from './data/sample';

const deck = ref<CardDeck | null>(null);

onMounted(async () => {

	const { data: collections } = await sampleProvider.collections();

	const { data: decks } = await collections[0].decks();

	deck.value = decks[0];
});

</script>

<template>
	<PlayView v-if="deck" :deck="deck" />
	<div v-else>
		No content loaded
	</div>
</template>
