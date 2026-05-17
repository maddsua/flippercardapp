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
import { downloadFile, pickUploadFiles, type CollectionBundleContent, type CollectionBundle, type ContentBundle, type CollectionBundleDeckContent } from '../../../../content_io';
import ContentExporterStatus from './ContentExporterStatus.vue';

const route = useRoute();
const router = useRouter();
const client = useClient();

const backHref = '/app/dashboard/content';

const state = reactive({
	data: null as Collection | null,
	error: null as string | null,

	contentIO: {
		operation: null as string | null,
		active: false,
		progress: 0,
		warn: null as string | null,
		error: null as string | null
	},
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
	router.push(`/app/dashboard/content/collection/${state.data!.id}/metadata`);
};

const openDeckEditor = (params: { deckID?: string; collectionID?: string }) => {
	if (params.deckID) {
		router.push(`/app/editor/deck/${params.deckID}/editor`);
	} else if (params.collectionID) {
		router.push(`/app/editor/deck/editor?collection_id=${params.collectionID}`);
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

const exportCollection = async () => {

	if (!state.data) {
		throw new Error('Logic error');
	}

	state.contentIO.active = true;
	state.contentIO.operation = 'Exporting collection...';
	state.contentIO.warn = null;
	state.contentIO.error = null;

	const collection: CollectionBundleContent = {
		meta: {
			name: state.data.name,
			description: state.data.description,
			created: state.data.created,
			updated: state.data.updated
		},
		decks: []
	};

	for (const entry of state.data.decks) {

		if (!state.contentIO.active) {
			return;
		}

		const { data, error } = await client.decks.load(entry.id);
		if (!data || error) {
			state.contentIO.warn = error?.message || 'Unable to export deck';
			continue;
		}

		collection.decks.push({
			meta: {
				name: entry.name,
				description: entry.description,
			},
			cards: data.cards,
		});

		state.contentIO.progress = collection.decks.length / state.data.decks.length;
	}

	state.contentIO.progress = 1;

	if (state.data.decks.length !== collection.decks.length) {

		if (collection.decks.length === 0) {
			state.contentIO.error = 'Unable to load collection decks';
			state.contentIO.active = false;
			return;
		}

		const missingCount = state.data.decks.length - collection.decks.length;
		state.contentIO.warn = `Unable to export ${missingCount} decks`;
	}

	const bundle: CollectionBundle = {
		type: 'collection_bundle',
		content: [collection]
	};

	const name = state.data.name.replace(/[^a-z0-9]/gi, '_');

	downloadFile(JSON.stringify(bundle), `${name}-export.json`);

	state.contentIO.operation = 'Export done';

	setTimeout(() => state.contentIO.active = false, 3_000);
};

const importContent = async () => {

	const files = await pickUploadFiles();
	if (!files?.length) {
		return;
	}

	state.contentIO.active = true;
	state.contentIO.operation = 'Parsing bundles...';

	for (const file of files) {
		try {
			const bundle: ContentBundle | null = await file.text().then(data => JSON.parse(data));
			await importContentBundle(bundle);
		} catch (error) {
			state.contentIO.error = error instanceof Error ? error.message : 'Unable to parse bundle';
		}
	}

	if (!state.contentIO.error) {
		setTimeout(() => state.contentIO.active = false, 3_000);
	}
};

const importContentBundle = async (bundle: ContentBundle | null) => {

	let decks: CollectionBundleDeckContent[] = [];

	switch (bundle?.type) {

		case 'collection_bundle':

			if (!confirm(`Are you sure you want to import content from ${bundle.content.length} bundles?`)) {
				state.contentIO.error = 'Import cancelled by user';
				return;
			}

			decks = bundle.content.map(item => item.decks).flat();
			break;

		case 'deck_bundle':
			decks = bundle.content;
			break;

		default:
			state.contentIO.error = 'Unsupported bundle file';
			return;
	}

	state.contentIO.operation = 'Importing decks...';

	let imported = 0;

	for (const deck of decks) {

		if (!state.contentIO.active) {
			return;
		}
		
		const { data, error } = await client.decks.create({
			... deck.meta,
			collection_id: state.data!.id,
			cards: deck.cards,
		});

		if (!data || error) {
			state.contentIO.warn = error?.message || 'Unable to import deck';
			continue;
		}

		imported++;
		state.contentIO.progress = decks.length / imported;

		state.data?.decks.push(data);
	}

	state.contentIO.progress = 1;

	if (imported !== decks.length) {

		if (decks.length === 0) {
			state.contentIO.error = 'Unable to load collection decks';
			state.contentIO.active = false;
			return;
		}

		const missingCount = decks.length - imported;
		state.contentIO.warn = `Unable to export ${missingCount} decks`;
	}

	state.contentIO.operation = 'Import done';
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
				<GenericButton variant="thin" theme="red" :disabled="!state.data" @click="editCollectionMetadata">
					Manage
				</GenericButton>
				<GenericButton variant="thin" :disabled="!state.data || state.contentIO.active" @click="exportCollection">
					Export
				</GenericButton>
			</div>
		</div>

		<div class="decks-section">

			<div class="header">
				<div class="title">
					Collection's decks
				</div>
				<div class="actions">
					<GenericButton variant="thin" :disabled="!state.data" @click="openDeckEditor({ collectionID: state.data.id })">
						Add deck
					</GenericButton>
					<GenericButton variant="thin" theme="green" :disabled="!state.data" @click="importContent">
						Upload decks
					</GenericButton>
				</div>
			</div>

			<ContentExporterStatus v-if="state.contentIO.active"
				:operation="state.contentIO.operation"
				:progress="state.contentIO.progress"
				:warn="state.contentIO.warn"
				:error="state.contentIO.error"
				@cancel="state.contentIO.active = false" />

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
