<script setup lang="ts">
import type { CardDeckMetadata, ResourceVisibility } from '@/api_models';
import type { CardNode } from '@/content';
import EditorModal from '../EditorModal.vue';
import { unwrapErrorMessage, useClient } from '@/api';
import { computed, reactive } from 'vue';
import GenericButton from '@/components/App/Inputs/GenericButton.vue';
import GenericInput from '@/components/App/Inputs/GenericInput.vue';
import InlineErrorMessage from '@/components/App/Messages/InlineErrorMessage.vue';
import { resourceVisibilityOptions } from '@/inputs';
import GenericDropdown from '@/components/App/Inputs/GenericDropdown.vue';
import InputLabel from '@/components/App/Inputs/InputLabel.vue';

interface Origin {
	deckID: string | null,
	collectionID: string | null,
};

interface Content {
	meta: ContentMeta;
	cards: CardNode[]
};

interface ContentMeta {
	name: string;
	description: string | null;
	visibility: ResourceVisibility;
};

interface ContentChanges {
	meta: boolean;
	cards: boolean;
};

const props = defineProps<{
	origin: Origin;
	content: Content;
	changes: ContentChanges;
}>();

const emit = defineEmits<{
	(e: 'done'): void;
	(e: 'update', meta: CardDeckMetadata): void;
}>();

const client = useClient();

const state = reactive({
	busy: false,
	meta: {
		name: props.content.meta.name,
		description: props.content.meta.description || '',
		visibility: props.content.meta.visibility,
	},
	versionLabel: '',
	error: null as string | null,
});

const metaChanged = computed(() =>
	props.changes.meta ||
	state.meta.name !== props.content.meta.name ||
	state.meta.description !== props.content.meta.description ||
	state.meta.visibility !== props.content.meta.visibility);

const publishVersion = async () => {

	if (!props.origin || state.busy) {
		state.error = 'Invalid origin data';
		return;
	}

	state.busy = true;
	state.error = null;

	const { deckID, collectionID } = props.origin;

	let meta: CardDeckMetadata | null = null;

	if (deckID) {
		meta = await publishChanges(deckID);
	} else if (collectionID) {
		meta = await publishNew(collectionID);
	} else {
		state.error = 'Invalid editor state'
		state.busy = false;
	}

	if (meta) {
		emit('update', meta);
	}

	if (!state.error) {
		emit('done');
	}

	state.busy = false;
};

const publishNew = async (collectionID: string): Promise<CardDeckMetadata | null> => {

	const { data, error } = await client.decks.create({
		meta: state.meta,
		content: { cards: props.content.cards },
		collection_id: collectionID,
		label: state.versionLabel || null,
	});

	if (!data || error) {
		state.error = unwrapErrorMessage(error);
		state.busy = false;
		return null;
	}

	return data;
};

const publishChanges = async (deckID: string): Promise<CardDeckMetadata | null> => {

	const { data, error } = await client.decks.update(deckID, {
		meta: metaChanged.value ? state.meta : null,
		content: props.changes.cards ? { cards: props.content.cards } : null,
		label: state.versionLabel || null,
	});

	if (!data || error) {
		state.error = unwrapErrorMessage(error);
		state.busy = false;
		return null;
	}

	return data;
};

</script>

<template>
	<EditorModal title="Publish version" variant="compact" @close="emit('done')">

		<div class="publish-form">

			<InputLabel variant="slick">
				<template v-slot:label>
					Deck name
				</template>
				<GenericInput placeholder="Deck name" v-model="state.meta.name" />
			</InputLabel>

			<InputLabel variant="slick">
				<template v-slot:label>
					Deck summary
				</template>
				<GenericInput placeholder="Deck summary" :multiline="true" v-model="state.meta.description" />
			</InputLabel>

			<InputLabel variant="slick">
				<template v-slot:label>
					Deck visibility
				</template>
				<GenericDropdown :options="resourceVisibilityOptions" v-model="state.meta.visibility" />
			</InputLabel>

			<InputLabel variant="slick">
				<template v-slot:label>
					Version label
				</template>
				<GenericInput placeholder="Label, e.g. 'New version'" v-model="state.versionLabel" />
			</InputLabel>
	
			<GenericButton variant="thin" :disabled="state.busy" :spinner="state.busy" @click="publishVersion">
				Publish
			</GenericButton>

			<InlineErrorMessage v-if="state.error">
				<template v-slot:title>
					Unable to publish changes
				</template>
				{{ state.error }}
			</InlineErrorMessage>

		</div>

	</EditorModal>
</template>

<style lang="scss" scoped>
	.publish-form {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}
</style>
