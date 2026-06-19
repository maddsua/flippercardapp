<script setup lang="ts">
import { ref } from 'vue';
import type { CardContentFace } from '@/content';
import { addModelNode } from '@/content';
import FullscreenMessage from '@/components/App/Messages/FullscreenMessage.vue';
import CardThumbnail from '../CardThumbnail.vue';
import CardFaceEditorPreviewOverlay from './FacePreview/CardFaceEditorPreviewOverlay.vue';
import EditablePollNode from './NodeEditor/EditablePollNode.vue';
import EditableTextNode from './NodeEditor/EditableTextNode.vue';
import EditableImageNode from './NodeEditor/EditableImageNode.vue';
import EditableTitleNode from './NodeEditor/EditableTitleNode.vue';
import EditorGroup from './EditorGroup.vue';
import GenericButton from '@/components/App/Inputs/GenericButton.vue';

const model = defineModel<CardContentFace | null>();

const props = defineProps<{
	isFront?: boolean;
	active?: boolean;
}>();

const previewHovered = ref(false);

const removeNode = (idx: number) => model.value?.content?.splice(idx, 1);

const reorderNode = (idx: number, delta: number) => {

	const newIdx = idx + delta;
	if (!model.value || newIdx < 0 || newIdx >= (model.value.content?.length ?? 0)) {
		return;
	}

	const node = model.value.content[idx];

	model.value.content[idx] = model.value.content[newIdx];
	model.value.content[newIdx] = node;
};

</script>

<template>
	<div class="face-editor" :class="{ active }">

		<div class="face-editor-header">

			<div class="face-editor-summary">

				<div class="face-editor-title">
					<template v-if="props.isFront">
						Front
					</template>
					<template v-else>
						Back
					</template>

					face
				</div>

			</div>

			<div class="preview-trigger-attach">
				<CardThumbnail v-if="model"
					size="small"
					:face="model"
					:hidden="!active"
					@mouseover="previewHovered = true"
					@mouseleave="previewHovered = false" />
			</div>

		</div>

		<hr />

		<div v-if="model" class="face-editor-canvas">

			<CardFaceEditorPreviewOverlay v-if="model && (previewHovered || !active)" :face="model" />

			<EditorGroup title="Content" :scrollable="true">

				<div class="content-tree">
					<template v-for="(item, idx) of model.content">

						<EditableTitleNode v-if="item.type === 'title'"
							v-model="item.content"
							@up="reorderNode(idx, -1)"
							@down="reorderNode(idx, 1)"
							@remove="removeNode(idx)" />

						<EditableImageNode v-else-if="item.type === 'image'"
							v-model="item.media_id"
							@up="reorderNode(idx, -1)"
							@down="reorderNode(idx, 1)"
							@remove="removeNode(idx)" />

						<EditableTextNode v-else-if="item.type === 'textbox'"
							v-model="item.content"
							@up="reorderNode(idx, -1)"
							@down="reorderNode(idx, 1)"
							@remove="removeNode(idx)" />

						<EditablePollNode v-else-if="item.type === 'poll'"
							v-model="item.content"
							@up="reorderNode(idx, -1)"
							@down="reorderNode(idx, 1)"
							@setQuizFlag="flag => item.is_quiz = flag"
							@remove="removeNode(idx)" />
					</template>
				</div>

				<div v-if="!model.content?.length" class="add-actions">

					<GenericButton variant="thin" @click="addModelNode(model.content, 'title')">
						+ Add title
					</GenericButton>

					<GenericButton variant="thin" @click="addModelNode(model.content, 'image')">
						+ Add image
					</GenericButton>

					<GenericButton variant="thin" @click="addModelNode(model.content, 'textbox')">
						+ Add textbox
					</GenericButton>

				</div>

			</EditorGroup>

		</div>

		<FullscreenMessage v-else>
			No page selected
		</FullscreenMessage>

	</div>
</template>

<style lang="scss" scoped>
	.face-editor {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		min-height: 0;
		height: 100%;

		border-radius: 0.25rem;
		padding: 0.25rem 0.5rem;
		border: 2px solid transparent;

		&.active {
			border-color: var(--app-theme-sporty-yellow);
		}

		.face-editor-header {
			position: relative;
			display: flex;
			flex-flow: row nowrap;
			align-items: start;
			justify-content: space-between;
			gap: 2rem;
			padding: 0 1rem;

			.face-editor-summary {
				display: flex;
				flex-direction: column;
				gap: 1rem;

				.face-editor-title {
					font-size: 1.125rem;
					font-weight: 300;
				}
			}

			.preview-trigger-attach {
				position: absolute;
				top: 0;
				right: 0;
				z-index: 25;
			}
		}

		hr {
			display: block;
			background-color: var(--app-theme-ghostly-glow);
			width: 100%;
			height: 1px;
			outline: none;
			border: none;
		}

		.content-tree {
			display: flex;
			flex-direction: column;
			gap: 1rem;
			height: 100%;
		}

		.add-actions {
			display: flex;
			flex-direction: column;
			align-items: center;
			gap: 1.25rem;
		}

		.face-editor-canvas {
			position: relative;
			width: 100%;
			height: 100%;
			min-height: 0;
			padding: 1rem 0;
		}
	}
</style>
