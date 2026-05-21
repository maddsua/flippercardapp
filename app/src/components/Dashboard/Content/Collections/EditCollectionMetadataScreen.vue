<script setup lang="ts">
import { computed, onMounted, reactive } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useClient } from '../../../../api';
import type { CollectionMetadata } from '../../../../api_models';
import AppUiHeader from '../../../App/AppUiHeader.vue';
import Skeleton from '../../../App/Skeleton.vue';
import CollectionFormWrapper from './CollectionFormWrapper.vue';
import CollectionFormHeader from './CollectionFormHeader.vue';
import InputLabel from '../../../App/InputLabel.vue';
import GenericInput from '../../../App/GenericInput.vue';
import InlineErorrMessage from '../../../App/InlineErorrMessage.vue';
import GenericButton from '../../../App/GenericButton.vue';
import ErrorMessage from '../../../App/ErrorMessage.vue';
import LoadingMessage from '../../../App/LoadingMessage.vue';
import FullscreenMessage from '../../../App/FullscreenMessage.vue';
import InputRow from '../../../App/InputRow.vue';
import { downloadBlob } from '../../../../files';

const route = useRoute();
const router = useRouter();
const client = useClient();


const state = reactive({
	data: null as CollectionMetadata | null,
	dataValid: false,
	inputs: {
		name: '',
		description: ''
	},
	exporter: {
		busy: false,
		error: null as string | null,
	},
	error: null as string | null,
});

const backHref = computed(() => {
	if (state.data && state.dataValid) {
		return `/app/dashboard/content/collection/${state.data.id}`;
	}
	return '/app/dashboard/content';
})

const formValid = computed(() => state.inputs.name.trim().length > 0);

onMounted(async () => {

	const { collection_id } = route.params;
 
	const { data, error } = await client.collections.list({ ids: Array.isArray(collection_id) ? collection_id : [collection_id] })
	if ( error) {
		state.error = error.message;
		return;
	}

	const collection = data?.entries.find(item => item.id === collection_id);
	if (!collection) {
		state.error = 'Collection not found';
		return;
	}

	state.data = collection;
	state.dataValid = true;

	state.inputs = {
		name: collection.name,
		description: collection.description || '',
	};
});

const updateCollection = async () => {

	if (!state.data) {
		throw new Error('Invalid condition');
	}

	const { data, error } = await client.collections.update(state.data.id, state.inputs);
	if (!data || error) {
		state.error = error?.message || 'Unable to update collection';
		return;
	}

	router.push(backHref.value);
};

const exportCollection = async () => {

	if (!state.data) {
		return;
	}

	state.exporter = { busy: true, error: null };

	const { blob, error } = await client.collections.exportBundle(state.data.id);
	if (!blob || error) {
		state.exporter.busy = false;
		state.error = error?.message || 'Unable to export collection';
		return;
	}

	downloadBlob(blob, `${state.data.name}-export-${new Date().getTime()}.cardbundle`);

	state.exporter.busy = false;
};

const deleteCollection = async () => {

	if (!state.data) {
		throw new Error('Invalid condition');
	}

	if (state.data.size > 0) {

		const confirmText = 'Delete collection and content';
		while (true) {

			const response = prompt(`Type '${confirmText}' to proceed`)?.trim();
			if (!response) {
				return;
			}

			if (response.toLowerCase() === confirmText.toLowerCase()) {
				break;
			}
		}

	} else if (!confirm('Really delete this collection?')) {
		return;
	}

	const { error } = await client.collections.remove(state.data.id, { recursive: state.data.size > 0 });
	if (error) {
		state.error = error?.message || 'Unable to delete collection';
		return;	
	}

	state.dataValid = false;

	router.push(backHref.value);
};

</script>

<template>

	<AppUiHeader :backHref="backHref">
		<template v-slot:title>
			Edit collection
		</template>
		<template v-slot:summary>
			<template v-if="state.data">
				Editing: {{ state.data.name }}
			</template>
			<Skeleton v-else-if="!state.error">
				Loading collection
			</Skeleton>
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

	<CollectionFormWrapper v-else>

		<CollectionFormHeader>
			<template v-slot:overscript>
				Edit Collection
			</template>
			<template v-slot:title>
				Changing some things huh?
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
			<GenericButton theme="green" :disabled="!formValid" @click="updateCollection">
				Update collection →
			</GenericButton>
		</InputRow>

		<InputRow>
			<GenericButton variant="thin" :spinner="state.exporter.busy" :disabled="!state.data || state.exporter.busy" @click="exportCollection">
				Export collection bundle
			</GenericButton>
			<GenericButton  variant="thin" theme="red" @click="deleteCollection">
				✗ Delete collection
			</GenericButton>
		</InputRow>

	</CollectionFormWrapper>
	
</template>
