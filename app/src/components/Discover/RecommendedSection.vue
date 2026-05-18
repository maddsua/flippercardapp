<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import type { CollectionMetadata } from '../../api_models';
import LoadingMessage from '../App/LoadingMessage.vue';
import CollectionList from '../Collections/CollectionList.vue';
import CollectionListEntry from '../Collections/CollectionListEntry.vue';
import { intl, useLanguage } from '../../intl';
import { useClient } from '../../api';
import { useStorage } from '../../storage';
import { useRouter } from 'vue-router';
import ErrorMessage from '../App/ErrorMessage.vue';
import CentralMessage from '../App/CentralMessage.vue';

const client = useClient();
const store = useStorage();
const router = useRouter();
const lang = useLanguage();

interface RecommendedEntry extends CollectionMetadata {
	starred: boolean;
};

const state = reactive({
	data: null as RecommendedEntry[] | null,
	error: null as string | null,
});

const handleSelect = async (entry: RecommendedEntry) => {
	await store.collections.add(entry.id);
	router.push(`/app/collection/${entry.id}`);
};

onMounted(async () => {
	
	const { data, error } = await client.collections.list({ limit: 5 });
	if (!data || error) {
		state.error = error?.message || 'Unabale to load recommendations';
		return;
	}

	const starred = new Set(await store.collections.entries());
	state.data = data.entries.map(item => ({ ... item, starred: starred.has(item.id) }));
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

		<CollectionList v-else-if="state.data.length">
			<CollectionListEntry v-for="item of state.data"
				:title="item.name"
				:summary="item.description"
				:starrable="true"
				:starred="item.starred"
				@click="handleSelect(item)" />
		</CollectionList>

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
