<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router';
import { computed, onMounted, reactive, toRaw, watch } from 'vue';
import DeckEditorStatusBar from './DeckEditorStatusBar.vue';
import CardFaceComponent from '../Cards/CardFace.vue';
import type { CardContentFace, CardNode, CardImageElement, CardContentNode } from '../../content';
import DeckCardFacePreviewSlot from './DeckCardFacePreviewSlot.vue';
import EditorCanvasColumn from './EditorCanvasColumn.vue';
import CardFaceContentEditor from './CardContentEditor/CardFaceContentEditor.vue';
import CardFaceThemeEditor from './CardContentEditor/CardFaceThemeEditor.vue';
import DeckCardList from './DeckCardList.vue';
import { useStorage } from '../../storage';
import { useClient } from '../../api';
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

interface ErrorState {
	message: string;
	details?: string;
};

interface ImportState {
	total: number;
	progress: number;
	warns: string[];
	error: string | null;
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

	meta: {
		id: null as string | null,
		collectionID: null as string | null,
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

	snapshotSaved: false,
	exportActive: false,
	importer: null as ImportState | null,

	loading: false,
	locked: false,
	error: null as ErrorState | null,
});

interface ResumableState extends Pick<typeof state, 'content' | 'view' | 'meta'> {};

const isEdited = computed(() => state.changes.cards || state.changes.meta);
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

const cardSelectorList = computed(() => state.content.cards.map(item => item.front));

