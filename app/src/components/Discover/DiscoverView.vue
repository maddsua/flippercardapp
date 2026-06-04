<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import { intl, useLanguage } from '../../intl';
import { useRouter } from 'vue-router';
import { useStorage } from '../../storage/storage';
import AppUI from '../App/AppUI.vue';
import AppUiHeader from '../App/AppUiHeader.vue';
import RecommendedSection from './RecommendedSection.vue';
import SearchSection from './SearchSection.vue';

const lang = useLanguage();
const store = useStorage();
const router = useRouter();

const state = reactive({
	searching: false,
	starred: new Set<string>(),
});

onMounted(async () => {
	state.starred = new Set(await store.collections.starred.all().catch(() => []));
});

const openCollection = async (entry: { id: string }) => {
	store.collections.starred.add(entry.id).catch(() => null);
	router.push(`/collection/${entry.id}`);
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

		<SearchSection
			:starred="state.starred"
			@active="val => state.searching = val"
			@open="openCollection" />
		
		<KeepAlive>
			<RecommendedSection v-if="!state.searching"
				:starred="state.starred"
				@open="openCollection" />
		</KeepAlive>

	</AppUI>
</template>
