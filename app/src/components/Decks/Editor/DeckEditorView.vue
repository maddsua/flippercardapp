<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, reactive, toRaw, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { unwrapErrorMessage, useClient } from '@/api';
import type { CardDeckMetadata, ResourceVisibility } from '@/api_models';
import type { CardNode } from '@/content';
import { useStorage, type DeckEditorHistoryMetaEntry } from '@/storage/storage';
import GenericButton from '@/components/App/Inputs/GenericButton.vue';
import LoadingMessage from '@/components/App/Messages/LoadingMessage.vue';
import EditorContentExporter from './EditorModals/EditorContentExporter.vue';
import EditorContentImporter from './EditorModals/EditorContentImporter.vue';
import DeckEditorStatusBar from './DeckEditorStatusBar.vue';
import DeckCardList from './EditorCardNavigationList.vue';
import EditorScreenOverlay from './EditorScreenOverlay.vue';
import EditorVersionControl from './EditorModals/EditorVersionControl.vue';
import OverlayErrorMessage from '@/components/App/Messages/OverlayErrorMessage.vue';
import CardFaceEditor from './FaceEditor/CardFaceEditor.vue';
import EditorVersionPublishModal from './EditorModals/EditorVersionPublishModal.vue';

const route = useRoute();
const router = useRouter();
const store = useStorage();
const client = useClient();

const defaultDeckMeta = () => ({
	name: 'Unnamed deck',
	description: null as string | null,
	visibility: 'HIDDEN' as ResourceVisibility,
});

const state = reactive({
	
	content: {
		meta: defaultDeckMeta(),
		cards: [] as CardNode[],
	},

	origin: {
		deckID: null as string | null,
		collectionID: null as string | null,
	},

	editor: {
		ready: false,
		error: null as string | null,
		prevAppTitle: null as string | null,
		view: {
			cardIdx: 0,
		},
		modals: {
			versions: false,
			importer: false,
			exporter: false,
			publish: false,
		},
		changes: {
			meta: false,
			cards: false,
		},
		snapshots: {
			timer: null as NodeJS.Timeout | null,
			lock: false,
			writtenVersion: null as DeckEditorHistoryMetaEntry | null,
			loadedVersion: null as DeckEditorHistoryMetaEntry | null,
		},
	},
});

const editorReady = computed(() => state.editor.ready && !state.editor.error);
const contentEdited = computed(() => !!state.content.cards.length && (state.editor.changes.cards || state.editor.changes.meta));
const changesSaved = computed(() => !state.editor.snapshots.timer && !!(state.editor.snapshots.loadedVersion || state.editor.snapshots.writtenVersion));

const activeCard = computed(() => {

	const entry = state.content.cards[state.editor.view.cardIdx];
	if (!entry) {
		return { id: null, front: null, back: null };
	}

	return { id: entry.id, front: entry.front, back: entry.back };
});

const cardSelectorList = computed(() => state.content.cards.map(item => item.front));

const addCard = (node: CardNode) => {
	state.content.cards.push(node);
	state.editor.view.cardIdx = state.content.cards.length - 1;
};

const createCard = () => {
	const newFace = () => ({ content: [] });
	addCard({ id: crypto.randomUUID(), front: newFace(), back: newFace() });
};

const selectCard = (idx: number) => {
	state.editor.view.cardIdx = idx;
};

const duplicateCard = (idx: number) => {
	const clonedNode = state.content.cards[idx];
	if (!clonedNode) {
		return;
	}
	addCard(structuredClone(toRaw(clonedNode)));
};

const removeCard = (idx: number) => {

	if (!confirm('Remove card?')) {
		return;
	}

	state.content.cards.splice(idx, 1);

	if (state.editor.view.cardIdx >= state.content.cards.length) {
		state.editor.view.cardIdx = state.content.cards.length - 1;
	}
};

interface ResumableState extends DeckEditorHistoryMetaEntry {
	editor: {
		view: Pick<typeof state['editor']['view'], 'cardIdx'>;
	};
	publisher: Pick<typeof state['origin'], 'deckID' | 'collectionID'>;
	content: typeof state['content'];
};

