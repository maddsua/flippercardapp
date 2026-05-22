<script setup lang="ts">
import { reactive } from 'vue';
import { useClient } from '../../api';
import { useStorage } from '../../storage';
import { useRouter } from 'vue-router';
import { intl, useLanguage } from '../../intl';
import type { CollectionSearchResult } from '../../api_models';
import Searchbar from './Searchbar.vue';
import LoadingMessage from '../App/LoadingMessage.vue';
import ContentList from '../Content/ContentList.vue';
import ContentListEntry from '../Content/ContentListEntry.vue';
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
	score: number;
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

	const scoreMap = await store.playStats.collectionScores();
	const starred = new Set(await store.collections.entries());

	state.data = data.entries.map(item => ({
		... item,
		starred: starred.has(item.id),
		score: scoreMap.get(item.id) || 0,
	}));
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
	await store.collections.add(entry.id)
	router.push(`/app/collection/${entry.id}`);
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

		<ContentList v-else-if="state.data.length">
			<ContentListEntry v-for="item of state.data"
				:title="item.name"
				:summary="item.description"
				:visibility="item.visibility"
				:starrable="true"
				:starred="item.starred"
				:deckCount="item.size"
				:score="item.score"
				@click="handleSelect(item)" />
		</ContentList>

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
