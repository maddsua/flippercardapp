<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import type { CardCollection } from '../../content';
import { useContent } from '../../content.loaders';
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

const router = useRouter();

const state = reactive({
	collections: [] as CardCollection[],
	ready: false,
	error: null as string | null,
});

onMounted(async () => {

	const { data, error } = await useContent().collections();
	if (!data || error) {
		state.error = error?.message || 'Unable to load collections';
		return;
	}

	state.collections = data;
	setTimeout(() => state.ready = true, 250);
});

const openCollection = (id: string) => {
	router.push(`/app/collection/${id}`);
};

const openExplore = () => {
	router.push('/app/explore');
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

		<CollectionList v-if="state.ready && state.collections.length">
			<CollectionListEntry v-for="item of state.collections" :title="item.name" @click="openCollection(item.id)" />
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

			<LoadingMessage v-else-if="!state.ready">
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
