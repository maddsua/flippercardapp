<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router';
import { computed, onMounted, reactive, toRaw, watch } from 'vue';
import DeckEditorStatusBar from './DeckEditorStatusBar.vue';
import CardFaceComponent from '../Cards/CardFace.vue';
import type { CardContentFace, CardNode, CardImageNode, CardContentNode } from '../../content';
import DeckCardFacePreviewSlot from './DeckCardFacePreviewSlot.vue';
import EditorCanvasColumn from './EditorCanvasColumn.vue';
import CardFaceContentEditor from './CardContentEditor/CardFaceContentEditor.vue';
import CardFaceThemeEditor from './CardContentEditor/CardFaceThemeEditor.vue';
import DeckCardList from './DeckCardList.vue';
import { useStorage } from '../../storage';
import { unwrapErrorMessage, useClient } from '../../api';
import FullscreenMessage from '../App/FullscreenMessage.vue';
import type { Card as CardType, CardDeck, ResourceVisibility } from '../../api_models';
import { blobToJson, downloadBlob, escapeFileName, pickLocalFiles } from '../../files';
import EditorScreenOverlay from './EditorScreenOverlay.vue';
import LoadingMessage from '../App/LoadingMessage.vue';
import ErrorMessage from '../App/ErrorMessage.vue';
import GenericButton from '../App/GenericButton.vue';
import WarnList from '../App/WarnList.vue';

const route = useRoute();
const router = useRouter();
const store = useStorage();
const client = useClient();

interface ActiveFace {
	id: string;
	face: CardContentFace;
};

