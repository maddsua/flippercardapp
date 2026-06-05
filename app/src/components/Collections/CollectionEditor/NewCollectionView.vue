<script setup lang="ts">
import { computed, reactive } from 'vue';
import { useRouter } from 'vue-router';
import type { ResourceVisibility } from '@/api_models';
import { resourceVisibilityOptions } from '@/inputs';
import { useClient } from '@/api';
import AppUI from '@/components/App/Layout/AppUI.vue';
import AppUiHeader from '@/components/App/Layout/AppUiHeader.vue';
import GenericButton from '@/components/App/Inputs/GenericButton.vue';
import GenericDropdown from '@/components/App/Inputs/GenericDropdown.vue';
import GenericInput from '@/components/App/Inputs/GenericInput.vue';
import InlineErorrMessage from '@/components/App/Messages/InlineErorrMessage.vue';
import InputLabel from '@/components/App/Inputs/InputLabel.vue';
import InputRow from '@/components/App/Inputs/InputRow.vue';
import CollectionFormHeader from './CollectionFormHeader.vue';
import CollectionFormWrapper from './CollectionFormWrapper.vue';

const client = useClient();
const router = useRouter();

const state = reactive({
	inputs: {
		name: '',
		description: '',
		visibility: 'HIDDEN' as ResourceVisibility,
	},
	error: null as string | null,
});

const formValid = computed(() => state.inputs.name.trim().length > 0);

const backHref = '/collections/all';

const createCollection = async () => {

	const { data, error } = await client.collections.create(state.inputs);
	if (!data || error) {
		state.error = error?.message || 'Unable to create collection';
		return;
	}

	openCollection(data.id);
};

const openCollection = (id: string) => {
	router.push(`/collection/${id}`);
};

</script>

<template>

	<AppUI>

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
	
			<InputLabel>
	
				<template v-slot:label>
					Visibility
				</template>
	
				<GenericDropdown :options="resourceVisibilityOptions" v-model="state.inputs.visibility" />
	
			</InputLabel>
	
			<InlineErorrMessage v-if="state.error">
				{{ state.error }}
			</InlineErorrMessage>
	
			<InputRow>
				<GenericButton :disabled="!formValid" @click="createCollection">
					Create collection →
				</GenericButton>
			</InputRow>
	
		</CollectionFormWrapper>

	</AppUI>

</template>
