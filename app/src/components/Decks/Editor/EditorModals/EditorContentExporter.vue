<script setup lang="ts">
import { unparse } from 'papaparse';
import { computed, onMounted, onUnmounted, reactive } from 'vue';
import { useClient, type Result } from '@/api';
import type { ResourceVisibility } from '@/api_models';
import {
	stringifyTextBoxContent,
	type CardContentCSVRow,
	type CardImageNode,
	type CardNode,
	type CardPollNode,
	type CardTextBoxNode,
	type ContentBundle,
	type ImageBlobBundle
} from '@/content';
import { blobToJson, downloadBlob, escapeFileName } from '@/files';
import GenericButton from '@/components/App/Inputs/GenericButton.vue';
import GenericToggle from '@/components/App/Inputs/GenericToggle.vue';
import GenericDropdown from '@/components/App/Inputs/GenericDropdown.vue';
import GenericInput from '@/components/App/Inputs/GenericInput.vue';
import InlineErrorMessage from '@/components/App/Messages/InlineErrorMessage.vue';
import InlineProgressIndicator from '@/components/App/Messages/InlineProgressIndicator.vue';
import EditorCardSelectorGrid from '../EditorCardSelectorGrid.vue';
import EditorModal from '../EditorModal.vue';

const client = useClient();

const props = defineProps<{
	content: Content;
	meta: Meta;
}>();

interface Meta {
	deckID: string | null;
	collectionID: string | null;
};

interface ContentMeta {
	name: string;
	description: string | null;
	visibility: ResourceVisibility;
};

interface Content {
	meta: ContentMeta;
	cards: CardNode[];
};

const emit = defineEmits<{
	(e: 'done'): void;
}>();

enum Format {
	JSON = 'json',
	JSONGZ = 'json+gz',
	CSV = 'csv',
};

interface FormatOption {
	value: Format;
	label: string;
	flags?: {
		media?: boolean;
		metadata?: boolean;
	};
};

const formatOptions: FormatOption[] = [
	{
		value: Format.JSONGZ,
		label: 'JSON bundle (compressed)',
		flags: {
			media: true,
			metadata: true
		}
	},
	{
		value: Format.JSON,
		label: 'JSON bundle',
		flags: {
			media: true,
			metadata: true
		}
	},
	{
		value: Format.CSV,
		label: 'CSV (text only)'
	},
];

const activeFormatOption = computed(() => formatOptions.find(item => item.value === state.options.format));

const state = reactive({
	selectedCards: new Set<string>(),
	options: {
		format: Format.JSONGZ,
		exportImages: true,
		exportDetails: true,
		filename: `card-export-${new Date().getTime()}`,
	},
	progress: 0,
	total: 0,
	busy: false,
	cancelled: false,
	error: null as string | null,
});

const selectedCards = computed(() => props.content.cards.filter(item => state.selectedCards.has(item.id)));
const stateValid = computed(() => selectedCards.value.length > 0 && state.options.filename.length > 0 && !state.error && !state.cancelled);

const clearCardSelection = () => state.selectedCards.clear();

const selectAllCards = () => {

	clearCardSelection();

	for (const card of props.content.cards) {
		state.selectedCards.add(card.id);
	}
};

const exitTool = () => {
	if (state.busy) {
		state.error = 'Export cancelled by user';
		console.debug(state.error);
		state.cancelled = true;
	}
	emit('done');
};

const handleKeys = (event: KeyboardEvent) => {
	const key = event.key.toLowerCase();
	if (key === 'escape' || key === 'esc') {
		event.stopImmediatePropagation();
		event.stopPropagation();
		exitTool();
	}
};

onMounted(() => {
	selectAllCards();
	state.options.filename = `${escapeFileName(props.content.meta.name)}-export-${new Date().getTime()}`;
	document.addEventListener('keydown', handleKeys);
});

onUnmounted(() => document.removeEventListener('keydown', handleKeys));

const exportDeck = async () => {

	state.busy = true;
	state.progress = 0;
	state.total = 1;

	switch (state.options.format) {
		case Format.JSONGZ:
			await exportDeckJSON({ compress: true });
			break;
		case Format.JSON:
			await exportDeckJSON();
			break;
		case Format.CSV:
			await exportDeckCSV();
			break;
	}

	state.progress = 0;
	state.total = 0;
	state.busy = false;
	state.cancelled = false;

	exitTool();
};

