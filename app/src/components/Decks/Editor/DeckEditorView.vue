<script setup lang="ts">
import { computed, onMounted, onUnmounted, reactive, toRaw, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { unwrapErrorMessage, useClient } from '@/api';
import type { CardDeck, ResourceVisibility } from '@/api_models';
import type { CardContentFace, CardNode } from '@/content';
import { useStorage } from '@/storage/storage';
import FullscreenMessage from '@/components/App/Messages/FullscreenMessage.vue';
import GenericButton from '@/components/App/Inputs/GenericButton.vue';
import LoadingMessage from '@/components/App/Messages/LoadingMessage.vue';
import CardFaceComponent from '@/components/Cards/CardFace.vue';
import EditorContentExporter from './EditorModals/EditorContentExporter.vue';
import EditorContentImporter from './EditorModals/EditorContentImporter.vue';
import CardFaceContentEditor from './ContentNodeEditable/CardFaceContentEditor.vue';
import CardFaceThemeEditor from './ContentNodeEditable/CardFaceThemeEditor.vue';
import DeckCardFacePreviewSlot from './DeckCardFacePreviewSlot.vue';
import DeckEditorStatusBar from './DeckEditorStatusBar.vue';
import EditorCanvasColumn from './EditorCanvasColumn.vue';
import DeckCardList from './EditorCardNavigationList.vue';
import EditorScreenOverlay from './EditorScreenOverlay.vue';
import EditorContentVersionControl from './EditorModals/EditorContentVersionControl.vue';
import OverlayErrorMessage from '@/components/App/Messages/OverlayErrorMessage.vue';

const route = useRoute();
const router = useRouter();
const store = useStorage();
const client = useClient();

interface ActiveFace {
	id: string;
	face: CardContentFace;
};

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

	publisher: {
		deckID: null as string | null,
		collectionID: null as string | null,
		busy: false,
		error: null as string | null,
	},

	editor: {
		ready: false,
		error: null as string | null,
		saved: false,
		oldTitle: null as string | null,
		modals: {
			versions: false,
			importer: false,
			exporter: false,
		},
		view: {
			cardIdx: 0,
			frontFace: true,
			previewAnimationTurn: false,
		},
		changes: {
			meta: false,
			cards: false,
		},
	},
});

const isReady = computed(() =>
	state.editor.ready &&
	!state.editor.error &&
	!state.publisher?.busy &&
	!state.publisher?.error);

const isEdited = computed(() => !!state.content.cards.length && (state.editor.changes.cards || state.editor.changes.meta));

const activeCardFace = computed((): ActiveFace | null => {
	const card = state.content.cards[state.editor.view.cardIdx];
	if (!card) {
		return null;
	}

	const makeActive = (face: CardContentFace, id: string) => {

		// it is a bit of a fucking hack but it's better than putting 100500 nested proxies and/or event handlers
		if (!face.theme) {
			face.theme = { card: {}, interactives: { }};
		} else if (!face.theme.interactives) {
			face.theme.interactives = {};
		}

		return { face, id };
	};

	return state.editor.view.frontFace ?
		makeActive(card.front, `${card.id}-front`)
		: makeActive(card.back, `${card.id}-back`);
});

const cardSelectorList = computed(() => state.content.cards.map(item => item.front));

const publishChanges = async () => {

	if (!state.publisher || state.publisher.busy) {
		return;
	}

	state.publisher.busy = true;

	const { deckID, collectionID } = state.publisher;

	if (deckID) {

		const { data, error } = await client.decks.update(deckID, {
			meta: state.editor.changes.meta ? state.content.meta : null,
			content: state.editor.changes.cards ? { cards: state.content.cards } : null,
		});

		if (!data || error) {
			state.publisher.error = unwrapErrorMessage(error);
			state.publisher.busy = false;
			return;
		}

	} else if (collectionID) {

		const { data, error } = await client.decks.create({
			meta: state.content.meta,
			content: { cards: state.content.cards },
			collection_id: collectionID,
		});

		if (!data || error) {
			state.publisher.error = unwrapErrorMessage(error);
			state.publisher.busy = false;
			return;
		}

		const { id: deckID } = data;
		state.publisher.deckID = deckID;
		router.push(`/decks/editor/${deckID}`);

	} else {
		state.publisher.error = 'Invalid editor state'
		state.publisher.busy = false;
		return;
	}

	state.editor.changes = { meta: false, cards: false };
	state.editor.saved = false;

	await clearStateSnapshot();

	state.publisher.busy = false;
};

