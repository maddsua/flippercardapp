<script setup lang="ts">
import { reactive } from 'vue';
import { useClient } from '@/api';
import type { CollectionSearchResult } from '@/api_models';
import { intl, useLanguage } from '@/intl';
import LoadingMessage from '@/components/App/Messages/LoadingMessage.vue';
import ContentList from '@/components/Content/ContentList.vue';
import ContentListEntry from '@/components/Content/ContentListEntry.vue';
import Searchbar from './Searchbar.vue';
import { useStorage } from '@/storage/storage';
import InlineErrorMessage from '@/components/App/Messages/InlineErrorMessage.vue';
import { collectionCompletionMetric } from '@/play';

const props = defineProps<{
	starred: Set<string>;
}>();

const emit = defineEmits<{
	(e: 'active', state: boolean): void;
	(e: 'open', entry: SearchResultState): void;
}>();

const client = useClient();
const lang = useLanguage();
const store = useStorage();

interface SearchResultState extends CollectionSearchResult {
	starred: boolean;
	score: number;
};

const state = reactive({
	data: [] as SearchResultState[],
	prompted: false,
	busy: false,
	locked: false,
	timer: null as NodeJS.Timeout | null,
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

	const collectionStats = new Map(await store.collections.stats.aggregated(data.entries.map(item => item.id)).catch(() => []))

	state.data = data.entries.map(item => ({
		... item,
		starred: props.starred.has(item.id),
		score: collectionCompletionMetric(collectionStats, item),
	}));
};

const searchInput = (value?: string) => {

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

</script>

<template>
	<div class="search-section">

		<Searchbar @update:modelValue="searchInput" />

		<InlineErrorMessage v-if="state.error">
			<template v-slot:title>
				{{ intl(lang, {
					en: 'Failed to search for collections',
					de: 'Suchmaschine Fehler',
					uk: 'Помилка пошуку',
				}) }}
			</template>
			{{ state.error }}
		</InlineErrorMessage>

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
				:date="item.updated"
				:starrable="true"
				:starred="item.starred"
				:deckCount="item.size"
				:completion="item.score"
				@click="emit('open', item)" />
		</ContentList>

		<div v-if="!state.busy" class="searh-summary">
			<div class="summary">
				<template v-if="state.data.length">
					{{ state.data.length }} result(s)
				</template>
				<template v-else-if="!state.prompted">
					{{ intl(lang, {
						en: 'Thinking about anything special?',
						de: 'Wonach suchen wir?',
						uk: 'Що шукаємо?',
					}) }}
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
