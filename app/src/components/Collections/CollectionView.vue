<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import { useRoute, useRouter } from 'vue-router';
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
import { useClient } from '../../api';
import type { CardDeckMetadata, CollectionMetadata } from '../../api_models';
import { useStorage } from '../../storage';

const router = useRouter();
const route = useRoute();
const client = useClient();
const store = useStorage();

interface CardDeckMetadataState extends CardDeckMetadata {
	starred: boolean;
};

interface CollectionState extends CollectionMetadata{
	decks: CardDeckMetadataState[];
};

const state = reactive({
	data: null as CollectionState | null,
	starred: false,
	error: null as string | null,
});

onMounted(async () => {

	const id = route.params['collection_id'];
	if (!id || typeof id !== 'string') {
		state.error = 'Collection ID required'
		return;
	}

	const { data, error } = await client.collections.load(id);
	if (!data || error) {
		state.error = error?.message || 'Unable to load collection data';
		return;
	}

	const deckStars = new Set(await store.starred());
	const decks = data.decks.map(item => ({ ... item, starred: deckStars.has(item.id) }));

	const collectionStarred = new Set(await store.collections()).has(data.id);

	setTimeout(() => {
		state.data = ({  ... data, decks, });
		state.starred = collectionStarred;
	}, 350);
});

const openDeck = (id: string) => {
	router.push(`/app/play/deck/${id}`);
};

const closeCollection = () => {
	router.push('/app/collections');
};

const lang = useLanguage();

const toggleStar = async () => {

	state.starred = !state.starred;

	if (state.data) {
		const { id } = state.data;
		if (state.starred) {
			await store.addCollection(id);
		} else {
			await store.removeCollection(id);
		}
	}
};

</script>

<template>
	<AppUI>

		<AppUiHeader backHref="/app/collections" :starrable="true" :starred="state.starred" @toggleStar="toggleStar">

			<template v-slot:title>

				<Skeleton v-if="!state.data">
					Name placeholder
				</Skeleton>

				<template v-else-if="state.data?.name">
					{{ state.data?.name }}
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

				<Skeleton v-if="!state.data">
					Deskcription placeholder
				</Skeleton>

				<template v-else-if="state.data?.description">
					{{ state.data?.description }}
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

		<CollectionList v-if="state.data && state.data.decks.length">
			<CollectionListEntry v-for="item of state.data.decks" :title="item.name" @click="openDeck(item.id)" />
		</CollectionList>

		<CentralMessage v-else>

			<ErrorMessage v-if="state.error">

				<template v-slot:message>
					{{ intl(lang, {
						en: 'Unable to display content',
						de: 'Inhalt kann nicht angezeigt werden',
						uk: 'Не вдається відобразити вміст'
					}) }}
				</template>
				
				<template v-slot:details>
					{{ state.error }}
				</template>

			</ErrorMessage>

			<LoadingMessage v-else-if="!state.data">
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

			<GenericButton @click="closeCollection">
				{{ intl(lang, {
					en: 'Back to the list',
					de: 'Zurück zur Liste',
					uk: 'Назад до списку'
				}) }}
			</GenericButton>

		</CollectionEndlistAction>
	
	</AppUI>
</template>
