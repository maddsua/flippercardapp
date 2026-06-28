<script setup lang="ts">
import { computed, reactive } from 'vue';
import { unwrapErrorMessage, useClient } from '@/api';
import type { CardDeckMeta, ResourceVisibility } from '@/api_models';
import { resourceVisibilityOptions } from '@/inputs';
import type { CardNode, ContentSummary } from '@/content';
import GenericButton from '@/components/App/Inputs/GenericButton.vue';
import GenericDropdown from '@/components/App/Inputs/GenericDropdown.vue';
import GenericInput from '@/components/App/Inputs/GenericInput.vue';
import InputLabel from '@/components/App/Inputs/InputLabel.vue';
import InlineErrorMessage from '@/components/App/Messages/InlineErrorMessage.vue';
import EditorModal from '../EditorModal.vue';
import GenericToggle from '@/components/App/Inputs/GenericToggle.vue';

interface Origin {
	deckID: string | null,
	collectionID: string | null,
};

interface Content {
	meta: ContentMeta;
	cards: CardNode[]
};

interface ContentMeta {
	summary: ContentSummary;
	visibility: ResourceVisibility;
};

interface ContentChanges {
	summary: boolean;
	cards: boolean;
};

const props = defineProps<{
	origin: Origin;
	content: Content;
	changes: ContentChanges;
}>();

const emit = defineEmits<{
	(e: 'done'): void;
	(e: 'publish', meta: CardDeckMeta): void;
}>();

const client = useClient();

const state = reactive({
	busy: false,
	meta: {
		name: props.content.meta.summary.name,
		description: props.content.meta.summary.description || '',
		visibility: props.content.meta.visibility,
	},
	editSummary: false,
	versionLabel: '',
	error: null as string | null,
});

const summaryChanged = computed(() =>
	props.changes.summary ||
	state.meta.name !== props.content.meta.summary.name ||
	state.meta.description !== props.content.meta.summary.description);

const publishVersion = async () => {

	if (!props.origin || state.busy) {
		state.error = 'Invalid origin data';
		return;
	}

	state.busy = true;
	state.error = null;

	const { deckID, collectionID } = props.origin;

	let meta: CardDeckMeta | null = null;

	if (deckID) {
		meta = await publishChanges(deckID);
	} else if (collectionID) {
		meta = await publishNew(collectionID);
	} else {
		state.error = 'Invalid editor state'
		state.busy = false;
	}

	if (meta) {
		emit('publish', meta);
	}

	if (!state.error) {
		emit('done');
	}

	state.busy = false;
};

const publishNew = async (collectionID: string): Promise<CardDeckMeta | null> => {

	const { data, error } = await client.decks.create({
		summary: {
			name: state.meta.name,
			description: state.meta.description,
		},
		visibility: state.meta.visibility,
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

const publishChanges = async (deckID: string): Promise<CardDeckMeta | null> => {

	const { data, error } = await client.decks.update(deckID, {
		summary: summaryChanged.value ? {
			name: state.meta.name,
			description: state.meta.description,
		} : null,
		visibility: state.meta.visibility,
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

			<template v-if="state.editSummary">

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

			</template>

			<InputLabel variant="slick">
				<template v-slot:label>
					Version label
				</template>
				<GenericInput placeholder="Label, e.g. 'New version'" v-model="state.versionLabel" />
			</InputLabel>

			<InputLabel variant="slick">
				<template v-slot:label>
					Deck visibility
				</template>
				<GenericDropdown :options="resourceVisibilityOptions" v-model="state.meta.visibility" />
			</InputLabel>

			<GenericToggle label="Edit summary" v-model="state.editSummary" />

			<GenericButton variant="thin-wide" :disabled="state.busy" :spinner="state.busy" @click="publishVersion">
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
		max-height: 100%;
		overflow: hidden auto;
		scrollbar-width: thin;
		padding-right: 0.5rem;
	}
</style>
