<script setup lang="ts">
import { parse } from 'papaparse';
import { computed, reactive } from 'vue';
import { useClient } from '@/api';
import {
	parseQuizOptions,
	parseTextBoxContent,
	type CardContentCSVRow,
	type CardImageNode,
	type CardNode,
	type ContentBundle,
} from '@/content';
import { pickLocalFiles } from '@/files';
import GenericButton from '@/components/App/Inputs/GenericButton.vue';
import GenericToggle from '@/components/App/Inputs/GenericToggle.vue';
import InlineErrorMessage from '@/components/App/Messages/InlineErrorMessage.vue';
import InlineProgressIndicator from '@/components/App/Messages/InlineProgressIndicator.vue';
import EditorCardSelectorGrid from '../EditorCardSelectorGrid.vue';
import EditorModal from '../EditorModal.vue';
import OverlayErrorMessage from '@/components/App/Messages/OverlayErrorMessage.vue';

const client = useClient();

interface Meta {
	name: string | null;
	description: string | null;
};

const emit = defineEmits<{
	(e: 'done'): void;
	(e: 'addCards', cards: CardNode[]): void;
	(e: 'replaceCards', cards: CardNode[]): void;
	(e: 'updateMeta', meta: Meta): void;
}>();

interface ImageData {
	media_id: string;
	media_url: string | null;
	source_name: string;
	blob: Blob;
};

const state = reactive({
	content: {
		meta: {
			name: null as string | null,
			description: null as string | null,
		},
		cards: [] as CardNode[],
		images: [] as ImageData[],
	},
	selectedCards: new Set<string>(),
	options: {
		setMetadata: {
			available: false,
			value: false,
		},
		overwriteContent: false,
	},
	source: {
		name: null as string | null,
		ready: false,
		error: null as string | null,
	},
	dragging: false,
	progress: 0,
	total: 0,
	busy: false,
	cancelled: false,
	error: null as string | null,
});

const sourceLoading = computed(() => !!state.source.name && !state.source.ready);
const stateValid = computed(() => state.source.ready && !state.error && !state.cancelled);
const selectedCards = computed(() => state.content.cards.filter(item => state.selectedCards.has(item.id)));

const clearCardSelection = () => state.selectedCards.clear();

const selectAllCards = () => {

	clearCardSelection();

	for (const card of state.content.cards) {
		state.selectedCards.add(card.id);
	}
};

const exitTool = () => {

	if (state.busy) {
		state.error = 'Import cancelled by user';
		console.debug(state.error);
		state.cancelled = true;
	}

	revokeBlobUrls();

	emit('done');
};

const revokeBlobUrls = () => {
	for (const img of state.content.images) {
		if (!img.media_url) {
			continue;
		}
		URL.revokeObjectURL(img.media_url);
		img.media_url = null;
	}
};

const filterImageNodes = (cards: CardNode[]): CardImageNode[] =>
	cards.map(item => [item.front, item.back])
		.flat()
		.map(item => item.content)
		.flat()
		.filter(item => item.type === 'image');

const pickFileUpload = async () => {

	const files = await pickLocalFiles({ accept: ['.carddeck', '.json', '.csv'] });
	if (!files?.length) {
		return;
	}

	await loadBundleFile(files[0]);
};

const dropFileUpload = (event: DragEvent) => {
	state.dragging = false;
	if (!event.dataTransfer?.files.length) {
		return;
	}
	loadBundleFile(event.dataTransfer.files[0]);
};

const loadBundleFile = async (file: File) => {

	const extIdx = file.name.lastIndexOf('.');
	const ext = extIdx > 0 ? file.name.slice(extIdx) : file.name;

	switch (ext) {
		case ".carddeck":
			loadJSONGZFile(file);
			break;
		case ".json":
			await loadJSONFile(file);
			break;
		case ".csv":
			await loadFileCSV(file);
			break;
		default:
			return;
	}

	setTimeout(() => selectAllCards(), 50);
};

