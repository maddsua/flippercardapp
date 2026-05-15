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

const client = useClient();
const router = useRouter();

const state = reactive({
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

	router.push(backHref);
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
		</InputRow>

	</CollectionFormWrapper>

</template>
