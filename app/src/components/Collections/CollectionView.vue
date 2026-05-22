<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import ContentList from '../Content/ContentList.vue';
import ErrorMessage from '../App/ErrorMessage.vue';
import ContentListEntry from '../Content/ContentListEntry.vue';
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

interface DeckEntry extends CardDeckMetadata {
	starred: boolean;
	score: number;
};

interface CollectionEntry extends CollectionMetadata {
	decks: DeckEntry[];
};

const state = reactive({
	data: null as CollectionEntry | null,
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

	const starredDecks = new Set(await store.starredDecks.entries());
	const playStats = new Map(await store.playStats.entries());

	const deckEntries = data.decks .map(item => ({
		... item,
		starred: starredDecks.has(item.id),
		score: playStats.get(item.id)?.score || 0,
	}));

	state.data = {
		... data,
		decks: deckEntries.sort((a,b) => (b.starred ? 1 : 0) - (a.starred ? 1 : 0)),
	};

	state.starred = await store.collections.contains(data.id);
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
			await store.collections.add(id);
		} else {
			await store.collections.remove(id);
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

		<ContentList v-if="state.data && state.data.decks.length">
			<ContentListEntry v-for="item of state.data.decks"
				:title="item.name"
				:summary="item.description"
				:visibility="item.visibility"
				:cardCount="item.size"
				:starred="item.starred"
				:score="item.score"
				@click="openDeck(item.id)" />
		</ContentList>

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