const loadJSONGZFile = async (file: File) => {

	state.source.name = file.name;

	const decompressor = new DecompressionStream('gzip');

	const bundle = await new Response(file.stream().pipeThrough(decompressor)).json().catch(() => null);
	if (!bundle) {
		state.source.error = 'Invalid JSON+GZ data source';
		return;
	}

	await parseBundleJSON(bundle);
};

const loadJSONFile = async (file: File) => {

	state.source.name = file.name;

	const bundle = await file.text().then(data => JSON.parse(data)).catch(() => null);
	if (!bundle) {
		state.source.error = 'Invalid JSON data source';
		return;
	}

	await parseBundleJSON(bundle);
};

const parseBundleJSON = async (bundle: ContentBundle | null) => {

	const cards: CardNode[] = [];
	const meta: Meta = { name: null, description: null, };

	if (bundle?.decks?.length) {

		for (const deck of bundle.decks) {

			if (!deck.deck_id && !deck.collection_id) {
				console.warn('Importer: Empty deck id bindings');
				continue;
			}

			if (!deck?.cards?.length) {
				console.warn('Importer: Empty deck content');
				continue;
			}

			for (const card of deck.cards) {

				if (!card?.id || (!card.front?.content?.length && !card.back?.content?.length)) {
					console.warn('Importer: Empty card content');
					return;
				}

				cards.push({
					id: card.id,
					front: card.front,
					back: card.back,
				});
			}

			if (!meta.name && deck.name.length) {
				meta.name = deck.name;
			}

			if (!meta.description && deck.description?.length) {
				meta.description = deck.description;
			}
		}
	}

	if (!cards.length) {
		state.source.error = 'Empty JSON data source';
		return;
	}

	const images: ImageData[] = [];

	if (bundle?.image_blobs?.length) {

		for (const image of bundle.image_blobs) {

			if (typeof image?.media_id !== 'string' || !image.media_id.length) {
				console.warn('Importer: Invalid image bundle: missing image ID');
				continue;
			} else if (typeof image?.data_url !== 'string' || !image.data_url.length) {
				console.warn('Importer: Invalid image bundle: missing image data');
				continue;
			}

			const blob = await fetch(image.data_url).then(res => res.blob()).catch(() => null);
			if (!blob) {
				console.error('Importer: Unable to parse image data url for image', image.media_id);
				continue;
			}

			images.push({
				media_id: image.media_id,
				media_url: URL.createObjectURL(blob),
				source_name: image.source_name || `Image upload ${new Date().getTime()}`,
				blob,
			});
		}
	}

	const imageUrlMap = new Map(images.map(item => ([ item.media_id, item.media_url ])));
	filterImageNodes(cards).forEach(item => item.media_url = imageUrlMap.get(item.media_id || ''));

	state.content = { cards, images, meta };
	state.options = { setMetadata: { available: true, value: true }, overwriteContent: true };
	state.source.ready = true;
};