const resetPublisher = () => {
	if (!state.publisher) {
		return;
	}
	state.publisher.busy = false;
	state.publisher.error = null;
};

const discardChanges = () => {

	if (isEdited.value && !confirm('Really discard your changes?')) {
		return;
	}

	exitEditor();
};

const backHref = computed(() => state.publisher?.collectionID ? `/collection/${state.publisher.collectionID}` : '/collections');

const exitEditor = () => {
	clearStateSnapshot();
	router.push(backHref.value);
};

const flipCardFace = () => {

	state.editor.view.previewAnimationTurn = true;

	setTimeout(() => {
		state.editor.view.frontFace = !state.editor.view.frontFace;
		setTimeout(() => {
			state.editor.view.previewAnimationTurn = false;
		}, 10);
	}, 150);
};

const addCard = (node: CardNode) => {
	state.content.cards.push(node);
	state.editor.view.cardIdx = state.content.cards.length - 1;
	state.editor.view.frontFace = true;
};

const createCard = () => {
	const newFace = () => ({ content: [] });
	addCard({ id: crypto.randomUUID(), front: newFace(), back: newFace() });
};

const selectCard = (idx: number) => {
	state.editor.view.cardIdx = idx;
	state.editor.view.frontFace = true;
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

interface ResumableState extends Pick<typeof state, 'content'> {
	editor: Pick<typeof state['editor']['view'], 'cardIdx'>;
	publisher: Pick<typeof state['publisher'], 'deckID' | 'collectionID'>;
};

const storeStateSnapshot = async () => {

	if (isReady.value || state.editor.saved) {
		return;
	}

	const snapshot: ResumableState = {
		content: state.content,
		editor: {
			cardIdx: state.editor.view.cardIdx,
		},
		publisher: {
			deckID: state.publisher?.deckID || null,
			collectionID: state.publisher?.collectionID || null,
		},
	};

	store.decks.editor.snapshot.store(snapshot);

	state.editor.saved = true;
};

const loadSnapshot = async () => await store.decks.editor.snapshot.load() as ResumableState | null;

const clearStateSnapshot = async () => store.decks.editor.snapshot.clear();

const initAutosave = () => {

	watch(() => state.content.meta, () => {
		state.editor.changes.meta = true;
		state.editor.saved = false;
	}, { deep: true });

	watch(() => state.content.cards, () => {
		state.editor.changes.cards = true;
		state.editor.saved = false;
	}, { deep: true });

	setInterval(async () => {

		if (!isEdited.value) {
			return;
		}

		await storeStateSnapshot();

	}, 1000);
};

const applyRemoteDeckState = (deck: CardDeck) => {

	state.content = {
		meta: {
			name: deck.name,
			description: deck.description || null,
			visibility: deck.visibility,
		},
		cards: deck.cards,
	};

	state.publisher = { deckID: deck.id, collectionID: deck.collection_id, busy: false, error: null };
	state.editor.changes = { meta: false, cards: false };
};

const applySnapshotState = (snapshot: ResumableState) => {
	state.publisher = { deckID: snapshot.publisher.deckID, collectionID: snapshot.publisher.collectionID, busy: false, error: null };
	state.content = snapshot.content;
	state.editor.changes = { meta: false, cards: false };
};

const resolveDeckState = async (deck: CardDeck) => {

	const snapshot = await loadSnapshot();
	if (!snapshot) {
		applyRemoteDeckState(deck);
		return null;
	}

	const { deckID: storedID } = snapshot.publisher;

	if (storedID === deck.id) {
		applySnapshotState(snapshot);
		return null
	}

	if (!confirm('Editor contains unsaved changed of another deck. Overwrite changes or load them?')) {
		applySnapshotState(snapshot);
		state.publisher = { deckID: deck.id, collectionID: deck.collection_id, busy: false, error: null };
		return;
	}

	await clearStateSnapshot();
	applyRemoteDeckState(deck);

	return null;
};

onMounted(async () => {

	state.editor.oldTitle = document.title;
	watch(() => state.content.meta.name, (name) => document.title = `${name || 'Unnamed'} | Deck editor`, { immediate: true });

	const { deck_id } = route.params;
	if (typeof deck_id === 'string') {

		const { data, error } = await client.decks.load(deck_id);
		if (!data || error) {
			state.editor.error = unwrapErrorMessage(error);
			return;
		}

		await resolveDeckState(data);

		initAutosave();
		
		state.editor.ready = true;
		return;
	}

	const { collection_id } = Object.fromEntries(new URLSearchParams(window.location.search).entries());
	if (collection_id) {

		const snapshot = await loadSnapshot();
		if (snapshot) {
			applySnapshotState(snapshot);
		}

		state.publisher = { deckID: null, collectionID: collection_id, busy: false, error: null };

		initAutosave();

		state.editor.ready = true;
		return;
	}

	state.editor.error = 'Invalid editor URL';
});

onUnmounted(() => document.title = state.editor.oldTitle || '');

const patchDeckMeta = (patch: { name: string | null; description: string | null; }) => {
	state.content.meta.name = patch.name || defaultDeckMeta().name;
	state.content.meta.description = patch.description;
};

const handleVersionRollback = async () => {
	await clearStateSnapshot();
	window.location.reload();
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

		<EditorScreenOverlay v-if="state.publisher?.busy || state.publisher?.error">

			<OverlayErrorMessage v-if="state.publisher.error" :backHref="backHref">
				
				Unable to publish changes

				<template v-slot:details>
					{{ state.publisher.error }}
				</template>

				<template v-slot:after>
					<GenericButton variant="thin" theme="orange" @click="resetPublisher">
						Dismiss
					</GenericButton>
				</template>

			</OverlayErrorMessage>

			<template v-else>
				<LoadingMessage>
					Publishing changes ...
				</LoadingMessage>
			</template>

		</EditorScreenOverlay>

		<DeckEditorStatusBar
			:meta="state.content.meta"
			:edited="isEdited"
			:valid="isReady"
			@updateMeta="meta => state.content.meta = meta"
			@flip="flipCardFace"
			@versions="state.editor.modals.versions = true"
			@publish="publishChanges"
			@import="state.editor.modals.importer = true"
			@export="state.editor.modals.exporter = true"
			@disacard="discardChanges" />

		<EditorScreenOverlay v-if="state.editor.modals.exporter">
			<EditorContentExporter
				:meta="state.publisher"
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

		<EditorScreenOverlay v-if="state.editor.modals.versions && state.publisher.deckID">
			<EditorContentVersionControl
				:deckID="state.publisher.deckID"
				@rollback="handleVersionRollback"
				@done="state.editor.modals.versions = false" />
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

				<EditorCanvasColumn>
					<template v-slot:title>

						Preview

						<template v-if="state.editor.view.frontFace">
							(front)
						</template>
						<template v-else>
							(back)
						</template>

					</template>

					<template v-slot:content>

						<DeckCardFacePreviewSlot v-if="activeCardFace" :turned="state.editor.view.previewAnimationTurn">
							<template v-if="state.editor.view.frontFace">
								<CardFaceComponent :key="activeCardFace.id" :entry="activeCardFace.face" decoration="question-mark" />
							</template>
							<template v-else>
								<CardFaceComponent :key="activeCardFace.id" :entry="activeCardFace.face" />
							</template>
						</DeckCardFacePreviewSlot>

						<template v-else>
							<FullscreenMessage>
								Preview not available
							</FullscreenMessage>
						</template>

					</template>

				</EditorCanvasColumn>

				<hr />

				<EditorCanvasColumn>

					<template v-slot:title>

						Editor

						<template v-if="state.editor.view.frontFace">
							(front)
						</template>
						<template v-else>
							(back)
						</template>

					</template>

					<template v-slot:content>

						<template v-if="activeCardFace">
							<CardFaceContentEditor v-model="activeCardFace.face.content" :isFront="state.editor.view.frontFace" />
							<CardFaceThemeEditor v-model="activeCardFace.face.theme" />
						</template>

						<template v-else>
							<FullscreenMessage>
								No page selected
							</FullscreenMessage>
						</template>

					</template>

				</EditorCanvasColumn>

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
				gap: 2rem;
				width: 100%;
				max-width: 70rem;
				height: 100%;
			}

			hr {
				display: block;
				background-color: var(--app-theme-powder-trail);
				width: 100%;
				height: 100%;
				outline: none;
				border: none;
			}
		}
	}
</style>
