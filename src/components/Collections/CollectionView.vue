<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import { useRoute } from 'vue-router';
import { useCollectionProvider } from '../../content.loaders';
import type { CardCollection, CardDeck } from '../../content';

const route = useRoute();

const state = reactive({
	collection: {
		entry: null as CardCollection | null,
		error: null as string | null,
	},
	decks: {
		entries: null as CardDeck[] | null,
		error: null as string | null,
	},
});

const loadCollection = async () => {

	const id = route.params['collection_id'];
	if (!id || typeof id !== 'string') {
		state.collection.error = 'Collection ID required'
		return;
	}

	const { data, error } = await useCollectionProvider().collections(id);
	if (!data || error) {
		state.collection.error = error?.message || 'Unable to load collection data';
		return;
	}

	if (data.length === 0) {
		state.collection.error = 'Collection not found';
		return;
	}

	state.collection.entry = data[0];
};

const loadDecks = async (collection: CardCollection) => {

	const { data, error } = await collection.decks();
	if (!data || error) {
		state.decks.error = error?.message || 'Unable to load decks';
		return;
	}

	state.decks.entries = data;
};

onMounted(async () => {

	await loadCollection();

	if (state.collection.entry) {
		await loadDecks(state.collection.entry);
	}
	
});

//	todo: make look nice

</script>

<template>
	<div class="collection-view">

		<p>
			Decks
		</p>

		<p>
			<RouterLink to="/collections">
				Back to the list
			</RouterLink>
		</p>

		<p>
			Collection: {{ state.collection.entry?.name }}
		</p>

		<ul v-if="state.decks.entries?.length">
			<li v-for="item of state.decks.entries">
				<RouterLink :to="`/play/deck/${item.id}`">
					{{ item.name }}
				</RouterLink>
			</li>
		</ul>
		<div v-else class="message">
			No decks available
		</div>
	</div>
</template>

<style lang="scss" scoped>

</style>