const queueStateSnapshot = () => {

	if (!state.editor.ready || state.editor.snapshots.lock) {
		return;
	}

	if (state.editor.snapshots.timer) {
		clearTimeout(state.editor.snapshots.timer);
	}

	const postSnapshot = () => {
		state.editor.snapshots.lock = false;
		state.editor.snapshots.timer = null;
	};

	const interval = 3_000;

	state.editor.snapshots.timer = setTimeout(async () => {

		if (state.editor.snapshots.lock) {
			return;
		}

		state.editor.snapshots.lock = true;

		const snapshot: ResumableState = {
			deck_id: state.origin?.deckID || 'new',
			timestamp: new Date(),
			content: structuredClone(toRaw(state.content)),
			editor: {
				view: {
					cardIdx: state.editor.view.cardIdx,
				},
			},
			publisher: {
				deckID: state.origin?.deckID || null,
				collectionID: state.origin?.collectionID || null,
			},
		};

		const { error } = await store.decks.editor.snapshots.push(snapshot)
			.then(() => ({ error: null }))
			.catch(error => ({ error }));

		if (error) {
			console.error('editor.snapshots.push', error);
			postSnapshot();
			return;
		}

		state.editor.snapshots.writtenVersion = { deck_id: snapshot.deck_id, timestamp: snapshot.timestamp };
		postSnapshot();

	}, interval);
};

const restoreStateSnapshot = async (deckID?: string) => {

	const snapshot = await store.decks.editor.snapshots.latest<ResumableState>(deckID || 'new').catch(() => null);
	if (!snapshot) {
		state.editor.snapshots.loadedVersion = null;
		return;
	}

	state.origin = { deckID: snapshot.publisher.deckID, collectionID: snapshot.publisher.collectionID };
	state.content = snapshot.content;
	state.editor.view.cardIdx = snapshot.editor.view.cardIdx;
	state.editor.changes = { meta: true, cards: true };

	state.editor.snapshots.loadedVersion = { deck_id: snapshot.deck_id, timestamp: snapshot.timestamp };
};

const clearStateSnapshot = async () => {

	const { ready } = state.editor;

	state.editor.ready = false;

	await store.decks.editor.snapshots.remove(state.origin.deckID || 'new')
		.catch(error => console.error('editor.snapshots.remove', error));

	if (state.editor.snapshots.timer) {
		clearTimeout(state.editor.snapshots.timer);
	}

	state.editor.snapshots.writtenVersion = null;
	state.editor.snapshots.loadedVersion = null;

	nextTick(() => state.editor.ready = ready);
};

const fetchRemoteState = async (deckID: string) => {

	const { data, error } = await client.decks.load(deckID);
	if (!data || error) {
		state.editor.error = unwrapErrorMessage(error);
		return;
	}

	state.content = {
		meta: {
			name: data.name,
			description: data.description || null,
			visibility: data.visibility,
		},
		cards: data.cards,
	};

	state.origin = { deckID: data.id, collectionID: data.collection_id };
	state.editor.changes = { meta: false, cards: false };
};

const watchContentEdits = () => {

	watch(() => state.content.meta, () => {

		if (!state.editor.ready) {
			return;
		}

		state.editor.changes.meta = true;
		queueStateSnapshot();

	}, { deep: true });

	watch(() => state.content.cards, () => {

		if (!state.editor.ready) {
			return;
		}

		state.editor.changes.cards = true;
		queueStateSnapshot();

	}, { deep: true });
};

const updateAppTitle = () => {
	state.editor.prevAppTitle = document.title;
	watch(() => state.content.meta.name, (name) => document.title = `${name || 'Unnamed'} | Deck editor`, { immediate: true });
};

const restoreAppTitle = () => {
	document.title = state.editor.prevAppTitle || '';
};

onMounted(async () => {

	updateAppTitle();
	watchContentEdits();

	const { deck_id } = route.params;
	if (typeof deck_id === 'string') {

		await restoreStateSnapshot(deck_id)
			.catch(error => void console.error('restoreStateSnapshot', error)) || null;

		if (!state.editor.snapshots.loadedVersion) {
			await fetchRemoteState(deck_id);
		}

		state.editor.ready = true;
		return;
	}

	const { collection_id } = Object.fromEntries(new URLSearchParams(window.location.search).entries());
	if (collection_id) {

		await restoreStateSnapshot().catch(error => console.error('restoreStateSnapshot', error));

		state.origin = { deckID: null, collectionID: collection_id };

		state.editor.ready = true;
		return;
	}

	state.editor.error = 'Invalid editor URL';
});

onUnmounted(() => {
	restoreAppTitle();
});

