<script setup lang="ts">
import { computed, onMounted, reactive } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import type { CollectionMeta, ResourceVisibility } from '@/api_models';
import { resourceVisibilityOptions } from '@/inputs';
import { useClient } from '@/api';
import AppUI from '@/components/App/Layout/AppUI.vue';
import AppUiHeader from '@/components/App/Layout/AppUiHeader.vue';
import FullscreenMessage from '@/components/App/Messages/FullscreenMessage.vue';
import GenericButton from '@/components/App/Inputs/GenericButton.vue';
import GenericDropdown from '@/components/App/Inputs/GenericDropdown.vue';
import GenericInput from '@/components/App/Inputs/GenericInput.vue';
import InlineErrorMessage from '@/components/App/Messages/InlineErrorMessage.vue';
import InputLabel from '@/components/App/Inputs/InputLabel.vue';
import InputRow from '@/components/App/Inputs/InputRow.vue';
import LoadingMessage from '@/components/App/Messages/LoadingMessage.vue';
import Skeleton from '@/components/App/Messages/Skeleton.vue';
import CollectionFormHeader from './CollectionFormHeader.vue';
import CollectionFormWrapper from './CollectionFormWrapper.vue';
import OverlayErrorMessage from '@/components/App/Messages/OverlayErrorMessage.vue';

const route = useRoute();
const router = useRouter();
const client = useClient();

const state = reactive({
	data: null as CollectionMeta | null,
	dataValid: false,
	inputs: {
		name: '',
		description: '',
		visibility: 'HIDDEN' as ResourceVisibility,
	},
	exporter: {
		busy: false,
		error: null as string | null,
	},
	error: null as string | null,
});

const backHref = computed(() => {
	if (state.data && state.dataValid) {
		return `/collection/${state.data.id}`;
	}
	return '/collections/all';
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
		visibility: collection.visibility,
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

	<AppUI>

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

		<OverlayErrorMessage v-if="state.error">

			Unable to load collection

			<template v-slot:details>
				{{ state.error }}
			</template>

			<template v-slot:after>
				<GenericButton variant="thin" @click="router.push(backHref)">
					Go back
				</GenericButton>
			</template>

		</OverlayErrorMessage>

		<FullscreenMessage v-else-if="!state.data">
			<LoadingMessage>
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

			<InputLabel>

				<template v-slot:label>
					Visibility
				</template>

				<GenericDropdown :options="resourceVisibilityOptions" v-model="state.inputs.visibility" />

			</InputLabel>

			<InlineErrorMessage v-if="state.error">
				{{ state.error }}
			</InlineErrorMessage>

			<InputRow>
				<GenericButton :disabled="!formValid" @click="updateCollection">
					Update collection →
				</GenericButton>
			</InputRow>

			<InputRow>
				<GenericButton variant="thin" theme="red" @click="deleteCollection">
					✗ Delete collection
				</GenericButton>
			</InputRow>

		</CollectionFormWrapper>

	</AppUI>

</template>
