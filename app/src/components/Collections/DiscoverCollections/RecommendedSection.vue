<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import { useClient } from '@/api';
import type { CollectionMetadata } from '@/api_models';
import { intl, useLanguage } from '@/intl';
import CentralMessage from '@/components/App/Messages/CentralMessage.vue';
import LoadingMessage from '@/components/App/Messages/LoadingMessage.vue';
import ContentList from '@/components/Content/ContentList.vue';
import ContentListEntry from '@/components/Content/ContentListEntry.vue';
import { useStorage } from '@/storage/storage';
import InlineErrorMessage from '@/components/App/Messages/InlineErrorMessage.vue';
import { distributeCollectionPlayScore } from '@/play';

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
		score: distributeCollectionPlayScore(collectionStats, item),
	}));
});

</script>

<template>
	<div class="recommended-section">

		<div class="section-header">
			Recommended collections
		</div>

		<InlineErrorMessage v-if="state.error">
			<template v-slot:title>
				Failed to search for collections
			</template>
			{{ state.error }}
		</InlineErrorMessage>

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

		.section-header {
			font-size: 0.85rem;
			font-weight: 600;
		}
	}
</style>
