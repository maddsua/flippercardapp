<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useClient } from '../../api';
import type { CollectionMetadata } from '../../api_models';
import { intl, useLanguage } from '../../intl';
import { useStorage } from '../../storage/storage';
import AppUI from '../App/Layout/AppUI.vue';
import AppUiHeader from '../App/Layout/AppUiHeader.vue';
import CentralMessage from '../App/Messages/CentralMessage.vue';
import GenericButton from '../App/Inputs/GenericButton.vue';
import LoadingMessage from '../App/Messages/LoadingMessage.vue';
import ContentList from '../Content/ContentList.vue';
import ContentListEntry from '../Content/ContentListEntry.vue';
import CollectionBreak from './CollectionBreak.vue';
import CollectionEndlistAction from './CollectionEndlistAction.vue';
import InlineErrorMessage from '../App/Messages/InlineErrorMessage.vue';
import { distributeCollectionPlayScore } from '@/play';

const router = useRouter();
const client = useClient();
const store = useStorage();

interface Entry extends CollectionMetadata {
	score: number;
};

const state = reactive({
	data: null as Entry[] | null,
	error: null as string | null,
});

onMounted(async () => {

	const ids = await store.collections.starred.all().catch(() => []);
	if (!ids.length) {
		state.data = [];
		return;
	}

	const { data, error } = await client.collections.list({ ids, limit: ids.length });
	if (!data || error) {
		state.error = error?.message || 'Unable to load collections';
		return;
	}

	const collectionStats = new Map(await store.collections.stats.aggregated(data.entries.map(item => item.id)).catch(() => []));

	state.data = data.entries.map(item => ({
		...item,
		score: distributeCollectionPlayScore(collectionStats, item),
	}));
});

const openCollection = (id: string) => {
	router.push(`/collection/${id}`);
};

const openCollectionDiscover = () => {
	router.push('/collections/discover');
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

		<ContentList v-if="state.data !== null && state.data.length">
			<ContentListEntry v-for="item of state.data"
				:title="item.name"
				:summary="item.description"
				:visibility="item.visibility"
				:date="item.created"
				:starred="true"
				:deckCount="item.size"
				:score="item.score"
				@click="openCollection(item.id)" />
		</ContentList>

		<CentralMessage v-else>

			<InlineErrorMessage v-if="state.error">

				<template v-slot:title>
					{{ intl(lang, {
						en: `Unable to display collections`,
						de: 'Karten können nicht angezeigt werden',
						uk: 'Не вдалося відобразити колекції'
					}) }}
				</template>
				
				{{ state.error }}

			</InlineErrorMessage>

			<LoadingMessage v-else-if="state.data === null">
				{{ intl(lang, {
					en: 'Loading...',
					de: 'Lädt...',
					uk: 'Один момент...'
				}) }}
			</LoadingMessage>

			<div v-else class="welcome-message">

				<img class="demo" src="/src/assets/images/card-stack-demo.svg" width="720" height="720" alt="Card demo image" />

				<p>
					{{ intl(lang, {
						en: `You haven't added any collections yet!`,
						de: 'Sie haben noch keine Karten hinzugefügt!',
						uk: 'Ви ще не додали жодної колекції!'
					}) }}
				</p>

			</div>

		</CentralMessage>

		<CollectionBreak />

		<CollectionEndlistAction>

			<GenericButton theme="orange" variant="thin" @click="openCollectionDiscover">
				{{ intl(lang, {
					en: 'Explore cards',
					de: 'Mehr Karten finden',
					uk: 'Знайти картки'
				}) }}
			</GenericButton>

		</CollectionEndlistAction>
	
	</AppUI>

</template>

<style lang="scss" scoped>

	.welcome-message {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.75rem;

		.demo {
			display: block;
			width: 10rem;
			max-width: 100%;
			height: auto;
			opacity: 0.75;
		}
	}

</style>