const patchDeckMeta = (patch: { name: string | null; description: string | null; }) => {
	state.content.meta.name = patch.name || defaultDeckMeta().name;
	state.content.meta.description = patch.description;
};

const handleVersionRollback = async () => {
	await clearStateSnapshot();
	window.location.reload();
};

const handleVersionPublish = async (meta: CardDeckMetadata) => {

	state.editor.ready = false;

	state.origin = { deckID: meta.id, collectionID: meta.collection_id };
	state.content.meta = { name: meta.name, description: meta.description || null, visibility: meta.visibility };

	await clearStateSnapshot();

	state.editor.changes = { meta: false, cards: false };

	nextTick(() => state.editor.ready = true);
};

const backHref = computed(() => state.origin?.collectionID ? `/collection/${state.origin.collectionID}` : '/collections');

const discardChanges = () => {

	if (contentEdited.value && !confirm('Really discard your changes?')) {
		return;
	}

	exitEditor();
};

const exitEditor = () => {
	clearStateSnapshot();
	router.push(backHref.value);
};

</script>

<template>

	<div class="deck-editor">

		<EditorScreenOverlay v-if="!state.editor.ready || state.editor.error">

			<OverlayErrorMessage v-if="state.editor.error" :backHref="backHref">
				
				Unable to load editor

				<template v-slot:details>
					{{ state.editor.error }}
				</template>

				<template v-slot:after>
					<GenericButton variant="thin" @click="exitEditor">
						Go back
					</GenericButton>
				</template>

			</OverlayErrorMessage>

			<template v-else>
				<LoadingMessage>
					Loading editor ...
				</LoadingMessage>
			</template>

		</EditorScreenOverlay>

		<DeckEditorStatusBar
			:meta="state.content.meta"
			:edited="contentEdited"
			:valid="editorReady"
			:autosaved="changesSaved"
			@versions="state.editor.modals.versions = true"
			@publish="state.editor.modals.publish = true"
			@import="state.editor.modals.importer = true"
			@export="state.editor.modals.exporter = true"
			@disacard="discardChanges" />

		<EditorScreenOverlay v-if="state.editor.modals.exporter">
			<EditorContentExporter
				:meta="state.origin"
				:content="state.content"
				@done="state.editor.modals.exporter = false" />
		</EditorScreenOverlay>

		<EditorScreenOverlay v-if="state.editor.modals.importer">
			<EditorContentImporter 
				@addCards="cards => state.content.cards.push(...cards)"
				@replaceCards="cards => state.content.cards = cards"
				@updateMeta="patchDeckMeta"
				@done="state.editor.modals.importer = false" />
		</EditorScreenOverlay>

		<EditorScreenOverlay v-if="state.editor.modals.versions && state.origin.deckID">
			<EditorVersionControl
				:deckID="state.origin.deckID"
				@rollback="handleVersionRollback"
				@done="state.editor.modals.versions = false" />
		</EditorScreenOverlay>

		<EditorScreenOverlay v-if="state.editor.modals.publish">
			<EditorVersionPublishModal
				:origin="state.origin"
				:content="state.content"
				:changes="state.editor.changes"
				@publish="handleVersionPublish"
				@done="state.editor.modals.publish = false" />
		</EditorScreenOverlay>

		<div class="editor-canvas">

			<div class="canvas-grid">

				<DeckCardList
					:list="cardSelectorList"
					:pointer="state.editor.view.cardIdx"
					@select="selectCard"
					@add="createCard()"
					@duplicate="duplicateCard"
					@remove="removeCard" />

				<CardFaceEditor v-model="activeCard.front" :isFront="true" />
				<hr />
				<CardFaceEditor v-model="activeCard.back" />

			</div>
		</div>

	</div>

</template>

<style lang="scss" scoped>
	.deck-editor {
		display: flex;
		flex-direction: column;
		gap: 2rem;
		width: 100%;
		height: 100%;
		position: relative;

		.editor-canvas {
			display: flex;
			justify-content: center;
			width: 100%;
			height: 100%;
			min-height: 0;

			.canvas-grid {
				display: grid;
				grid-template-columns: auto 1fr 1px 1fr;
				gap: 1rem;
				width: 100%;
				max-width: 70rem;
				height: 100%;
			}

			hr {
				display: block;
				background-color: var(--app-theme-powder-trail);
				width: 1px;
				height: 100%;
				outline: none;
				border: none;
			}
		}
	}
</style>
