<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import { intl, useLanguage } from '../../intl';
import AppUI from '../App/AppUI.vue';
import AppUiHeader from '../App/AppUiHeader.vue';
import CollectionList from '../Collections/CollectionList.vue';
import CollectionListEntry from '../Collections/CollectionListEntry.vue';
import type { CardDeckMetadata } from '../../api_models';
import { useRouter } from 'vue-router';
import LoadingMessage from '../App/LoadingMessage.vue';
import ErrorMessage from '../App/ErrorMessage.vue';
import CentralMessage from '../App/CentralMessage.vue';
import { useStorage } from '../../storage';
import { useClient } from '../../api';

const state = reactive({
	data: null as CardDeckMetadata[] | null,
	error: null as string | null
});

const router = useRouter();
const lang = useLanguage();
const store = useStorage();
const client = useClient();

onMounted(async () => {

	const starredIDs = await store.starred();
	if (!starredIDs.length) {
		state.data = [];
		return;
	}

	const { data, error } = await client.listDecks({ ids: starredIDs });
	if (!data || error) {
		state.error = error?.message || 'Unable to load decks';
		return;
	}

	state.data = data.entries;
});

const openDeck = (id: string) => {
	router.push(`/app/play/deck/${id}`);
};

</script>

<template>
	<AppUI>

		<AppUiHeader>
			<template v-slot:title>
				{{ intl(lang, {
					en: 'Starred cards',
					de: 'Markierte Karten',
					uk: 'Фаворити',
				}) }}
			</template>
		</AppUiHeader>

		<CollectionList v-if="state.data && state.data.length">
			<CollectionListEntry v-for="item of state.data" :title="item.name" :starred="true" @click="openDeck(item.id)" />
		</CollectionList>

		<CentralMessage v-else>

			<ErrorMessage v-if="state.error">

				<template v-slot:message>
					{{ intl(lang, {
						en: 'Unable to display content',
						de: 'Inhalt kann nicht angezeigt werden',
						uk: 'Не вдається відобразити вміст'
					}) }}
				</template>
				
				<template v-slot:details>
					{{ state.error }}
				</template>

			</ErrorMessage>

			<LoadingMessage v-else-if="!state.data">
				{{ intl(lang, {
					en: 'Loading...',
					de: 'Lädt...',
					uk: 'Один момент...'
				}) }}
			</LoadingMessage>

			<p v-else>
				{{ intl(lang, {
					en: `You haven't starred anything`,
					de: 'Sie habend noch keine Sterne gegeben',
					uk: 'Ви ще нічого не зберегли'
				}) }}
			</p>

		</CentralMessage>

	</AppUI>
</template>
