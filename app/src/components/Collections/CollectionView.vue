<script setup lang="ts">
import { computed, onMounted, reactive } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useClient } from '@/api';
import type { AuthState, CardDeckMetadata, CollectionMetadata } from '@/api_models';
import { intl, useLanguage } from '@/intl';
import { useStorage } from '@/storage/storage';
import AppUI from '../App/Layout/AppUI.vue';
import AppUiHeader from '../App/Layout/AppUiHeader.vue';
import CentralMessage from '../App/Messages/CentralMessage.vue';
import GenericButton from '../App/Inputs/GenericButton.vue';
import LoadingMessage from '../App/Messages/LoadingMessage.vue';
import Skeleton from '../App/Messages/Skeleton.vue';
import ContentList from '../Content/ContentList.vue';
import ContentListEntry from '../Content/ContentListEntry.vue';
import ContentEntryBadge from '../Content/ContentEntryBadge.vue';
import InlineErrorMessage from '../App/Messages/InlineErrorMessage.vue';

const router = useRouter();
const route = useRoute();
const client = useClient();
const store = useStorage();
const lang = useLanguage();

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
	auth: null as  AuthState | null,
	error: null as string | null,
});

const backRef = computed(() => state.auth?.actor?.permissions.team_member ? '/collections/all' : '/');

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

	const starredDecks = new Set(await store.decks.starred.all().catch(() => []));

	const loadedDeckIDSet = new Set(data.decks.map(item => item.id));

	const deckScoreMap = new Map(await store.decks.stats
		.filter(val => loadedDeckIDSet.has(val.deck_id))
		.catch(() => [])
		.then(entries => entries.map(entry => [entry.deck_id, entry.score])));

	const deckEntries = data.decks.map(item => ({
		... item,
		starred: starredDecks.has(item.id),
		score: deckScoreMap.get(item.id) || 0,
	}));

	state.data = {
		... data,
		decks: deckEntries.sort((a,b) => (b.starred ? 1 : 0) - (a.starred ? 1 : 0)),
	};

	state.starred = await store.collections.starred.has(data.id).catch(() => false);

	state.auth = await client.auth.whoami({ cached: true }).then(res => res.data || null);
});

const playDeck = (id: string) => {
	router.push(`/play/deck/${id}`);
};

const toggleStar = async () => {

	if (!state.data) {
		return;
	}

	const { id } = state.data;

	if (state.starred) {
		state.starred = await store.collections.starred.del(id).then(() => false).catch(() => false);
	} else {
		state.starred = await store.collections.starred.add(id).then(() => true).catch(() => false);
	}
};

const openMetadataEditor = () => {
	router.push(`/collection/${state.data!.id}/edit`);
};

const openNewDeckEditor = () => {
	router.push(`/decks/editor?collection_id=${state.data!.id}`);
};

const openDeckEditor = (id: string) => {
	router.push(`/decks/editor/${id}`);
};

const deleteDeck = async (deckID: string) => {

	if (!state.data) {
		throw new Error('Invalid state');
	}

	if (!confirm('Delete deck?')) {
		return;
	}

	const { error } = await client.decks.remove(deckID);
	if (error) {
		console.error('Unable to delete collection deck:', error.message);
		return;
	}

	state.data.decks = state.data.decks?.filter(item => item.id !== deckID) || null;
};

const visibilityIcons = {
	'PUBLIC': 'world',
	'HIDDEN': 'link',
	'PRIVATE': 'lock',
} as const;

const fmtDate = (date: string) => new Date(date).toLocaleDateString('en-UK', {
	year: 'numeric',
	month: 'numeric',
	day: 'numeric',
	hour: 'numeric',
	minute: 'numeric',
	second: 'numeric',
});

const capitalize = (text: string) => text.slice(0, 1).toUpperCase() + text.slice(1).toLowerCase();

</script>

<template>
	<AppUI>

		<AppUiHeader :backHref="backRef" :starrable="true" :starred="state.starred" @toggleStar="toggleStar">

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

			<template v-if="state.data" v-slot:badges>

				<div class="badges-row">

					<ContentEntryBadge :icon="visibilityIcons[state.data.visibility]">
						{{ capitalize(state.data.visibility) }}
					</ContentEntryBadge>

					<ContentEntryBadge icon="clock">
						{{ fmtDate(state.data.updated) }}
					</ContentEntryBadge>

					<ContentEntryBadge icon="decks">
						{{ state.data?.size.toFixed(0) }}
					</ContentEntryBadge>

				</div>

			</template>

			<template v-if="state.data && state.auth?.actor?.permissions.team_member" v-slot:actions>
				<GenericButton variant="thin" theme="green" @click="openNewDeckEditor">
					+ Add deck
				</GenericButton>
				<GenericButton variant="thin" theme="orange" @click="openMetadataEditor">
					Manage
				</GenericButton>
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
				:editable="state.auth?.actor?.permissions.content_edit"
				:playable="true"
				@click="playDeck(item.id)"
				@edit="openDeckEditor(item.id)"
				@delete="deleteDeck(item.id)" />
		</ContentList>

		<CentralMessage v-else>

			<InlineErrorMessage v-if="state.error">

				<template v-slot:title>
					{{ intl(lang, {
						en: 'Unable to display content',
						de: 'Inhalt kann nicht angezeigt werden',
						uk: 'Не вдається відобразити вміст'
					}) }}
				</template>

				{{ state.error }}

			</InlineErrorMessage>

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
	
	</AppUI>
</template>

<style lang="scss" scoped>
	.badges-row {
		display: flex;
		flex-flow: row wrap;
		font-size: 0.75rem;
		gap: 1em;
		align-items: center;
	}
</style>
