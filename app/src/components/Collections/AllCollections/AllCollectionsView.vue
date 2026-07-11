<script setup lang="ts">
import { useClient, type Pagination } from '@/api';
import type { AuthState, CollectionMeta } from '@/api_models';
import AppUI from '@/components/App/Layout/AppUI.vue';
import AppUiHeader from '@/components/App/Layout/AppUiHeader.vue';
import CentralMessage from '@/components/App/Messages/CentralMessage.vue';
import GenericButton from '@/components/App/Inputs/GenericButton.vue';
import LoadingMessage from '@/components/App/Messages/LoadingMessage.vue';
import ContentList from '@/components/Content/ContentList.vue';
import ContentListEntry from '@/components/Content/ContentListEntry.vue';
import { genericPageState, pageControls } from '@/dataloader';
import { intl, useLanguage } from '@/intl';
import { onMounted, reactive } from 'vue';
import { useRouter } from 'vue-router';
import InlineErrorMessage from '@/components/App/Messages/InlineErrorMessage.vue';
import { useStorage } from '@/storage/storage';
import { collectionCompletionMetric } from '@/play';

const lang = useLanguage();
const router = useRouter();
const client = useClient();
const store = useStorage();

interface Entry extends CollectionMeta {
	completion: number;
	starred: boolean;
};

const state = reactive({
	page: genericPageState<Entry>(),
	auth: null as  AuthState | null,
});

const { more: loadMore } = pageControls<Entry>(state.page, async (pagination: Pagination) => {

	const { data, error } = await client.collections.list(pagination);
	if (!data || error) {
		return { data: null, error };
	}

	const collectionStats = new Map(await store.collections.stats.aggregated(data.entries.map(item => item.id)).catch(() => []));
	const starSet = new Set(await store.collections.starred.all().catch(() => []));

	return {
		data: {
			... data,
			entries: data.entries.map(entry => ({
				... entry,
				completion: collectionCompletionMetric(collectionStats, entry),
				starred: starSet.has(entry.id),
			}))
		},
		error: null,
	}
});

onMounted(loadMore);

onMounted(async () => {
	state.auth = await client.auth.whoami({ cached: true }).then(res => res.data || null);
});

const openNewCollectionEditor = () => {
	router.push('/collections/new');
};

const openCollection = async (id: string) => {
	router.push(`/collection/${id}`);
};

</script>

<template>
	<AppUI>

		<AppUiHeader backHref="/collections">

			<template v-slot:title>
				{{ intl(lang, {
					en: 'All collections',
					uk: 'Всі колекції',
					de: 'Alle Sammlungen'
				}) }}
			</template>

			<template v-slot:summary>
				Full content list
			</template>

			<template v-if="state.auth?.actor?.permissions.team_member" v-slot:actions>
				<GenericButton variant="thin" theme="green" @click="openNewCollectionEditor">
					+ Add collection
				</GenericButton>
			</template>

		</AppUiHeader>

		<InlineErrorMessage v-if="state.page.error">
			<template v-slot:title>
				Failed to load collection list
			</template>
			{{ state.page.error }}
		</InlineErrorMessage>

		<LoadingMessage v-else-if="!state.page.ready">
			{{ intl(lang, {
				en: 'Loading collection...',
				uk: 'Завантажуємо колекції...',
				de: 'Sammlungen laden...'
			}) }}
		</LoadingMessage>

		<template v-else>

			<template v-if="state.page.entries.length">

				<ContentList>
					<ContentListEntry v-for="entry of state.page.entries"
						:title="entry.name"
						:summary="entry.description"
						:date="entry.updated"
						:visibility="entry.visibility"
						:deckCount="entry.size"
						:starrable="true"
						:starred="entry.starred"
						:completion="entry.completion"
						:themeColor="entry.theme_color"
						@click="openCollection(entry.id)" />
				</ContentList>

				<div v-if="state.page.has_next" class="actions">
					<GenericButton theme="green" variant="thin" @click="loadMore">
						Load more
					</GenericButton>
				</div>

			</template>

			<template v-else>
				<CentralMessage>
					No collections yet
				</CentralMessage>
			</template>

		</template>

	</AppUI>
</template>

<style lang="scss" scoped>
	.actions {
		display: flex;
		flex-flow: row nowrap;
		justify-content: center;
		align-items: center;
	}
</style>
