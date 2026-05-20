<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router';
import { computed, onMounted, reactive, toRaw, watch } from 'vue';
import DeckEditorStatusBar from './DeckEditorStatusBar.vue';
import CardFaceComponent from '../Cards/CardFace.vue';
import type { CardContentFace, CardContentNode } from '../../content';
import DeckCardFacePreviewSlot from './DeckCardFacePreviewSlot.vue';
import EditorCanvasColumn from './EditorCanvasColumn.vue';
import CardFaceContentEditor from './CardContentEditor/CardFaceContentEditor.vue';
import CardFaceThemeEditor from './CardContentEditor/CardFaceThemeEditor.vue';
import DeckCardList from './DeckCardList.vue';
import { useStorage } from '../../storage';
import { useClient } from '../../api';
import FullscreenMessage from '../App/FullscreenMessage.vue';
import EditorErrorScreen from './EditorErrorScreen.vue';
import EditorLoadingScreen from './EditorLoadingScreen.vue';
import { downloadFile, escapeName, pickUploadFiles, type CardDeckBundle } from '../../content_io';
import type { CardDeck } from '../../api_models';

const route = useRoute();
const router = useRouter();
const store = useStorage();
const client = useClient();

interface ActiveFace {
	id: string;
	face: CardContentFace;
};

interface ErrorState {
	message: string;
	details?: string;
};

const state = reactive({
	
	content: {
		details: {
			name: 'Unnamed deck',
			description: null as string | null,
		},
		cards: [] as CardContentNode[],
	},

	view: {
		cardIdx: 0,
		frontFace: true,
		previewAnimationTurn: false,
	},

	meta: {
		id: null as string | null,
		collectionID: null as string | null,
	},

	editor: {
		detailsChanged: false,
		cardsChanged: false,
		snapshotSaved: false,
	},

	loading: false,
	locked: false,
	error: null as ErrorState | null,
});

interface ResumableState extends Pick<typeof state, 'content' | 'view' | 'meta'> {};

const isEdited = computed(() => state.editor.cardsChanged || state.editor.detailsChanged);
const isValid = computed(() => !state.error && !state.loading && state.content.cards.length > 0);

const activeCardFace = computed((): ActiveFace | null => {
	const card = state.content.cards[state.view.cardIdx];
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

	return state.view.frontFace ?
		makeActive(card.front, `${card.id}-front`)
		: makeActive(card.back, `${card.id}-back`);
});

const publishNewDeck = async () => {

	const { data, error } = await client.decks.create({
		details: state.content.details,
		content: { cards: state.content.cards },
		collection_id: state.meta.collectionID,
	});

	if (!data || error) {
		state.error = {
			message: 'Failed to publish deck',
			details: error?.message
		};
		return;
	}

	state.meta = { id: data.id, collectionID: null };
	state.editor = { detailsChanged: false, cardsChanged: false, snapshotSaved: false };

	clearStateSnapshot();
};

const patchDeckExisting = async (id: string) => {

	const { data, error } = await client.decks.update(id, {
		details: state.editor.detailsChanged ? state.content.details : null,
		content: state.editor.cardsChanged ? { cards: state.content.cards } : null,
	});

	if (!data || error) {
		state.error = {
			message: 'Failed to update deck',
			details: error?.message
		};
		return;
	}

	state.editor.detailsChanged = false;
	state.editor.cardsChanged = false;

	clearStateSnapshot();
};

const publishChanges = async () => {

	if (state.locked) {
		return;
	}

	state.loading = true;
	state.locked = true;

	if (state.meta.id) {
		await patchDeckExisting(state.meta.id);
	} else {
		await publishNewDeck();
	}

	state.loading = false;
	state.locked = false;
};

const discardChanges = () => {

	if (isEdited.value && !confirm('Really discard your changes?')) {
		return;
	}

	exitEditor();
};

const exitEditor = () => {

	clearStateSnapshot();

	if (state.meta.collectionID) {
		router.push(`/app/dashboard/content/collection/${state.meta.collectionID}`);
		return;
	}

	router.push('/app/dashboard/content');
};

const flipCardFace = () => {

	state.view.previewAnimationTurn = true;

	setTimeout(() => {
		state.view.frontFace = !state.view.frontFace;
		setTimeout(() => {
			state.view.previewAnimationTurn = false;
		}, 10);
	}, 150);
};

const addCard = (node: CardContentNode) => {
	state.content.cards.push(node);
	state.view.cardIdx = state.content.cards.length - 1;
	state.view.frontFace = true;
};

const createCard = () => {
	const newFace = () => ({ content: [] });
	addCard({ id: crypto.randomUUID(), front: newFace(), back: newFace() });
};

const selectCard = (idx: number) => {
	state.view.cardIdx = idx;
	state.view.frontFace = true;
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

	if (state.view.cardIdx >= state.content.cards.length) {
		state.view.cardIdx = state.content.cards.length - 1;
	}
};

const storeStateSnapshot = async () => {

	if (state.locked || state.editor.snapshotSaved) {
		return;
	}

	state.locked = true;

	const snapshot: ResumableState = { content: state.content, view: state.view, meta: state.meta };
	await store.deckEditor.store(snapshot);

	state.editor.snapshotSaved = true;
	state.locked = false;
};