const loadFileCSV = async (file: File) => {

	state.source.name = file.name;

	const { data, errors } = parse<Partial<CardContentCSVRow>>(await file.text(), { header: true, skipEmptyLines: true });
	if (!data || errors.length) {
		state.source.error = errors.at(0)?.message || 'Invalid CSV data source';
		return;
	}

	const cards: CardNode[] = [];

	interface CardNodeCandidate extends CardNode {
		valid: boolean;
	};

	for (const row of data) {

		const next: CardNodeCandidate = { id: crypto.randomUUID(), front: { content: [] }, back: { content: [] }, valid: false };

		if (typeof row.front_title === 'string' && row.front_title.length) {
			next.valid = true;
			next.front.content.push({ type: 'title', content: row.front_title });
		}

		if (typeof row.front_image === 'string' && row.front_image.length) {
			next.valid = true;
			next.back.content.push({ type: 'image', media_id: row.front_image });
		}

		if (typeof row.front_textarea === 'string' && row.front_textarea.length) {
			const content = parseTextBoxContent(row.front_textarea);
			if (content.length) {
				next.valid = true;
				next.front.content.push({ type: 'textbox', content });
			}
		}

		if (typeof row.front_quiz === 'string' && row.front_quiz.length) {
			const options = parseQuizOptions(row.front_quiz);
			if (options.length) {
				next.valid = true;
				next.front.content.push({ type: 'poll', content: options, is_quiz: true });
			}
		}

		if (typeof row.back_title === 'string' && row.back_title.length) {
			next.valid = true;
			next.back.content.push({ type: 'title', content: row.back_title });
		}

		if (typeof row.back_image === 'string' && row.back_image.length) {
			next.valid = true;
			next.back.content.push({ type: 'image', media_id: row.back_image });
		}

		if (typeof row.back_textarea === 'string' && row.back_textarea.length) {

			const content = parseTextBoxContent(row.back_textarea);
			if (content.length) {
				next.valid = true;
				next.back.content.push({ type: 'textbox', content });
			}
		}

		if (!next.valid) {
			continue;
		}

		cards.push({ id: next.id, front: next.front, back: next.back });
	}

	if (!cards.length) {
		state.source.error = 'Empty CSV data source';
		return;
	}

	state.content = { cards, images: [], meta: { name: null, description: null } };
	state.options = { setMetadata: { available: false, value: false }, overwriteContent: true };
	state.source.ready = true;
};

const importData = async () => {

	state.busy = true;

	const imageNodes = filterImageNodes(selectedCards.value);

	state.total = imageNodes.length + 1;

	const imageBlobs = new Map<string, ImageData>(state.content.images.map(item => ([ item.media_id, item ])));
	const imageRemap = new Map<string, string>();

	for (const node of imageNodes) {

		state.progress++;

		if (state.cancelled) {
			return;
		}

		if (!node.media_id) {
			continue;
		}

		const remappedMediaID = imageRemap.get(node.media_id);
		if (remappedMediaID) {
			node.media_id = remappedMediaID;
			continue;
		}

		const image = imageBlobs.get(node.media_id);
		if (!image) {
			console.warn('Importer: Referenced image media id not found', node.media_id);
			continue;
		}

		const { data, error } = await client.images.upload(image.blob, image.source_name);
		if (!data || error) {
			console.error('Importer: Unable to upload image', node.media_id, error?.message);
			continue;
		}

		node.media_id = data.id;
		node.media_url = null;

		imageRemap.set(node.media_id, data.id);
	}

	if (state.options.overwriteContent) {
		emit('replaceCards', selectedCards.value);
	} else {
		emit('addCards', selectedCards.value);
	}

	if (state.options.setMetadata.value) {
		emit('updateMeta', state.content.meta);
	}

	state.busy = false;

	exitTool();
};

</script>

