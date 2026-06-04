<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import { useClient } from '../../api';
import type { CollectionMetadata } from '../../api_models';
import { intl, useLanguage } from '../../intl';
import CentralMessage from '../App/CentralMessage.vue';
import ErrorMessage from '../App/ErrorMessage.vue';
import LoadingMessage from '../App/LoadingMessage.vue';
import ContentList from '../Content/ContentList.vue';
import ContentListEntry from '../Content/ContentListEntry.vue';
import { useStorage } from '../../storage/storage';

const client = useClient();
const lang = useLanguage();
const store = useStorage();

const props = defineProps<{
	starred: Set<string>;
}>();

const emit = defineEmits<{
	(e: 'open', entry: RecommendedEntry): void;
}>();

interface RecommendedEntry extends CollectionMetadata {
	starred: boolean;
	score: number;
};

const state = reactive({
	data: null as RecommendedEntry[] | null,
	error: null as string | null,
});

onMounted(async () => {
	
	const { data, error } = await client.collections.list({ limit: 5 });
	if (!data || error) {
		state.error = error?.message || 'Unabale to load recommendations';
		return;
	}

	const collectionStats = new Map(await store.collections.stats.aggregated(data.entries.map(item => item.id)).catch(() => []))

	state.data = data.entries.map(item => ({
		... item,
		starred: props.starred.has(item.id),
		//	todo: export
		score: (() => {
			const stat = collectionStats.get(item.id);
			if (!stat) {
				return 0;
			}
			return stat.avg_score * (stat.decks_played / (item.size ?? 1));
		})(),
	}));
});

</script>

<template>
	<div class="recommended-section">

		<div class="section-header">
			Recommended collections
		</div>

		<ErrorMessage v-if="state.error">
			<template v-slot:message>
				Failed to search for collections
			</template>
			<template v-slot:details>
				{{ state.error }}
			</template>
		</ErrorMessage>

		<LoadingMessage v-else-if="!state.data">
			{{ intl(lang, {
				en: 'Loading recommended...',
				uk: 'Завантажуємо рекомендовані...',
				de: 'Empfehlungen lädt...'
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
				@click="emit('open', item)" />
		</ContentList>

		<template v-else>
			<CentralMessage>
				Nothing to recommend yet
			</CentralMessage>
		</template>

	</div>
</template>

<style lang="scss" scoped>
	.recommended-section {
		display: flex;
		flex-direction: column;
		gap: 2rem;
		padding: 2rem 0;
		border-top: 1px solid var(--app-theme-ghostly-glow);

		.section-header {
			font-size: 0.85rem;
			font-weight: 600;
		}
	}
</style>
