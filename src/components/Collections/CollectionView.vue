<script setup lang="ts">
import { computed, onMounted, reactive } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useCollectionProvider } from '../../content.loaders';
import type { CardCollection, CardDeck } from '../../content';
import CollectionList from './CollectionList.vue';
import FullscreenMessage from '../App/FullscreenMessage.vue';
import ErrorMessage from '../App/ErrorMessage.vue';
import CollectionListEntry from './CollectionListEntry.vue';
import CollectionContainer from './CollectionContainer.vue';
import CollectionHeader from './CollectionHeader.vue';
import CollectionEndlistAction from './CollectionEndlistAction.vue';
import CollectionBreak from './CollectionBreak.vue';
import GenericButton from '../App/GenericButton.vue';
import AppUI from '../App/AppUI.vue';

const router = useRouter();
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

const stateError = computed(() => state.collection.error || state.decks.error || null);

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

const openDeck = (id: string) => {
	router.push(`/app/play/deck/${id}`);
};

const closeDeck = () => {
	router.push('/app/collections');
};

</script>

<template>
	<AppUI>

		<CollectionContainer>
	
			<CollectionHeader backHref="/app/collections">
	
				<template v-slot:title>
	
					<template v-if="state.collection.entry?.name">
						{{ state.collection.entry?.name }}
					</template>
	
					<template v-else>
						Unnamed collection
					</template>
	
				</template>
	
				<template v-slot:summary>
	
					<template v-if="state.collection.entry?.description">
						{{ state.collection.entry?.description }}
					</template>
	
					<template v-else>
						No description provided
					</template>
	
				</template>
	
			</CollectionHeader>
	
			<CollectionList v-if="state.decks.entries?.length">
				<CollectionListEntry v-for="item of state.decks.entries" :title="item.name" @click="openDeck(item.id)" />
			</CollectionList>
	
			<FullscreenMessage v-else>
	
				<ErrorMessage v-if="stateError">
	
					<template v-slot:message>
						Unable to display collection
					</template>
					
					<template v-slot:details>
						{{ stateError }}
					</template>
	
				</ErrorMessage>
	
				<p v-else>
					This collection doesn't have any cards yet!
				</p>
				
			</FullscreenMessage>
	
			<CollectionBreak />
	
			<CollectionEndlistAction>
	
				<GenericButton @click="closeDeck">
					Back to the collections
				</GenericButton>
	
			</CollectionEndlistAction>
	
		</CollectionContainer>
	</AppUI>
</template>
