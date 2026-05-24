<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import { intl, useLanguage } from '../../intl';
import AppUI from '../App/AppUI.vue';
import AppUiHeader from '../App/AppUiHeader.vue';
import ContentList from '../Content/ContentList.vue';
import ContentListEntry from '../Content/ContentListEntry.vue';
import type { CardDeckMetadata } from '../../api_models';
import { useRouter } from 'vue-router';
import LoadingMessage from '../App/LoadingMessage.vue';
import ErrorMessage from '../App/ErrorMessage.vue';
import CentralMessage from '../App/CentralMessage.vue';
import { useStorage } from '../../storage';
import { useClient } from '../../api';

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

	const ids = await store.starredDecks.entries();
	if (!ids.length) {
		state.data = [];
		return;
	}

	const { data, error } = await client.decks.list({ ids, limit: ids.length });
	if (!data || error) {
		state.error = error?.message || 'Unable to load decks';
		return;
	}

	const playStats = new Map(await store.playStats.entries());
	state.data = data.entries.map(item => ({ ... item, score: playStats.get(item.id)?.score || 0 }));
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