const loadSnapshot = async () => await store.deckEditor.load() as ResumableState | null;

const clearStateSnapshot = async () => store.deckEditor.store(null);

const initAutosave = () => {

	watch(() => state.content.details, () => {
		state.editor.detailsChanged = true;
		state.editor.snapshotSaved = false;
	}, { deep: true });

	watch(() => state.content.cards, () => {
		state.editor.cardsChanged = true;
		state.editor.snapshotSaved = false;
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
		details: {
			name: deck.name,
			description: deck.description || null,
		},
		cards: deck.cards,
	};

	state.meta.id = deck.id;
	state.meta.collectionID = deck.collection_id;

	state.editor.cardsChanged = true;
	state.editor.detailsChanged = true;
};

const applySnapshotState = (snapshot: ResumableState) => {
	state.meta = snapshot.meta;
	state.content = snapshot.content;
};

const resolveDeckState = async (data: CardDeck) => {

	const snapshot = await loadSnapshot();
	if (!snapshot) {
		applyRemoteDeckState(data);
		return null;
	}

	const { id: storedID } = snapshot.meta;

	if (storedID === data.id) {
		applySnapshotState(snapshot);
		return null
	}

	if (!confirm('Editor contains unsaved changed of another deck. Overwrite changes or load them?')) {
		applySnapshotState(snapshot);
		state.meta = { id: data.id, collectionID: null };
		return;
	}

	await clearStateSnapshot();
	applyRemoteDeckState(data);

	return null;
};

onMounted(async () => {

	const { deck_id } = route.params;
	if (typeof deck_id === 'string') {

		const { data, error } = await client.decks.load(deck_id);
		if (!data || error) {
			state.error = {
				message: 'Unable to load deck',
				details: error?.message,
			};
			return;
		}

		await resolveDeckState(data);

		initAutosave();
		
		return;
	}

	const { collection_id } = Object.fromEntries(new URLSearchParams(window.location.search).entries());
	if (collection_id) {

		const snapshot = await loadSnapshot();
		if (snapshot) {
			applySnapshotState(snapshot);
		}

		state.meta.collectionID = collection_id;

		initAutosave();

		return;
	}

	state.error = {
		message: 'Invalid editor URL',
		details: 'Collection or deck id parameter not provided'
	};
});

const importDeckBundle = async () => {

	if (state.loading || state.locked) {
		return;
	}

	const files = await pickUploadFiles({ accept: ['.json'] });
	if (!files) {
		return;
	}

	const bundle: CardDeckBundle | null = await files[0].text()
		.then(data => JSON.parse(data))
		.catch(() => null);

	if (!bundle || typeof bundle !== 'object' || bundle.type !== 'maddsua:flippercarddapp:bundle:deck') {
		console.warn('Invalid bundle file: Invalid JSON object or unsupported bundle type');
		return;
	}

	if (isEdited.value) {
		if (!confirm('Importing a bundle will discard current changes. Continue anyway?')) {
			return;
		}
	}

	state.content = {
		details: {
			name: bundle.name,
			description: bundle.description,
		},
		cards: bundle.cards,
	};

	state.meta = {
		id: bundle.id,
		collectionID: bundle.collection_id,
	};
};

const exportDeckBundle = async () => {

	const bundle: CardDeckBundle = {
		type: 'maddsua:flippercarddapp:bundle:deck',
		id: state.meta.id,
		collection_id: state.meta.collectionID,
		name: state.content.details.name,
		description: state.content.details.description,
		cards: state.content.cards,
	};

	const name = escapeName(state.content.details.name) || 'unnamed_deck';

	downloadFile(JSON.stringify(bundle), `${name}-export-${new Date().getTime()}.json`);
};

const updateDetails = (val: { name: string, description?: string | null }) => {
	state.content.details.name = val.name;
	state.content.details.description = val.description || null;
};

</script>

<template>

	<div class="deck-editor">

		<EditorLoadingScreen v-if="state.loading" />
		<EditorErrorScreen v-else-if="state.error" :error="state.error" />

		<DeckEditorStatusBar
			:meta="state.content.details"
			:edited="isEdited"
			:valid="isValid"
			@updateDetails="updateDetails"
			@flip="flipCardFace"
			@publish="publishChanges"
			@import="importDeckBundle"
			@export="exportDeckBundle"
			@disacard="discardChanges" />

		<div class="editor-canvas">

			<div class="canvas-grid">

				<DeckCardList
					:size="state.content.cards.length"
					:activeIdx="state.view.cardIdx"
					@select="selectCard"
					@add="createCard()"
					@duplicate="duplicateCard"
					@remove="removeCard" />

				<EditorCanvasColumn>
					<template v-slot:title>

						Preview

						<template v-if="state.view.frontFace">
							(front)
						</template>
						<template v-else>
							(back)
						</template>

					</template>

					<template v-slot:content>

						<DeckCardFacePreviewSlot v-if="activeCardFace" :turned="state.view.previewAnimationTurn">
							<template v-if="state.view.frontFace">
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

						<template v-if="state.view.frontFace">
							(front)
						</template>
						<template v-else>
							(back)
						</template>

					</template>

					<template v-slot:content>

						<template v-if="activeCardFace">
							<CardFaceContentEditor v-model="activeCardFace.face.content" :isFront="state.view.frontFace" />
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