const state = reactive({
	
	content: {
		meta: {
			name: 'Unnamed deck',
			description: null as string | null,
			visibility: 'HIDDEN' as ResourceVisibility,
		},
		cards: [] as CardNode[],
	},

	publisher: {
		deckID: null as string | null,
		collectionID: null as string | null,
		busy: false,
		error: null as string | null,
	},

	io: {
		busy: false,
		isImport: false,
		tasks: 0,
		tasksDone: 0,
		warns: null as string[] | null,
		error: null as string | null,
	},

	editor: {
		ready: false,
		error: null as string | null,
		saved: false,
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
	!state.publisher?.error &&
	!state.io.busy &&
	!state.io.error);

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

		state.publisher.deckID = data.id;

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

const resetState = () => {
	state.content = { meta: { name: 'Unnamed deck', description: null, visibility: 'HIDDEN' }, cards: [] };
	state.editor.changes = { meta: false, cards: false };
	state.io = { busy: false, isImport: false, tasks: 0, tasksDone: 0, warns: null, error: null };
	state.editor.saved = false;
	state.editor.view = { cardIdx: 0, previewAnimationTurn: false, frontFace: true };
};

const discardChanges = () => {

	if (isEdited.value && !confirm('Really discard your changes?')) {
		return;
	}

	exitEditor();
};

const exitEditor = () => {

	clearStateSnapshot();

	if (state.publisher?.collectionID) {
		router.push(`/dashboard/content/collection/${state.publisher.collectionID}`);
		return;
	}

	router.push('/dashboard/content');
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

	await store.deckEditor.store(snapshot);

	state.editor.saved = true;
};

const loadSnapshot = async () => await store.deckEditor.load() as ResumableState | null;

const clearStateSnapshot = async () => store.deckEditor.store(null);

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

const importDeckBundle = async () => {

	if (!isReady.value) {
		return;
	}

	const files = await pickLocalFiles({ accept: ['.carddeck', '.json'] });
	if (!files) {
		return;
	}

	state.io = { busy: true, isImport: true, tasks: 1, tasksDone: 0, warns: null, error: null };

	switch (files[0].type) {
		case 'application/json':
			return importJsonDeckTemplate(files[0]);
		default:
			return importCompressedDeckBundle(files[0]);
	}
};

interface CardDeckTemplate {
	type: 'flippercardtemplate';
	content: {
		name: string;
		description: string | null;
		cards: Omit<CardNode, 'id'>[];
	};
};

const importJsonDeckTemplate = async (file: File) => {

	const templ: CardDeckTemplate | null = await file.text().then(data => JSON.parse(data)).catch(() => null);
	if (!templ) {
		state.io.error = 'Invalid template file: Invalid format';
		return;
	} else if (templ.type !== 'flippercardtemplate') {
		state.io.error = 'Invalid template file: Invalid JSON schema';
		return;
	}

	if (!state.publisher.deckID) {
		state.content.meta = {
			name: templ.content.name,
			description: templ.content.description,
			visibility: 'HIDDEN',
		};
	}

	state.content.cards = templ.content.cards.map(item => ({ ...item, id: crypto.randomUUID() }));

	state.io.busy = false;
};

interface CardDeckBundle {
	type: 'flippercarddeckbundle';
	deck_id: string | null;
	collection_id: string | null;
	name: string;
	description: string | null;
	cards: Array<Omit<CardType, 'created' | 'updated'>>;
	image_blobs: ImageBlobBundle[];
};

interface ImageBlobBundle {
	id: string;
	source_name: string;
	data_url: string;
};

const importCompressedDeckBundle = async (file: File) => {

	const ds = new DecompressionStream('gzip');
	const bundleStream = file.stream().pipeThrough(ds)
	const bundle: CardDeckBundle | null = await new Response(bundleStream).json().catch(() => null);

	if (!bundle || typeof bundle !== 'object') {
		state.io.error = 'Invalid bundle file: Invalid format';
		return;
	} else if (bundle.type !== 'flippercarddeckbundle') {
		state.io.error = 'Invalid bundle file: Invalid JSON schema';
		return;
	}

	if (isEdited.value) {
		if (!confirm('Importing a bundle will discard current changes. Continue anyway?')) {
			state.io = { busy: false, isImport: false, tasks: 0, tasksDone: 0, warns: null, error: null };
			return;
		}
	}

	state.io.tasksDone++;
	state.io.tasks += bundle.image_blobs.length;

	const imageNodes = bundle.cards
		.map(item => [item.front.content, item.back.content])
		.flat()
		.flat()
		.filter(item => item.type === 'image' && item.media_id?.length) as CardImageNode[];

	for (const image of bundle.image_blobs) {

		state.io.tasksDone++;

		if (!image.data_url) {
			state.io.warns = [ ...state.io.warns || [], `Image import: Data URL is missing for ${image.id}`];
			continue;
		}

		const blob = await fetch(image.data_url).then(res => res.blob()).catch(() => null);
		if (!blob) {
			state.io.warns = [ ...state.io.warns || [], `Image import failed: Unable to parse blob for ${image.id}`];
			continue;
		}

		const { data, error } = await client.images.upload(blob, image.source_name);
		if (!data || error) {
			state.io.warns = [ ...state.io.warns || [], `Image upload failed: ${error?.message}`];
			continue;
		}

		for (const node of imageNodes) {
			if (node.media_id !== image.id) {
				continue;
			}
			node.media_id = data.id;
		}
	}

	state.content = {
		meta: {
			name: bundle.name,
			description: bundle.description,
			visibility: 'HIDDEN',
		},
		cards: bundle.cards,
	};

	state.publisher = { deckID: bundle.deck_id, collectionID: bundle.collection_id, busy: false, error: null };
	state.io.busy = false;
};

const downloadImageBlob = async (node: CardContentNode) => {

	if (node.type !== 'image' || !node.media_id?.length) {
		return null;
	}

	const { data } = await client.images.metadata(node.media_id);
	if (!data) {
		return null;
	}

	const { blob } = await client.images.blob(node.media_id);
	if (!blob) {
		return null;
	}

	return { id: node.media_id, source_name: data.source_name, blob };
};

const exportDeckBundle = async () => {

	if (!isReady.value || !state.publisher) {
		return;
	}

	state.io = { busy: false, isImport: false, tasks: 1, tasksDone: 0, warns: null, error: null };

	const bundle: CardDeckBundle = {
		type: 'flippercarddeckbundle',
		deck_id: state.publisher.deckID,
		collection_id: state.publisher.collectionID,
		name: state.content.meta.name,
		description: state.content.meta.description,
		cards: state.content.cards,
		image_blobs: await Promise.all(state.content.cards
			.map(item => [item.front.content, item.back.content])
			.flat()
			.flat()
			.filter(item => item.type === 'image' && item.media_id?.length)
			.map(item => downloadImageBlob(item)))
			.then(items => items.filter(item => item?.blob?.size))
			.then(items => Promise.all(items.map(async item => ({
				id: item!.id,
				source_name: item!.source_name,
				data_url: await blobToJson(item!.blob),
			})))),
	};

	const baseName = escapeFileName(state.content.meta.name) || 'unnamed_deck';
	const name = `${baseName}-export-${new Date().getTime()}.carddeck`;

	const rawBlob = new Blob([JSON.stringify(bundle)]);
	const compressor = new CompressionStream('gzip');
	const compressedBlob = await new Response(rawBlob.stream().pipeThrough(compressor)).blob();

	state.io.tasksDone++;

	downloadBlob(compressedBlob, name);

	state.io.busy = false;
};

const resetIOState = async () => {

	if (state.io.error) {
		await clearStateSnapshot();
		resetState();
		return;
	}

	state.io = { busy: false, isImport: false, tasks: 0, tasksDone: 0, warns: null, error: null };
};

</script>

<template>

	<div class="deck-editor">

		<EditorScreenOverlay v-if="!state.editor.ready || state.editor.error">

			<template v-if="state.editor.error">
				<ErrorMessage>
					<template v-slot:message>
						Unable to load editor
					</template>
					<template v-slot:details>
						{{ state.editor.error }}
					</template>
				</ErrorMessage>
				<GenericButton variant="thin" theme="orange" @click="exitEditor">
					Close editor
				</GenericButton>
			</template>

			<template v-else>
				<LoadingMessage>
					Loading editor ...
				</LoadingMessage>
			</template>

		</EditorScreenOverlay>

		<EditorScreenOverlay v-else-if="state.io.busy || state.io.warns?.length || state.io.error">

			<ErrorMessage v-if="state.io.error">
				<template v-slot:message>
					<template v-if="state.io.isImport">
						Unable to import cards
					</template>
					<template v-else>
						Unable to export cards
					</template>
				</template>
				<template v-slot:details>
					{{ state.io.error }}
				</template>
			</ErrorMessage>

			<LoadingMessage v-else-if="state.io.busy">
				<template v-if="state.io.isImport">
					Importing cards ...
				</template>
				<template v-else>
					Exporting cards ...
				</template>
			</LoadingMessage>		

			<WarnList v-if="state.io.warns?.length" :messages="state.io.warns" />

			<GenericButton v-if="state.io.error || state.io.warns?.length" variant="thin" theme="orange" @click="resetIOState">
				Dismiss
			</GenericButton>

		</EditorScreenOverlay>

		<EditorScreenOverlay v-if="state.publisher?.busy || state.publisher?.error">

			<template v-if="state.publisher.error">
				<ErrorMessage>
					<template v-slot:message>
						Unable to publish changes
					</template>
					<template v-slot:details>
						{{ state.publisher.error }}
					</template>
				</ErrorMessage>

				<GenericButton variant="thin" theme="orange" @click="resetPublisher">
					Dismiss
				</GenericButton>

			</template>

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
			@publish="publishChanges"
			@import="importDeckBundle"
			@export="exportDeckBundle"
			@disacard="discardChanges" />

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
