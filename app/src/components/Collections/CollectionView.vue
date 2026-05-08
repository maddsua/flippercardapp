<script setup lang="ts">
import { computed, onMounted, reactive } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useContent } from '../../content.loaders';
import type { CardCollection, CardDeck } from '../../content';
import CollectionList from './CollectionList.vue';
import ErrorMessage from '../App/ErrorMessage.vue';
import CollectionListEntry from './CollectionListEntry.vue';
import CollectionEndlistAction from './CollectionEndlistAction.vue';
import CollectionBreak from './CollectionBreak.vue';
import GenericButton from '../App/GenericButton.vue';
import AppUI from '../App/AppUI.vue';
import { intl, useLanguage } from '../../intl';
import CentralMessage from '../App/CentralMessage.vue';
import LoadingMessage from '../App/LoadingMessage.vue';
import Skeleton from '../App/Skeleton.vue';
import AppUiHeader from '../App/AppUiHeader.vue';

const router = useRouter();
const route = useRoute();

const state = reactive({
	collection: {
		entry: null as CardCollection | null,
		ready: false,
		error: null as string | null,
	},
	decks: {
		entries: [] as CardDeck[],
		ready: false,
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

	const { data, error } = await useContent().collections(id);
	if (!data || error) {
		state.collection.error = error?.message || 'Unable to load collection data';
		return;
	}

	if (data.length === 0) {
		state.collection.error = 'Collection not found';
		return;
	}

	state.collection.entry = data[0];
	setTimeout(() => state.collection.ready = true, 250);
};

const loadDecks = async (collection: CardCollection) => {

	const { data, error } = await collection.decks();
	if (!data || error) {
		state.decks.error = error?.message || 'Unable to load decks';
		return;
	}

	state.decks.entries = data;
	setTimeout(() => state.decks.ready = true, 350);
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

const lang = useLanguage();

</script>

<template>
	<AppUI>

		<AppUiHeader backHref="/app/collections">

			<template v-slot:title>

				<Skeleton v-if="!state.collection.ready">
					Name placeholder
				</Skeleton>

				<template v-else-if="state.collection.entry?.name">
					{{ state.collection.entry?.name }}
				</template>

				<template v-else>
					{{ intl(lang, {
						en: 'Unnamed collection',
						de: 'Unbenannte Sammlung',
						uk: 'Безіменна колекція'
					}) }}
				</template>

			</template>

			<template v-slot:summary>

				<Skeleton v-if="!state.collection.ready">
					Deskcription placeholder
				</Skeleton>

				<template v-else-if="state.collection.entry?.description">
					{{ state.collection.entry?.description }}
				</template>

				<template v-else>
					{{ intl(lang, {
						en: 'No description provided',
						de: 'Keine Beschreibung vorhanden',
						uk: 'Опис не надано'
					}) }}
				</template>

			</template>

		</AppUiHeader>

		<CollectionList v-if="state.decks.ready && state.decks.entries.length">
			<CollectionListEntry v-for="item of state.decks.entries" :title="item.name" @click="openDeck(item.id)" />
		</CollectionList>

		<CentralMessage v-else>

			<ErrorMessage v-if="stateError">

				<template v-slot:message>
					{{ intl(lang, {
						en: 'Unable to display content',
						de: 'Inhalt kann nicht angezeigt werden',
						uk: 'Не вдається відобразити вміст'
					}) }}
				</template>
				
				<template v-slot:details>
					{{ stateError }}
				</template>

			</ErrorMessage>

			<LoadingMessage v-else-if="!state.decks.ready">
				{{ intl(lang, {
					en: 'Loading...',
					de: 'Lädt...',
					uk: 'Один момент...'
				}) }}
			</LoadingMessage>

			<p v-else>
				{{ intl(lang, {
					en: `This collection doesn't have any cards yet!`,
					de: 'Diese Sammlung enthält noch keine Karten!',
					uk: 'У цій колекції ще немає жодної картки!'
				}) }}
			</p>
			
		</CentralMessage>

		<CollectionBreak />

		<CollectionEndlistAction>

			<GenericButton @click="closeDeck">
				{{ intl(lang, {
					en: 'Back to the list',
					de: 'Zurück zur Liste',
					uk: 'Назад до списку'
				}) }}
			</GenericButton>

		</CollectionEndlistAction>
	
	</AppUI>
</template>