const exportDeckJSON = async (opts?: { compress?: boolean }) => {

	const { name, description } = props.content.meta;

	const contentNodes = selectedCards.value.map(item => [item.front, item.back])
		.flat().map(item => item.content).flat();

	state.total = contentNodes.length;

	const imageBlobs: ImageBlobBundle[] = [];

	for (const node of contentNodes) {

		if (state.cancelled) {
			console.debug('Export cancelled');
			return;
		}

		state.progress++;

		if (node.type !== 'image') {
			continue;
		}

		const { data: blob, error } = await downloadImageBlob(node);
		if (!blob || error) {
			console.error('Download image blob:', error?.message || 'Unable to download image');
			continue;
		}

		imageBlobs.push(blob);
	}

	if (state.cancelled) {
		console.debug('Export cancelled');
		return;
	}

	const bundle: ContentBundle = {
		decks: [{
			deck_id: props.meta.deckID || 'unknown',
			collection_id: props.meta.collectionID || 'unknown',
			name,
			description,
			cards: selectedCards.value
		}],
		image_blobs: imageBlobs.length ? imageBlobs : undefined,
	};

	state.progress = state.total;

	const rawBlob = new Blob([JSON.stringify(bundle)]);

	if (opts?.compress) {
		const compressor = new CompressionStream('gzip');
		const compressedBlob = await new Response(rawBlob.stream().pipeThrough(compressor)).blob();
		downloadBlob(compressedBlob, `${state.options.filename}.carddeck`);
		return;
	}
	
	downloadBlob(rawBlob, `${state.options.filename}.json`);
};

const downloadImageBlob = async (node: CardImageNode): Promise<Result<ImageBlobBundle>> => {

	if (node.type !== 'image' || !node.media_id?.length) {
		return { data: null, error: new Error('Not an image or doesnt have a valid media ID') };
	}

	const { data, error: metadataError } = await client.images.metadata(node.media_id);
	if (!data || metadataError) {
		return { data: null, error: metadataError };
	}

	const { blob, error: blobError } = await client.images.blob(node.media_id);
	if (!blob || blobError) {
		return { data: null, error: blobError };
	}

	return {
		data: {
			media_id: node.media_id,
			source_name: data.source_name,
			data_url: await blobToJson(blob),
		},
		error: null,
	};
};

const exportDeckCSV = async () => {

	const serializePoll = (node?: CardPollNode | null): string => {

		if (!node?.content.length) {
			return '';
		}

		return [
			node.content.find(item => item.is_answer),
			...node.content.filter(item => !item.is_answer),
		].filter(item => item?.value).map(item => item?.value).join(',');
	};

	const serializeTextbox = (node?: CardTextBoxNode | null): string => {
		if (!node?.content.length) {
			return '';
		}
		return stringifyTextBoxContent(node.content);
	};

	const blob = new Blob([unparse(props.content.cards.map(item => ({
		front_title: item.front.content.find(item => item.type === 'title')?.content || '',
		front_image: item.front.content.find(item => item.type === 'image')?.media_id || '',
		front_textarea: serializeTextbox(item.front.content.find(item => item.type === 'textbox')),
		front_quiz: serializePoll(item.front.content.find(item => item.type === 'poll')),
		back_title: item.back.content.find(item => item.type === 'title')?.content || '',
		back_image: item.back.content.find(item => item.type === 'image')?.media_id || '',
		back_textarea: serializeTextbox(item.back.content.find(item => item.type === 'textbox')),
	} satisfies CardContentCSVRow)))]);

	downloadBlob(blob, `${state.options.filename}.csv`);
};

</script>

<template>
	<EditorModal title="Export deck content" @close="exitTool">
		<div class="exporter-tool">

			<div class="content-selection">

				<div class="selector-actions">
					<template v-if="state.selectedCards.size === props.content.cards.length">
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

				<EditorCardSelectorGrid v-if="content.cards.length" :cards="content.cards" v-model="state.selectedCards" />
				<div v-else class="placeholder">No cards to display</div>

			</div>

			<div class="options-selection">

				<div class="options-group">
					<div class="title">
						Export format
					</div>
					<GenericDropdown :options="formatOptions" v-model="state.options.format" />
				</div>

				<div class="options-group">
					<div class="title">
						Export file name
					</div>
					<GenericInput v-model="state.options.filename" :disabled="state.busy" />
				</div>

				<div class="options-group">
					<div class="title">
						Deck details
					</div>
					<GenericToggle label="Include deck details" v-model="state.options.exportDetails" :disabled="!activeFormatOption?.flags?.metadata" />
				</div>

				<div class="options-group">
					<div class="title">
						Media
					</div>
					<GenericToggle label="Include images" v-model="state.options.exportImages" :disabled="!activeFormatOption?.flags?.media" />
				</div>

				<InlineErrorMessage v-if="state.error">
					{{ state.error }}
				</InlineErrorMessage>

				<InlineProgressIndicator v-if="state.busy" title="Exporting cards" :total="state.total" :done="state.progress" />

				<GenericButton variant="thin" :disabled="!stateValid || state.busy" :spinner="state.busy" @click="exportDeck">
					Export data
				</GenericButton>

			</div>
		</div>
	</EditorModal>
</template>

<style lang="scss" scoped>
	.exporter-tool {
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

		.options-selection {
			display: flex;
			flex-direction: column;
			gap: 1rem;

			.options-group {
				display: flex;
				flex-direction: column;
				gap: 0.5rem;

				.title {
					font-size: 0.65rem;
					font-weight: 600;
				}
			}
		}
	}
</style>
