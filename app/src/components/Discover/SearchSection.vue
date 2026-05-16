<script setup lang="ts">
import { reactive } from 'vue';
import { useClient } from '../../api';
import { useStorage } from '../../storage';
import { useRouter } from 'vue-router';
import { intl, useLanguage } from '../../intl';
import type { CollectionSearchResult } from '../../api_models';
import Searchbar from './Searchbar.vue';
import LoadingMessage from '../App/LoadingMessage.vue';
import CollectionList from '../Collections/CollectionList.vue';
import CollectionListEntry from '../Collections/CollectionListEntry.vue';
import ErrorMessage from '../App/ErrorMessage.vue';

const emit = defineEmits<{
	(e: 'active', state: boolean): void;
}>();

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
	locked: false,
	timer: null as number | null,
	error: null as string | null,
});

const execSearchQuery = async (term: string) => {
	
	state.error = null;
	state.data = [];

	const { data, error } = await client.collections.search(term);
	if (!data || error) {
		state.error = error?.message || 'Unabale to run a search';
		return;
	}

	const starSet = new Set(await store.collections());

	state.data = data.entries.map(item => ({ ... item, starred: starSet.has(item.id) }));
};

const handleSearchInput = (value?: string) => {

	if (!value || value.length < 2) {
		emit('active', false);
		state.data = [];
		state.prompted = false;
		return;
	}

	state.prompted = true;
	state.busy = true;

	if (state.locked) {
		return;
	}

	state.locked = true;

	if (state.timer) {
		clearTimeout(state.timer);
	}
	
	state.timer = setTimeout(async () => {
		await execSearchQuery(value);
		state.locked = false;
		state.busy = false;
	}, 250);

	emit('active', true);
};

const handleSelect = async (entry: SearchResultState) => {
	if (!await store.addCollection(entry.id)) {
		router.push(`/app/collection/${entry.id}`);
	}
	entry.starred = true;
};

</script>

<template>
	<div class="search-section">

		<Searchbar @update:modelValue="handleSearchInput" />

		<ErrorMessage v-if="state.error">
			<template v-slot:message>
				Failed to search for collections
			</template>
			<template v-slot:details>
				{{ state.error }}
			</template>
		</ErrorMessage>

		<LoadingMessage v-else-if="state.busy">
			{{ intl(lang, {
				en: 'Looking for matches...',
				uk: 'Шукаємо...',
				de: 'Suchen...'
			}) }}
		</LoadingMessage>

		<CollectionList v-else-if="state.data.length">
			<CollectionListEntry v-for="item of state.data" :title="item.name" :summary="item.description" :starrable="true" :starred="item.starred" @click="handleSelect(item)" />
		</CollectionList>

		<div v-if="!state.busy" class="searh-summary">
			<div class="summary">
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
			</div>
		</div>

	</div>
</template>

<style lang="scss" scoped>
	.search-section {
		display: flex;
		flex-direction: column;
		gap: 2rem;

		.searh-summary {
			display: flex;
			justify-content: center;

			.summary {
				font-size: 0.85rem;
				color: var(--app-theme-mysterious-white);
			}
		}
	}
</style>
