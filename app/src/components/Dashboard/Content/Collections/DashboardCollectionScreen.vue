<script setup lang="ts">
import {  onMounted, reactive } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import AppUiHeader from '../../../App/AppUiHeader.vue';
import Skeleton from '../../../App/Skeleton.vue';
import ErrorMessage from '../../../App/ErrorMessage.vue';
import LoadingMessage from '../../../App/LoadingMessage.vue';
import FullscreenMessage from '../../../App/FullscreenMessage.vue';
import { useClient } from '../../../../api';
import type { Collection } from '../../../../api_models';
import GenericButton from '../../../App/GenericButton.vue';
import CentralMessage from '../../../App/CentralMessage.vue';
import DeckEntry from './DeckEntry.vue';

const route = useRoute();
const router = useRouter();
const client = useClient();

const backHref = '/dashboard/content';

const state = reactive({
	data: null as Collection | null,
	error: null as string | null,
});

onMounted(async () => {

	const { collection_id } = route.params;

	const id = typeof collection_id === 'string' ? collection_id : collection_id[0];

	const { data, error } = await client.collections.load(id);
	if (!data || error) {
		state.error = error?.message || 'Unable to load collection';
		return;
	}

	state.data = data;
});

const editCollectionMetadata = () => {
	router.push(`/dashboard/content/collection/${state.data!.id}/metadata`);
};

const openDeckEditor = (params: { deckID?: string; collectionID?: string }) => {
	if (params.deckID) {
		router.push(`/editor/deck/${params.deckID}/editor`);
	} else if (params.collectionID) {
		router.push(`/editor/deck/editor?collection_id=${params.collectionID}`);
	}
};

const deleteDeck = async (deckID: string) => {

	if (!state.data) {
		throw new Error('Invalid state');
	}

	if (!confirm('Delete deck?')) {
		return;
	}

	const { error } = await client.decks.remove(deckID);
	if (error) {
		console.error('Unable to delete collection deck:', error.message);
		return;
	}

	state.data.decks = state.data.decks?.filter(item => item.id !== deckID) || null;
};

</script>

<template>

	<AppUiHeader :backHref="backHref">
		<template v-slot:title>
			Manage collection
		</template>
	</AppUiHeader>

	<FullscreenMessage v-if="!state.data">

		<ErrorMessage v-if="state.error">
			<template v-slot:message>
				Unable to load collection
			</template>
			<template v-slot:details>
				{{ state.error }}
			</template>
		</ErrorMessage>
	
		<LoadingMessage v-else>
			Loading metadata...
		</LoadingMessage>

	</FullscreenMessage>

	<template v-else>

		<div class="metadata-section">

			<div class="content">
				<div class="title">
					<template v-if="state.data">
						{{ state.data.name }}
					</template>
					<Skeleton v-else>
						Loading name...
					</Skeleton>
				</div>
				<div class="description">
					<template v-if="state.data">
						{{ state.data.description || '[No description set]' }}
					</template>
					<Skeleton v-else-if="!state.error">
						Loading description...
					</Skeleton>
				</div>
			</div>

			<div class="actions">
				<GenericButton variant="thin" :disabled="!state.data" @click="editCollectionMetadata">
					Manage
				</GenericButton>
			</div>
		</div>

		<div class="decks-section">

			<div class="header">
				<div class="title">
					Collection's decks
				</div>
				<div class="actions">
					<GenericButton variant="thin" theme="green" :disabled="!state.data" @click="openDeckEditor({ collectionID: state.data.id })">
						New deck
					</GenericButton>
				</div>
			</div>

			<div v-if="state.data" class="deck-list">

				<CentralMessage v-if="!state.data.decks.length">
					No decks in this collection
				</CentralMessage>

				<template v-else>
					<DeckEntry v-for="entry of state.data.decks"
						:entry="entry"
						@edit="openDeckEditor({ deckID: entry.id })"
						@delete="deleteDeck(entry.id)" />
				</template>

			</div>

		</div>

	</template>
	
</template>

<style lang="scss" scoped>

	.metadata-section {
		display: flex;
		flex-flow: row nowrap;
		gap: 1rem;
		justify-content: space-between;
		padding: 0.5rem 1rem;
		background-color: var(--app-theme-ghostly-glow);
		border-radius: 0.5rem;

		.content {
			display: flex;
			flex-direction: column;
			gap: 0.25rem;
			min-width: 0;

			.title {
				font-size: 0.95rem;
				font-weight: 600;
			}

			.description {
				font-size: 0.75rem;
				color: var(--app-theme-mysterious-white);
			}

			.title, .description {
				max-width: 100%;
				overflow: hidden;
				white-space: nowrap;
				text-overflow: ellipsis;
			}
		}

		.actions {
			display: flex;
			flex-flow: row nowrap;
			gap: 0.5rem;
			align-items: center;
			flex-shrink: 0;
		}
	}

	.decks-section {
		display: flex;
		flex-direction: column;
		gap: 2rem;
		padding: 2rem 0;

		.header {
			display: flex;
			flex-flow: row nowrap;
			gap: 1rem;
			justify-content: space-between;

			.title {
				font-size: 0.95rem;
				font-weight: 600;
			}

			.actions {
				display: flex;
				flex-flow: row nowrap;
				gap: 0.5rem;
				align-items: center;
				flex-shrink: 0;
			}
		}

		.deck-list {
			display: flex;
			flex-direction: column;
			gap: 1rem;
		}
	}

</style>
