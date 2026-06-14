<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, reactive, toRaw, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { unwrapErrorMessage, useClient } from '@/api';
import type { CardDeckMetadata, CardDeckVersion, ResourceVisibility } from '@/api_models';
import type { CardNode } from '@/content';
import { useStorage, type DeckEditorHistoryMetaEntry } from '@/storage/storage';
import GenericButton from '@/components/App/Inputs/GenericButton.vue';
import LoadingMessage from '@/components/App/Messages/LoadingMessage.vue';
import EditorContentExporter from './EditorModals/EditorContentExporter.vue';
import EditorContentImporter from './EditorModals/EditorContentImporter.vue';
import DeckCardList from './EditorCardNavigationList.vue';
import EditorScreenOverlay from './EditorScreenOverlay.vue';
import EditorVersionControlModal from './EditorModals/EditorVersionControlModal.vue';
import OverlayErrorMessage from '@/components/App/Messages/OverlayErrorMessage.vue';
import CardFaceEditor from './FaceEditor/CardFaceEditor.vue';
import EditorVersionPublishModal from './EditorModals/EditorVersionPublishModal.vue';
import DeckEditorHeader from './Toolbar/DeckEditorHeader.vue';
import DeckEditorSummary from './Toolbar/DeckEditorSummary.vue';
import DeckEditorAutosaveIndicator from './Toolbar/DeckEditorAutosaveIndicator.vue';
import DeckEditorMenu from './Toolbar/DeckEditorMenu.vue';
import DeckEditorMenuEntry from './Toolbar/DeckEditorMenuEntry.vue';
import DeckEditorToolbarQuickAction from './Toolbar/DeckEditorToolbarQuickAction.vue';
import EditorDeckDetailsModal from './EditorModals/EditorDeckDetailsModal.vue';

const route = useRoute();
const router = useRouter();
const store = useStorage();
const client = useClient();

const defaultDeckSummary = () => ({
	name: 'Unnamed deck',
	description: null as string | null,
	visibility: 'HIDDEN' as ResourceVisibility,
});

const maxEditHistorySize = 20;