const publishNewDeck = async () => {

	const { data, error } = await client.decks.create({
		meta: state.content.meta,
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

	state.meta.id = data.id;
	state.changes = { meta: false, cards: false };
	state.snapshotSaved = false;

	await clearStateSnapshot();
};

const patchDeckExisting = async (id: string) => {

	const { data, error } = await client.decks.update(id, {
		meta: state.changes.meta ? state.content.meta : null,
		content: state.changes.cards ? { cards: state.content.cards } : null,
	});

	if (!data || error) {
		state.error = {
			message: 'Failed to update deck',
			details: error?.message
		};
		return;
	}

	state.changes = { meta: false, cards: false };
	state.snapshotSaved = false;

	await clearStateSnapshot();
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

const clearState = () => {
	state.content = { meta: { name: '', description: null, visibility: 'HIDDEN' }, cards: [] };
	state.changes = { cards: false, meta: false };
	state.exportActive = false;
	state.importer = null;
	state.snapshotSaved = false;
	state.view = { cardIdx: 0, previewAnimationTurn: false, frontFace: true };
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

const addCard = (node: CardNode) => {
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

	if (state.locked || state.snapshotSaved) {
		return;
	}

	state.locked = true;

	const snapshot: ResumableState = { content: state.content, view: state.view, meta: state.meta };
	await store.deckEditor.store(snapshot);

	state.snapshotSaved = true;
	state.locked = false;
};

const loadSnapshot = async () => await store.deckEditor.load() as ResumableState | null;

const clearStateSnapshot = async () => store.deckEditor.store(null);

const initAutosave = () => {

	watch(() => state.content.meta, () => {
		state.changes.meta = true;
		state.snapshotSaved = false;
	}, { deep: true });

	watch(() => state.content.cards, () => {
		state.changes.cards = true;
		state.snapshotSaved = false;
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

	state.meta.id = deck.id;
	state.meta.collectionID = deck.collection_id;

	state.changes.cards = false;
	state.changes.meta = false;
};

const applySnapshotState = (snapshot: ResumableState) => {
	state.meta = snapshot.meta;
	state.content = snapshot.content;
	state.changes.cards = true;
	state.changes.meta = true;
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

	const files = await pickLocalFiles({ accept: ['.carddeck', '.json'] });
	if (!files) {
		return;
	}

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

	state.importer = { total: 1, progress: 0, warns: [], error: null };

	const templ: CardDeckTemplate | null = await file.text().then(data => JSON.parse(data)).catch(() => null);
	if (!templ) {
		state.importer.error = 'Invalid template file: Invalid format';
		return;
	} else if (templ.type !== 'flippercardtemplate') {
		state.importer.error = 'Invalid template file: Invalid JSON schema';
		return;
	}

	state.content = {
		meta: {
			name: templ.content.name,
			description: templ.content.description,
			visibility: 'HIDDEN',
		},
		cards: templ.content.cards.map(item => ({ ...item, id: crypto.randomUUID() })),
	};

	state.importer = null;
};

interface CardDeckBundle {
	type: 'flippercarddeckbundle';
	id: string | null;
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

	state.importer = { total: 1, progress: 0, warns: [], error: null };

	const ds = new DecompressionStream('gzip');
	const bundleStream = file.stream().pipeThrough(ds)
	const bundle: CardDeckBundle | null = await new Response(bundleStream).json().catch(() => null);

	if (!bundle || typeof bundle !== 'object') {
		state.importer.error = 'Invalid bundle file: Invalid format';
		return;
	} else if (bundle.type !== 'flippercarddeckbundle') {
		state.importer.error = 'Invalid bundle file: Invalid JSON schema';
		return;
	}

	if (isEdited.value) {
		if (!confirm('Importing a bundle will discard current changes. Continue anyway?')) {
			state.importer = null;
			return;
		}
	}

	state.importer.progress++;
	state.importer.total += bundle.image_blobs.length;

	const imageNodes = bundle.cards
		.map(item => [item.front.content, item.back.content])
		.flat()
		.flat()
		.filter(item => item.type === 'image' && item.media_id?.length) as CardImageElement[];

	for (const image of bundle.image_blobs) {

		state.importer.progress++;

		if (!image.data_url) {
			state.importer.warns.push(`Image import: Data URL is missing for ${image.id}`);
			continue;
		}

		const blob = await fetch(image.data_url).then(res => res.blob()).catch(() => null);
		if (!blob) {
			state.importer.warns.push(`Image import failed: Unable to parse blob for ${image.id}`);
			continue;
		}

		const { data, error } = await client.images.upload(blob, image.source_name);
		if (!data || error) {
			state.importer.warns.push(`Image upload failed: ${error?.message}`);
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

	state.meta = {
		id: bundle.id,
		collectionID: bundle.collection_id,
	};

	state.importer = null;
};

const abortDeckImport = async () => {
	await clearStateSnapshot();
	clearState();
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

	state.exportActive = true;

	const bundle: CardDeckBundle = {
		type: 'flippercarddeckbundle',
		id: state.meta.id,
		collection_id: state.meta.collectionID,
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

	downloadBlob(compressedBlob, name);

	state.exportActive = false;
};

</script>

<template>

	<div class="deck-editor">

		<EditorScreenOverlay v-if="state.loading">
			<LoadingMessage />
		</EditorScreenOverlay>

		<EditorScreenOverlay v-else-if="state.error" :error="state.error">
			<ErrorMessage>
				<template v-slot:message>
					{{ state.error.message }}
				</template>
				<template v-if="state.error.details" v-slot:details>
					{{ state.error.details }}
				</template>
			</ErrorMessage>
		</EditorScreenOverlay>

		<EditorScreenOverlay v-else-if="state.exportActive">
			<LoadingMessage>
				Exporting content
			</LoadingMessage>
		</EditorScreenOverlay>

		<EditorScreenOverlay v-else-if="state.importer">

			<template v-if="state.importer.error">

				<ErrorMessage>
					<template v-slot:message>
						Import failed
					</template>
					<template v-slot:details>
						{{ state.importer.error }}
					</template>
				</ErrorMessage>

				<GenericButton variant="thin" theme="orange" @click="abortDeckImport">
					Dismiss
				</GenericButton>

			</template>

			<template v-else>
				<LoadingMessage>
					Importing part {{ state.importer.progress }}/{{ state.importer.total }}
				</LoadingMessage>
				<WarnList v-if="state.importer.warns.length" :messages="state.importer.warns" />
			</template>

		</EditorScreenOverlay>

		<DeckEditorStatusBar
			:meta="state.content.meta"
			:edited="isEdited"
			:valid="isValid"
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
					:pointer="state.view.cardIdx"
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
