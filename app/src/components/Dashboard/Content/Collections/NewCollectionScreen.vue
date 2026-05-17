<script setup lang="ts">
import { computed, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useClient } from '../../../../api';
import AppUiHeader from '../../../App/AppUiHeader.vue';
import CollectionFormWrapper from './CollectionFormWrapper.vue';
import CollectionFormHeader from './CollectionFormHeader.vue';
import InputLabel from '../../../App/InputLabel.vue';
import GenericInput from '../../../App/GenericInput.vue';
import InlineErorrMessage from '../../../App/InlineErorrMessage.vue';
import GenericButton from '../../../App/GenericButton.vue';
import InputRow from '../../../App/InputRow.vue';
import { pickUploadFiles, type ContentBundle } from '../../../../content_io';
import ContentExporterStatus from './ContentExporterStatus.vue';

const client = useClient();
const router = useRouter();

const state = reactive({
	importer: {
		operation: null as string | null,
		active: false,
		progress: 0,
		warn: null as string | null,
		error: null as string | null
	},
	inputs: {
		name: '',
		description: ''
	},
	error: null as string | null,
});

const formValid = computed(() => state.inputs.name.trim().length > 0);

const backHref = '/app/dashboard/content';

const createCollection = async () => {

	const { data, error } = await client.collections.create(state.inputs);
	if (!data || error) {
		state.error = error?.message || 'Unable to create collection';
		return;
	}

	openCollection(data.id);
};

const uploadFile = async () => {

	const files = await pickUploadFiles({ accept: ['application/json'] });
	if (!files?.length) {
		return;
	}

	state.importer.active = true;
	state.importer.operation = 'Uploading bundle';

	try {
		const bundle: ContentBundle | null = await files[0].text().then(data => JSON.parse(data));
		await importBundledCollection(bundle);
	} catch (error) {
		state.importer.error = error instanceof Error ? error.message : 'Unable to parse bundle';
	}
};

const importBundledCollection = async (bundle: ContentBundle | null) => {

	if (bundle?.type !== 'collection_bundle') {
		state.importer.error = 'Unsupported bundle type';
		return;
	} else if (bundle.content.length === 0) {
		state.importer.error = 'Empty bundle';
		return;
	} else if (bundle.content.length > 1) {
		state.importer.warn = 'Bundle contains multiple collections, only the first one will be imorted';
	}

	const { meta, decks } = bundle.content[0];

	const { data, error } = await client.collections.create(meta);
	if (!data || error) {
		state.importer.error = error?.message || 'Unable to create collection';
		return;
	}

	const { id: collection_id } = data;

	state.importer.operation = 'Importing decks...';

	let imported = 0;

	for (const deck of decks) {

		if (!state.importer.active) {
			return;
		}
		
		const { data, error } = await client.decks.create({
			... deck.meta,
			collection_id,
			cards: deck.cards,
		});

		if (!data || error) {
			state.importer.warn = error?.message || 'Unable to import deck';
			continue;
		}

		imported++;
		state.importer.progress = decks.length / imported;
	}

	state.importer.progress = 1;

	if (imported !== decks.length) {

		if (decks.length === 0) {
			state.importer.error = 'Unable to load collection decks';
			state.importer.active = false;
			return;
		}

		const missingCount = decks.length - imported;
		state.importer.warn = `Unable to export ${missingCount} decks`;
	}

	state.importer.operation = 'Import done';

	openCollection(collection_id);
};

const openCollection = (id: string) => {
	router.push(`/app/dashboard/content/collection/${id}`);
};

</script>

<template>

	<AppUiHeader :backHref="backHref">
		<template v-slot:title>
			Create new collection
		</template>
	</AppUiHeader>

	<CollectionFormWrapper>

		<CollectionFormHeader>
			<template v-slot:overscript>
				New Collection
			</template>
			<template v-slot:title>
				What is this thing about?
			</template>
		</CollectionFormHeader>

		<InputLabel>

			<template v-slot:label>
				Collection name
			</template>

			<GenericInput type="text" variant="borderless" placeholder="Pick a passing collection name" v-model="state.inputs.name" />

		</InputLabel>

		<InputLabel>

			<template v-slot:label>
				Description
			</template>

			<GenericInput type="text" variant="borderless" placeholder="Describe the purpose of this collection" v-model="state.inputs.description" />

		</InputLabel>

		<InlineErorrMessage v-if="state.error">
			{{ state.error }}
		</InlineErorrMessage>

		<InputRow>
			<GenericButton :disabled="!formValid" @click="createCollection">
				Create collection →
			</GenericButton>
			<GenericButton theme="orange" variant="thin" :disabled="state.importer.active" @click="uploadFile">
				Upload file
			</GenericButton>
		</InputRow>

		<ContentExporterStatus v-if="state.importer.active"
			:operation="state.importer.operation"
			:progress="state.importer.progress"
			:warn="state.importer.warn"
			:error="state.importer.error"
			@cancel="state.importer.active = false" />

	</CollectionFormWrapper>

</template>
