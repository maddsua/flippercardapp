<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import { useClient, type Pagination } from '../../../api';
import type { CardDeckMetadata, CollectionMetadata } from '../../../api_models';
import { useRouter } from 'vue-router';
import { genericPageState, pageControls } from '../../../dataloader';
import AppUiHeader from '../../App/AppUiHeader.vue';
import LoadingMessage from '../../App/LoadingMessage.vue';
import ErrorMessage from '../../App/ErrorMessage.vue';
import CollectionList from './Collections/CollectionList.vue';
import GenericButton from '../../App/GenericButton.vue';
import CollectionEntry from './Collections/CollectionEntry.vue';

const client = useClient();
const router = useRouter();

interface LazyLoadedCollection extends CollectionMetadata {
	decks?: CardDeckMetadata[] | null;
};

const state = reactive(genericPageState<LazyLoadedCollection>());

const { more: loadMore } = pageControls(state, async (pagination: Pagination) => {
	return client.collections.list(pagination);
});

onMounted(loadMore);

const openCollectionMetaEdit = (id: string) => {
	router.push(`/app/dashboard/content/collection/${id}/metadata`);
};

const openDeckEditor = (params: { deckID?: string; collectionID?: string }) => {
	if (params.deckID) {
		router.push(`/app/editor/deck/${params.deckID}/editor`);
	} else if (params.collectionID) {
		router.push(`/app/editor/deck/editor?collection_id=${params.collectionID}`);
	}
};

const loadCollectionDecks = async (entry: LazyLoadedCollection) => {

	const { data, error } = await client.decks.list({ collection_id: entry.id });
	if (error) {
		console.error('Unable to load collection decks:', error.message);
		return;
	}
	
	entry.decks = data?.entries || [];
};

const deleteCollectionDeck = async (entry: LazyLoadedCollection, deckID: string) => {

	if (!confirm('Delete deck?')) {
		return;
	}

	const { error } = await client.decks.remove(deckID);
	if (error) {
		console.error('Unable to delete collection deck:', error.message);
		return;
	}

	entry.decks = entry.decks?.filter(item => item.id !== deckID);
	entry.size--;
};

</script>

<template>

	<AppUiHeader backHref="/app/dashboard">
		<template v-slot:title>
			Content dashboard
		</template>
		<template v-slot:summary>
			Manage cards and collections
		</template>
	</AppUiHeader>

	<LoadingMessage v-if="!state.ready">
		Loading collections...
	</LoadingMessage>

	<ErrorMessage v-else-if="state.error">
		<template v-slot:message>
			Unable to load collections
		</template>
		<template v-slot:details>
			{{ state.error }}
		</template>
	</ErrorMessage>

	<CollectionList v-else>

		<template v-slot:actions_before>
			<GenericButton theme="blue" variant="thin" @click="router.push('/app/dashboard/content/collections/new')">
				+ Add collection
			</GenericButton>
		</template>

		<CollectionEntry v-for="entry of state.entries" :key="entry.id" :entry="entry"
			@edit="openCollectionMetaEdit(entry.id)"
			@showDecks="loadCollectionDecks(entry)"
			@addDeck="openDeckEditor({ collectionID: entry.id })"
			@editDeck="deckID => openDeckEditor({ deckID })"
			@deleteDeck="deckID => deleteCollectionDeck(entry, deckID)" />

		<template v-if="state.has_next" v-slot:actions_after>
			<GenericButton theme="green" variant="thin" @click="loadMore">
				Load more
			</GenericButton>
		</template>

	</CollectionList>

</template>

<style lang="scss" scoped>

</style>
