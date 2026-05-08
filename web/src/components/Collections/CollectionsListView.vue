<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import type { CardCollection } from '../../content';
import { useCollectionProvider } from '../../content.loaders';
import { useRouter } from 'vue-router';
import CollectionList from './CollectionList.vue';
import CollectionListEntry from './CollectionListEntry.vue';
import FullscreenMessage from '../App/FullscreenMessage.vue';
import ErrorMessage from '../App/ErrorMessage.vue';
import CollectionHeader from './CollectionHeader.vue';
import CollectionContainer from './CollectionContainer.vue';
import CollectionBreak from './CollectionBreak.vue';
import CollectionEndlistAction from './CollectionEndlistAction.vue';
import GenericButton from '../App/GenericButton.vue';
import AppUI from '../App/AppUI.vue';
import { intl, useLanguage } from '../../intl';

const router = useRouter();

const state = reactive({
	collections: null as CardCollection[] | null,
	error: null as string | null,
});

onMounted(async () => {

	const { data, error } = await useCollectionProvider().collections();
	if (!data || error) {
		state.error = error?.message || 'Unable to load collections';
		return;
	}

	state.collections = data;
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

		<CollectionContainer>
	
			<CollectionHeader>
	
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
	
			</CollectionHeader>
	
			<CollectionList v-if="state.collections?.length">
				<CollectionListEntry v-for="item of state.collections" :title="item.name" @click="openCollection(item.id)" />
			</CollectionList>
	
			<FullscreenMessage v-else>
	
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
	
				<p v-else>
					{{ intl(lang, {
						en: `You haven't added any collections yet!`,
						de: 'Sie haben noch keine Karten hinzugefügt!',
						uk: 'Ви ще не додали жодної колекції!'
					}) }}
				</p>
				
			</FullscreenMessage>
	
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
	
		</CollectionContainer>
	</AppUI>

</template>