<template>
	<EditorModal title="Import deck content" @close="exitTool">
		<div class="importer-tool">

			<OverlayErrorMessage v-if="state.source.error">

				Import failed

				<template v-slot:details>
					{{ state.source.error }}
				</template>

			</OverlayErrorMessage>

			<div v-else-if="state.source.ready" class="content-selection">

				<div class="selector-actions">
					<template v-if="state.selectedCards.size === state.content.cards.length">
						<GenericButton variant="thin" theme="orange" :disabled="state.busy" @click="clearCardSelection">
							Clear selection
						</GenericButton>
					</template>
					<template v-else>
						<GenericButton variant="thin" theme="blue" @click="selectAllCards">
							Select all
						</GenericButton>
					</template>
					<div class="selection-stats">
						{{ state.selectedCards.size }} item(s) selected
					</div>
				</div>

				<EditorCardSelectorGrid v-if="state.content.cards.length" :cards="state.content.cards" v-model="state.selectedCards" />
				<div v-else class="placeholder">No cards to display</div>

			</div>

			<div v-else class="file-drop-zone">

				<button class="drop-target"
					:class=" { active: state.dragging, busy: sourceLoading }"
					@dragover.prevent="state.dragging = true"
					@dragleave="state.dragging = false"
					@drop.prevent="dropFileUpload"
					@click="pickFileUpload">

					<div class="target-icon"></div>
					<div class="target-title">
						<template v-if="sourceLoading">
							Loading...
						</template>
						<template v-else>
							Drop files here
						</template>
					</div>

				</button>

			</div>

			<div class="options-selection">

				<div class="options-group">
					<div class="title">
						Source
					</div>
					<div class="value masked">
						{{ state.source.name || 'None' }}
					</div>
					<div v-if="state.content.meta.name" class="value">
						{{ state.content.meta.name }}
					</div>
					<div v-if="state.content.meta.description" class="value">
						{{ state.content.meta.description }}
					</div>
				</div>

				<div class="options-group">

					<div class="title">
						Content
					</div>

					<GenericToggle label="Replace existing metadata"
						v-model="state.options.setMetadata.value"
						:disabled="!state.options.setMetadata.available" />

					<GenericToggle label="Replace existing cards"
						v-model="state.options.overwriteContent" />

				</div>

				<InlineErrorMessage v-if="state.error">
					{{ state.error }}
				</InlineErrorMessage>

				<InlineProgressIndicator v-if="state.busy" title="Exporting cards" :total="state.total" :done="state.progress" />

				<GenericButton variant="thin" :disabled="!stateValid || state.busy" :spinner="state.busy" @click="importData">
					Import data
				</GenericButton>

			</div>
		</div>
	</EditorModal>
</template>

<style lang="scss" scoped>
	.importer-tool {
		display: grid;
		grid-template-columns: 3fr 1fr;
		gap: 3rem;
		padding: 0 2rem;
		height: 100%;
		min-height: 0;

		.content-selection {
			display: flex;
			flex-direction: column;
			gap: 1.5rem;
			min-height: 0;

			.selector-actions {
				display: flex;
				flex-flow: row nowrap;
				align-items: center;
				gap: 0.5rem;
				justify-content: space-between;

				.selection-stats {
					font-size: 0.75rem;
				}
			}

			.placeholder {
				display: flex;
				align-items: center;
				justify-content: center;
				flex-grow: 1;
				font-size: 0.8rem;
				color: var(--app-theme-mysterious-white);
			}
		}

		.file-drop-zone {
			display: flex;
			align-items: center;
			justify-content: center;
			flex-grow: 1;

			.drop-target {
				display: flex;
				flex-direction: column;
				gap: 2rem;
				align-items: center;
				background: unset;
				background-color: var(--app-theme-ghostly-glow);
				border-radius: 0.5rem;
				padding: 2rem 4rem;
				border: none;
				outline: none;
				transition: all 100ms ease;
				user-select: none;

				.target-icon {
					display: block;
					width: 4rem;
					height: 4rem;
					flex-shrink: 0;
					background-color: var(--app-theme-mysterious-white);
					mask-type: alpha;
					mask-image: url(/src/assets/icons/file-mask.svg);
					mask-position: center;
					mask-size: contain;
					mask-repeat: no-repeat;
				}

				.target-title {
					font-size: 1.25rem;
					font-weight: 300;
				}

				&:hover, &.active {
					cursor: pointer;
					background-color: var(--app-theme-powder-trail);
				}

				&.busy {
					cursor: not-allowed;
					pointer-events: none;
				}
			}
		}

		.options-selection {
			display: flex;
			flex-direction: column;
			gap: 1rem;
			min-width: 0;

			.options-group {
				display: flex;
				flex-direction: column;
				gap: 0.5rem;

				.title {
					font-size: 0.65rem;
					font-weight: 600;
				}

				.value {
					font-size: 0.65rem;
					font-weight: 600;
					padding: 0.125rem 0.5rem;
					overflow: hidden;
					text-overflow: ellipsis;
					white-space: nowrap;

					&.masked {
						background-color: var(--app-theme-mysterious-white);
						color: var(--app-theme-carbon);
						border-radius: 0.25rem;
					}
				}
			}
		}
	}
</style>
