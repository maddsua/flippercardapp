<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import { useRouter } from 'vue-router';
import CollectionList from './CollectionList.vue';
import CollectionListEntry from './CollectionListEntry.vue';
import ErrorMessage from '../App/ErrorMessage.vue';
import CollectionBreak from './CollectionBreak.vue';
import CollectionEndlistAction from './CollectionEndlistAction.vue';
import GenericButton from '../App/GenericButton.vue';
import AppUI from '../App/AppUI.vue';
import { intl, useLanguage } from '../../intl';
import CentralMessage from '../App/CentralMessage.vue';
import LoadingMessage from '../App/LoadingMessage.vue';
import AppUiHeader from '../App/AppUiHeader.vue';
import { useClient } from '../../api';
import type { CollectionMetadata } from '../../api_models';
import { useStorage } from '../../storage';

const router = useRouter();
const client = useClient();
const store = useStorage();

const state = reactive({
	data: null as CollectionMetadata[] | null,
	error: null as string | null,
});

onMounted(async () => {

	const myCollections = await store.collections();
	if (!myCollections.length) {
		state.data = [];
		return;
	}

	const { data, error } = await client.listCollections({ ids: myCollections });
	if (!data || error) {
		state.error = error?.message || 'Unable to load collections';
		return;
	}

	setTimeout(() => state.data = data.entries, 250);
});

const openCollection = (id: string) => {
	router.push(`/app/collection/${id}`);
};

const openExplore = () => {
	router.push('/app/discover');
};

const lang = useLanguage();

</script>

<template>

	<AppUI>

		<AppUiHeader>

			<template v-slot:title>
				{{ intl(lang, {
					en: `My collections`,
					de: 'Meine Karten',
					uk: 'Мої картки'
				}) }}
			</template>

			<template v-slot:summary>
				{{ intl(lang, {
					en: `Your card collections`,
					de: 'Ihre Kartensammlungen',
					uk: 'Ваші колекції карток'
				}) }}
			</template>

		</AppUiHeader>

		<CollectionList v-if="state.data !== null && state.data.length">
			<CollectionListEntry v-for="item of state.data" :title="item.name" @click="openCollection(item.id)" />
		</CollectionList>

		<CentralMessage v-else>

			<ErrorMessage v-if="state.error">

				<template v-slot:message>
					{{ intl(lang, {
						en: `Unable to display collections`,
						de: 'Karten können nicht angezeigt werden',
						uk: 'Не вдалося відобразити колекції'
					}) }}
				</template>
				
				<template v-slot:details>
					{{ state.error }}
				</template>

			</ErrorMessage>

			<LoadingMessage v-else-if="state.data === null">
				{{ intl(lang, {
					en: 'Loading...',
					de: 'Lädt...',
					uk: 'Один момент...'
				}) }}
			</LoadingMessage>

			<p v-else>
				{{ intl(lang, {
					en: `You haven't added any collections yet!`,
					de: 'Sie haben noch keine Karten hinzugefügt!',
					uk: 'Ви ще не додали жодної колекції!'
				}) }}
			</p>
			
		</CentralMessage>

		<CollectionBreak />

		<CollectionEndlistAction>

			<GenericButton theme="orange" @click="openExplore">
				{{ intl(lang, {
					en: 'Explore cards',
					de: 'Mehr Karten finden',
					uk: 'Знайти картки'
				}) }}
			</GenericButton>

		</CollectionEndlistAction>
	
	</AppUI>

</template>
