<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import { intl, useLanguage } from '@/intl';
import { useRouter } from 'vue-router';
import { useStorage } from '@/storage/storage';
import AppUI from '@/components/App/Layout/AppUI.vue';
import AppUiHeader from '@/components/App/Layout/AppUiHeader.vue';
import RecommendedSection from './RecommendedSection.vue';
import SearchSection from './SearchSection.vue';
import GenericButton from '@/components/App/Inputs/GenericButton.vue';
import type { AuthState } from '@/api_models';
import { useClient } from '@/api';
import CollectionBreak from '../CollectionBreak.vue';

const lang = useLanguage();
const store = useStorage();
const router = useRouter();
const client = useClient();

const state = reactive({
	searching: false,
	starred: new Set<string>(),
	auth: null as  AuthState | null,
});

onMounted(async () => {
	state.starred = new Set(await store.collections.starred.all().catch(() => []));
	state.auth = await client.auth.whoami({ cached: true }).then(res => res.data || null);
});

const openCollection = async (entry: { id: string }) => {
	store.collections.starred.add(entry.id).catch(() => null);
	router.push(`/collection/${entry.id}`);
};

const openNewCollectionEditor = () => {
	router.push('/collections/new');
};

const openFullList = () => {
	router.push('/collections/all');
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

			<template v-slot:summary>
				{{ intl(lang, {
					en: 'Get new stuff to exercise with',
					de: 'Hol dir neue Material fürs Training',
					uk: 'Знаходь нові штуки для тренувань',
				}) }}
			</template>

			<template v-slot:actions>
				<GenericButton v-if="state.auth?.actor?.permissions.team_member" variant="thin" theme="green" @click="openNewCollectionEditor">
					+ Add collection
				</GenericButton>
				<GenericButton variant="thin" @click="openFullList">
					{{ intl(lang, {
						en: 'Show all',
						de: 'Alle anzeigen',
						uk: 'Показати всі',
					}) }}
				</GenericButton>
			</template>
		</AppUiHeader>

		<SearchSection
			:starred="state.starred"
			@active="val => state.searching = val"
			@open="openCollection" />

		<CollectionBreak />
		
		<KeepAlive>
			<RecommendedSection v-if="!state.searching"
				:starred="state.starred"
				@open="openCollection" />
		</KeepAlive>

	</AppUI>
</template>