const state = reactive({

	content: {
		summary: defaultDeckSummary(),
		cards: [] as CardNode[],
	},

	origin: {
		deckID: null as string | null,
		collectionID: null as string | null,
		created: null as string | null,
		updated: null as string | null,
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
			details: false,
		},
		history: {
			entries: [] as ResumableState[],
			point: null as DeckEditorHistoryMetaEntry | null,
		},
		changes: {
			summary: false,
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
const contentEdited = computed(() => !!state.content.cards.length && (state.editor.changes.cards || state.editor.changes.summary));
const changesSaved = computed(() => !state.editor.snapshots.timer && !!(state.editor.snapshots.loadedVersion || state.editor.snapshots.writtenVersion));
const deckPublished = computed(() => !!state.origin.deckID);
const localDeckID = computed(() => state.origin.deckID || 'new');

const historyEditIdx = computed((): number => {

	if (!state.editor.history.point) {
		return 0;
	}

	const idx = state.editor.history.entries.findIndex(item => item.timestamp === state.editor.history.point!.timestamp);
	return idx > 0 ? idx : 0;
})

const historyCanUndo = computed((): boolean =>
	state.editor.history.entries.length > 0 &&
	historyEditIdx.value < state.editor.history.entries.length - 1);

const historyCanRedo = computed((): boolean =>
	state.editor.history.entries.length > 0 &&
	!!state.editor.history.point &&
	historyEditIdx.value > 0);

const canvasActiveCard = computed(() => {

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
	origin: typeof state.origin,
	content: typeof state['content'];
};

const cloneEditorState = (): ResumableState => {
	return {
		deck_id: localDeckID.value,
		timestamp: new Date(),
		content: structuredClone(toRaw(state.content)),
		editor: {
			view: {
				cardIdx: state.editor.view.cardIdx,
			},
		},
		origin: structuredClone(toRaw(state.origin)),
	};
};

const addHistoryVersion = (version: ResumableState) => {

	const limit = Math.min(state.editor.history.entries.length - historyEditIdx.value, maxEditHistorySize);

	state.editor.history.entries = state.editor.history.entries.slice(historyEditIdx.value, historyEditIdx.value + limit);
	state.editor.history.entries.unshift(version);

	if (state.editor.history.point) {
		const after = new Date(state.editor.history.point.timestamp.getTime() + 1);
		store.decks.editor.history.clear(localDeckID.value, after)
			.catch(error => console.debug('Trim persistent editor history:', error));
	}

	state.editor.history.point = null;
};

const editorHistoryBack = () => {

	const version = state.editor.history.entries[historyEditIdx.value + 1];
	if (!version) {
		console.warn('Cant execute editor undo');
		return;
	}

	applyEditorHistoryVersion(version);
};

const editorHistoryForward = () => {

	if (!state.editor.history.point) {
		editorHistorySkipToHead();
		return;
	}

	if (historyEditIdx.value <= 1) {
		editorHistorySkipToHead();
		return;
	}

	applyEditorHistoryVersion(state.editor.history.entries[historyEditIdx.value - 1]);
};

const editorHistorySkipToHead = () => {

	const head = state.editor.history.entries[0];
	if (head) {
		applyEditorHistoryVersion(head);
	}

	state.editor.history.point = null;
};

const applyEditorHistoryVersion = (version: ResumableState) => {

	const { ready } = state.editor;

	state.editor.ready = false;

	state.editor.history.point = { deck_id: version.deck_id, timestamp: version.timestamp };
	state.content = structuredClone(toRaw(version.content));

	nextTick(() => state.editor.ready = ready);
};

const autosaveStateSnapshot = () => {

	if (!state.editor.ready || state.editor.snapshots.lock) {
		return;
	}

	if (state.editor.snapshots.timer) {
		clearTimeout(state.editor.snapshots.timer);
	}

	const interval = 3_000;

	state.editor.snapshots.timer = setTimeout(async () => {
		const snapshot = cloneEditorState();
		addHistoryVersion(snapshot);
		await writeStateSnapshot(snapshot);
		state.editor.snapshots.timer = null;
	}, interval);
};

const writeStateSnapshot = async (snapshot: ResumableState) => {

	if (state.editor.snapshots.lock) {
		return;
	}

	state.editor.snapshots.lock = true;

	const { error } = await store.decks.editor.history.add(snapshot)
		.then(() => ({ error: null }))
		.catch(error => ({ error }));

	if (error) {
		console.error('editor.snapshots.push', error);
		state.editor.snapshots.lock = false;
		return;
	}

	state.editor.snapshots.writtenVersion = { deck_id: snapshot.deck_id, timestamp: snapshot.timestamp };
	state.editor.snapshots.lock = false;
};

const restoreStateSnapshot = async () => {

	const entries = await store.decks.editor.history.versions<ResumableState>(localDeckID.value, maxEditHistorySize).catch(() => null);
	if (!entries?.length) {
		state.editor.snapshots.loadedVersion = null;
		return;
	}

	const latest = structuredClone(entries[0]);

	state.origin = {
		deckID: latest.origin.deckID,
		collectionID: latest.origin.collectionID,
		created: latest.origin.created,
		updated: new Date().toISOString(),
	};

	state.content = latest.content;
	state.editor.view.cardIdx = latest.editor.view.cardIdx;

	state.editor.history.entries = entries;
	state.editor.changes = { summary: true, cards: true };
	state.editor.snapshots.loadedVersion = { deck_id: latest.deck_id, timestamp: latest.timestamp };
};

const clearStateSnapshot = async () => {

	const { ready } = state.editor;

	state.editor.ready = false;

	await store.decks.editor.history.clear(localDeckID.value)
		.catch(error => console.error('clearStateSnapshot', error));

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
		summary: {
			name: data.name,
			description: data.description || null,
			visibility: data.visibility,
		},
		cards: data.cards,
	};

	state.origin = {
		deckID: data.id,
		collectionID: data.collection_id,
		created: data.created,
		updated: data.updated,
	};

	state.editor.changes = { summary: false, cards: false };
	state.editor.view.cardIdx = 0;
};

const watchContentEdits = () => {

	watch(() => state.content.summary, () => {

		if (!state.editor.ready) {
			return;
		}

		state.editor.changes.summary = true;
		autosaveStateSnapshot();

	}, { deep: true });

	watch(() => state.content.cards, () => {

		if (!state.editor.ready) {
			return;
		}

		state.editor.changes.cards = true;
		autosaveStateSnapshot();

	}, { deep: true });
};

const updateAppTitle = () => {
	state.editor.prevAppTitle = document.title;
	watch(() => state.content.summary.name, (name) => document.title = `${name || 'Unnamed'} | Deck editor`, { immediate: true });
};

const restoreAppTitle = () => {
	document.title = state.editor.prevAppTitle || '';
};

onMounted(async () => {

	updateAppTitle();
	watchContentEdits();

	const { deck_id } = route.params;
	if (typeof deck_id === 'string') {

		state.origin.deckID = deck_id;

		await restoreStateSnapshot().catch(error => console.error('restoreStateSnapshot', error));

		if (!state.editor.snapshots.loadedVersion) {
			await fetchRemoteState(deck_id);
		}

		state.editor.ready = !state.editor.error;
		return;
	}

	const { collection_id } = Object.fromEntries(new URLSearchParams(window.location.search).entries());
	if (collection_id) {

		await restoreStateSnapshot().catch(error => console.error('restoreStateSnapshot', error));

		state.origin.deckID = null;
		state.origin.collectionID = collection_id;

		state.editor.ready = true;
		return;
	}

	state.editor.error = 'Invalid editor URL';
});

onUnmounted(() => {
	restoreAppTitle();
});

const patchDeckSummary = (patch: { name: string | null; description: string | null; }) => {
	state.content.summary.name = patch.name || defaultDeckSummary().name;
	state.content.summary.description = patch.description;
};

const handleVersionRollback = async () => {
	await clearStateSnapshot();
	window.location.reload();
};

const applyPulledVersion = (version: CardDeckVersion) => {

	state.content.cards = version.content.cards;
	state.origin.updated = version.created;

	state.editor.changes = { summary: false, cards: false };
	state.editor.view.cardIdx = 0;
};

const applyPublishedMeta = async (meta: CardDeckMetadata) => {

	state.editor.ready = false;

	state.origin = {
		deckID: meta.id,
		collectionID: meta.collection_id,
		created: meta.created,
		updated: meta.updated,
	};

	state.content.summary = {
		name: meta.name,
		description: meta.description || null,
		visibility: meta.visibility,
	};

	await clearStateSnapshot();

	state.editor.changes = { summary: false, cards: false };

	nextTick(() => state.editor.ready = true);
};

const dropLocalChanges = async () => {

	if (!confirm('Drop all local changes?')) {
		return;
	}

	state.editor.ready = false;

	state.editor.changes = { summary: false, cards: false };
	state.editor.history = { entries: [], point: null };

	await clearStateSnapshot();

	if (state.origin.deckID) {
		await fetchRemoteState(state.origin.deckID);
	} else {
		state.content.summary = defaultDeckSummary();
		state.content.cards = [];
	}

	if (state.editor.error) {
		return;
	}

	state.editor.ready = true;
};

const deleteDeckAndExit = async () => {

	if (!state.origin.deckID) {
		return;
	}

	if (!confirm('Delete deck?')) {
		return;
	}

	const { error } = await client.decks.remove(state.origin.deckID);
	if (error) {
		console.error('Unable to delete collection deck:', error.message);
		return;
	}

	await clearAndExitEditor();
};

const openPlayView = () => {
	if (!state.origin.deckID) {
		return;
	}
	window.open(`/play/deck/${state.origin.deckID}`, '_blank');
};

const backHref = computed(() => state.origin?.collectionID ? `/collection/${state.origin.collectionID}` : '/collections');

const clearAndExitEditorPrompt = () => {

	if (contentEdited.value && !confirm('Really discard your changes?')) {
		return;
	}

	clearAndExitEditor();
};

const clearAndExitEditor = async () => {
	await clearStateSnapshot();
	exitEditor();
};

const saveAndExitEditor = async () => {
	await writeStateSnapshot(cloneEditorState());
	exitEditor();
};

const exitEditor = () => router.push(backHref.value);

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
					<GenericButton variant="thin" @click="clearAndExitEditor">
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

		<DeckEditorHeader @exit="clearAndExitEditorPrompt">

			<template v-slot:autosave>
				<DeckEditorAutosaveIndicator :changed="contentEdited" :changesSaved="changesSaved" />
			</template>

			<template v-slot:ribbon>
				<DeckEditorMenu label="Deck">
					<DeckEditorMenuEntry label="Play deck" icon="play" :disabled="!editorReady || !deckPublished" @click="openPlayView" />
					<DeckEditorMenuEntry label="Details" icon="info" :disabled="!editorReady" @click="state.editor.modals.details = true" />
					<DeckEditorMenuEntry label="Versions" icon="history" :disabled="!editorReady" @click="state.editor.modals.versions = true" />
					<DeckEditorMenuEntry label="Import deck" icon="io" :disabled="!editorReady" @click="state.editor.modals.importer = true" />
					<DeckEditorMenuEntry label="Export deck" icon="io" :disabled="!editorReady" @click="state.editor.modals.exporter = true" />
					<DeckEditorMenuEntry label="Publish changes" icon="publish" :disabled="!editorReady" @click="state.editor.modals.publish = true" />
					<DeckEditorMenuEntry label="Discard local changes" icon="broom" :disabled="!editorReady || !contentEdited" @click="dropLocalChanges" />
					<DeckEditorMenuEntry label="Discard all and exit" icon="cross" :disabled="!editorReady || !contentEdited" @click="clearAndExitEditorPrompt" />
					<DeckEditorMenuEntry label="Delete deck" icon="delete" :disabled="!editorReady" @click="deleteDeckAndExit" />
					<DeckEditorMenuEntry label="Exit" icon="exit" @click="saveAndExitEditor" />
				</DeckEditorMenu>
				<DeckEditorMenu label="Changes">
					<DeckEditorMenuEntry label="Undo" icon="undo" :disabled="!editorReady || !historyCanUndo" @click="editorHistoryBack" />
					<DeckEditorMenuEntry label="Redo" icon="redo" :disabled="!editorReady || !historyCanRedo" @click="editorHistoryForward" />
				</DeckEditorMenu>
				<DeckEditorMenu label="Insert">
					<DeckEditorMenuEntry label="Insert title" :disabled="true" />
					<DeckEditorMenuEntry label="Insert text area" :disabled="true" />
					<DeckEditorMenuEntry label="Insert image" :disabled="true" />
					<DeckEditorMenuEntry label="Insert poll" :disabled="true" />
				</DeckEditorMenu>
			</template>

			<template v-slot:meta>
				<DeckEditorSummary
					:meta="state.content.summary"
					:changed="contentEdited"
					:changesSaved="changesSaved" />
			</template>

			<template v-slot:quickactions>
				<DeckEditorToolbarQuickAction :disabled="!editorReady" label="Delete deck" icon="delete" @click="deleteDeckAndExit" />
				<DeckEditorToolbarQuickAction :disabled="!editorReady || !deckPublished" label="Play deck" icon="play" @click="openPlayView" />
				<DeckEditorToolbarQuickAction :disabled="!editorReady || !contentEdited" label="Publish changes" icon="publish" @click="state.editor.modals.publish = true" />
			</template>

		</DeckEditorHeader>

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
				@updateMeta="patchDeckSummary"
				@done="state.editor.modals.importer = false" />
		</EditorScreenOverlay>

		<EditorScreenOverlay v-if="state.editor.modals.versions && state.origin.deckID">
			<EditorVersionControlModal
				:deckID="state.origin.deckID"
				@pull="applyPulledVersion"
				@rollback="handleVersionRollback"
				@done="state.editor.modals.versions = false" />
		</EditorScreenOverlay>

		<EditorScreenOverlay v-if="state.editor.modals.publish">
			<EditorVersionPublishModal
				:origin="state.origin"
				:content="state.content"
				:changes="state.editor.changes"
				@publish="applyPublishedMeta"
				@done="state.editor.modals.publish = false" />
		</EditorScreenOverlay>

		<EditorScreenOverlay v-if="state.editor.modals.details">
			<EditorDeckDetailsModal :content="state.content" :origin="state.origin" @done="state.editor.modals.details = false" />
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

				<CardFaceEditor v-model="canvasActiveCard.front" :isFront="true" />
				<hr />
				<CardFaceEditor v-model="canvasActiveCard.back" />

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
