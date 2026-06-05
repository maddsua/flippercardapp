<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useClient } from '../../api';
import type { CardDeckMetadata } from '../../api_models';
import { intl, useLanguage } from '../../intl';
import { useStorage } from '../../storage/storage';
import AppUI from '../App/Layout/AppUI.vue';
import AppUiHeader from '../App/Layout/AppUiHeader.vue';
import CentralMessage from '../App/Messages/CentralMessage.vue';
import ErrorMessage from '../App/Messages/ErrorMessage.vue';
import LoadingMessage from '../App/Messages/LoadingMessage.vue';
import ContentList from '../Content/ContentList.vue';
import ContentListEntry from '../Content/ContentListEntry.vue';

interface Entry extends CardDeckMetadata {
	score: number;
};

const state = reactive({
	data: null as Entry[] | null,
	error: null as string | null
});

const router = useRouter();
const lang = useLanguage();
const store = useStorage();
const client = useClient();

onMounted(async () => {

	const ids = await store.decks.starred.all().catch(() => []);
	if (!ids.length) {
		state.data = [];
		return;
	}

	const { data, error } = await client.decks.list({ ids, limit: ids.length });
	if (!data || error) {
		state.error = error?.message || 'Unable to load decks';
		return;
	}

	const loadedDeckIDSet = new Set(data.entries.map(item => item.id));

	const scoreMap = new Map(await store.decks.stats
		.filter(val => loadedDeckIDSet.has(val.deck_id))
		.catch(() => [])
		.then(entries => entries.map(entry => [entry.deck_id, entry.score])));

	state.data = data.entries.map(item => ({ ... item, score: scoreMap.get(item.id) || 0 }));
});

const openDeck = (id: string) => {
	router.push(`/play/deck/${id}`);
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

		<ContentList v-if="state.data && state.data.length">
			<ContentListEntry v-for="item of state.data"
				:title="item.name"
				:summary="item.description"
				:visibility="item.visibility"
				:starred="true"
				:cardCount="item.size"
				:score="item.score"
				@click="openDeck(item.id)" />
		</ContentList>

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
