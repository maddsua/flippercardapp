<script setup lang="ts">
import { reactive } from 'vue';
import { useClient } from '../../api';
import AppUI from '../App/AppUI.vue';
import AppUiHeader from '../App/AppUiHeader.vue';
import type { CollectionSearchResult } from '../../api_models';
import Searchbar from './Searchbar.vue';
import SearchSummary from './SearchSummary.vue';
import LoadingMessage from '../App/LoadingMessage.vue';
import CollectionList from '../Collections/CollectionList.vue';
import CollectionListEntry from '../Collections/CollectionListEntry.vue';
import { useStorage } from '../../storage';
import { useRouter } from 'vue-router';
import { intl, useLanguage } from '../../intl';

const client = useClient();
const store = useStorage();
const router = useRouter();
const lang = useLanguage();

interface SearchResultState extends CollectionSearchResult {
	starred: boolean;
};

const state = reactive({
	data: [] as SearchResultState[],
	prompted: false,
	busy: false,
	error: null as string | null,
});

let inputDelayTimer: number | null = null;

const handleSearch = async (term: string) => {

	if (state.busy) {
		return;
	}

	state.busy = true;
	state.error = null;
	state.data = [];

	const { data, error } = await client.collections.search(term);
	if (!data || error) {
		state.busy = false;
		state.error = error?.message || 'Unabale to run a search';
		return;
	}

	const starSet = new Set(await store.collections());

	state.data = data.entries.map(item => ({ ... item, starred: starSet.has(item.id) }));
	
	setTimeout(() => state.busy = false, 150);
};

const handleSearchInput = (event: InputEvent) => {

	const { value } = event.target as HTMLInputElement;

	if (value.length < 2) {
		state.data = [];
		state.prompted = false;
		return;
	}

	if (inputDelayTimer) {
		clearTimeout(inputDelayTimer);
	}

	inputDelayTimer = setTimeout(() => handleSearch(value), 250);
	state.prompted = true;
};

const handleSelect = async (entry: SearchResultState) => {
	if (!await store.addCollection(entry.id)) {
		router.push(`/app/collection/${entry.id}`);
	}
	entry.starred = true;
};

</script>

<template>
	<AppUI>
		<AppUiHeader>
			<template v-slot:title>
				{{ intl(lang, {
					en: 'Find collections',
					uk: 'Пошук карток',
					de: 'Nach Kärtchen suchen'
				}) }}
			</template>
		</AppUiHeader>

		<Searchbar @input="handleSearchInput" />

		<LoadingMessage v-if="state.busy">
			{{ intl(lang, {
				en: 'Looking for matches...',
				uk: 'Шукаємо...',
				de: 'Suchen...'
			}) }}
		</LoadingMessage>

		<CollectionList v-else-if="state.data.length">
			<CollectionListEntry v-for="item of state.data" :title="item.name" :summary="item.description" :starrable="true" :starred="item.starred" @click="handleSelect(item)" />
		</CollectionList>

		<SearchSummary v-if="!state.busy">
			<template v-if="state.data.length">
				{{ state.data.length }} result(s)
			</template>
			<template v-else-if="!state.prompted">
				Thinking about anything special?
			</template>
			<template v-else>
				{{ intl(lang, {
					en: 'No results',
					uk: 'Немає результатів',
					de: 'Keine Ergebnisse'
				}) }}
			</template>
		</SearchSummary>

	</AppUI>
</template>
