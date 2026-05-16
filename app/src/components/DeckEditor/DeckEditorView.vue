<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router';
import { computed, onMounted, reactive, watch } from 'vue';
import DeckEditorStatusBar from './DeckEditorStatusBar.vue';
import CardFaceComponent from '../Cards/CardFace.vue';
import type { CardFace, CardNode } from '../../content';
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

const route = useRoute();
const router = useRouter();
const store = useStorage();
const client = useClient();

interface ActiveFace {
	id: string;
	face: CardFace;
};

interface ErrorState {
	message: string;
	details?: string;
};

const state = reactive({
	
	data: {
		meta: {
			name: 'Unnamed deck',
			description: null as string | null,
		},
		cards: [] as CardNode[],
	},

	view: {
		cardIdx: 0,
		frontFace: true,
		previewAnimationTurn: false,
	},

	metadata: {
		id: null as string | null,
		collectionID: null as string | null,
	},

	editor: {
		cardsChanged: false,
		metaChanged: false,
		snapshotSaved: false,
	},

	loading: false,
	locked: false,
	error: null as ErrorState | null,
});

const isEdited = computed(() => state.editor.cardsChanged || state.editor.metaChanged);
const isValid = computed(() => !state.error && !state.loading && state.data.cards.length > 0);

const activeCardFace = computed((): ActiveFace | null => {
	const card = state.data.cards[state.view.cardIdx];
	if (!card) {
		return null;
	}

	const makeActive = (face: CardFace, id: string) => {

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

const wrapCards = () => state.data.cards.map(item => ({ id: item.id || null, content: item }));

const publishNewDeck = async () => {

	const { data, error } = await client.decks.add({
		...state.data.meta,
		cards: wrapCards(),
		collection_id: state.metadata.collectionID,
	});

	if (!data || error) {
		state.error = {
			message: 'Failed to publish deck',
			details: error?.message
		};
		return;
	}

	state.metadata = { id: data.id, collectionID: null };
	state.editor = { metaChanged: false, cardsChanged: false, snapshotSaved: false };

	store.deckEditor.store(null);
};

const patchDeckExisting = async (id: string) => {

	if (state.editor.metaChanged) {

		const { data, error } = await client.decks.updateMeta(id, state.data.meta);
		if (!data || error) {
			state.error = {
				message: 'Failed to update metadata',
				details: error?.message
			};
			return;
		}

		state.editor.metaChanged = false;
	}

	if (state.editor.cardsChanged) {

		const { data, error } = await client.decks.updateContent(id, { cards: wrapCards() });
		if (!data || error) {
			state.error = {
				message: 'Failed to update content',
				details: error?.message
			};
			return;
		}

		state.editor.cardsChanged = false;
	}
	
	store.deckEditor.store(null);
};

const handlePublish = async () => {

	if (state.locked) {
		return;
	}

	state.loading = true;
	state.locked = true;

	if (state.metadata.id) {
		await patchDeckExisting(state.metadata.id);
	} else {
		await publishNewDeck();
	}

	state.loading = false;
	state.locked = false;
};

const handleDiscard = () => {

	if (isEdited.value && !confirm('Really discard your changes?')) {
		return;
	}

	store.deckEditor.store(null);
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

const addNewCard = () => {
	const newFace = () => ({ content: [] });
	state.data.cards.push({ id: crypto.randomUUID(), front: newFace(), back: newFace() });
	state.view.cardIdx = state.data.cards.length - 1;
	state.view.frontFace = true;
};

const selectCard = (idx: number) => {
	state.view.cardIdx = idx;
	state.view.frontFace = true;
};

const removeCard = (idx: number) => {

	if (!confirm('Remove card?')) {
		return;
	}

	state.data.cards.splice(idx, 1);

	if (state.view.cardIdx >= state.data.cards.length) {
		state.view.cardIdx = state.data.cards.length - 1;
	}
};

const initAutosave = () => {

	watch(() => state.data.meta, () => {
		state.editor.metaChanged = true;
		state.editor.snapshotSaved = false;
	}, { deep: true });

	watch(() => state.data.cards, () => {
		state.editor.cardsChanged = true;
		state.editor.snapshotSaved = false;
	}, { deep: true });

	setInterval(async () => {

		if (!isEdited.value || state.locked || state.editor.snapshotSaved) {
			return;
		}

		await store.deckEditor.store(state.data);
		state.editor.snapshotSaved = true;

	}, 1000);
};

const loadDeckState = async (id: string) => {

	const { data, error } = await client.decks.load(id);
	if (!data || error) {
		state.error = {
			message: 'Unable to load deck',
			details: error?.message,
		};
		return;
	}

	state.data = {
		meta: {
			name: data.name,
			description: data.description || null,
		},
		cards: data.cards.map(item => ({ ...item.content, id: item.id })),
	};

	state.metadata.id = data.id;

	initAutosave();
};

onMounted(async () => {

	const { deck_id } = route.params;

	if (typeof deck_id === 'string') {
		state.loading = true;
		await loadDeckState(deck_id);
		state.loading = false;
		return;
	}

	const { collection_id } = Object.fromEntries(new URLSearchParams(window.location.search).entries());
	if (!collection_id) {
		state.error = {
			message: 'Invalid editor URL',
			details: 'collection id parameter not provided'
		};
		return;
	}

	state.metadata.collectionID = collection_id;

	const storedState = await store.deckEditor.load();
	if (storedState && typeof storedState === 'object') {
		state.data = storedState;
		state.editor = { metaChanged: true, cardsChanged: true, snapshotSaved: false };
	}

	initAutosave();
});

</script>

<template>

	<div class="deck-editor">

		<EditorLoadingScreen v-if="state.loading" />
		<EditorErrorScreen v-else-if="state.error" :error="state.error" />

		<DeckEditorStatusBar
			:meta="state.data.meta"
			:edited="isEdited"
			:valid="isValid"
			@editMeta="meta => state.data.meta = meta"
			@flip="flipCardFace"
			@publish="handlePublish"
			@disacard="handleDiscard" />

		<div class="editor-canvas">

			<div class="canvas-grid">

				<DeckCardList
					:size="state.data.cards.length"
					:activeIdx="state.view.cardIdx"
					@select="selectCard"
					@add="addNewCard()"
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
