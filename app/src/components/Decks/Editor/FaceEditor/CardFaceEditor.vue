<script setup lang="ts">
import { reactive, watch } from 'vue';
import type { CardContentFace } from '@/content';
import CardFaceContentEditor from './NodeEditor/CardFaceContentEditor.vue';
import FullscreenMessage from '@/components/App/Messages/FullscreenMessage.vue';
import CardThumbnail from '../CardThumbnail.vue';
import CardFaceThemeEditor from './ThemeEditor/CardFaceThemeEditor.vue';
import CardFaceEditorPreviewOverlay from './FacePreview/CardFaceEditorPreviewOverlay.vue';

const model = defineModel<CardContentFace | null>();

watch(model, (value) => {
	if (value && !value.theme) {
		value.theme = value.theme || {};
		value.theme.card = value.theme.card || {};
		value.theme.interactives = value.theme.interactives || {};
	}
})

const props = defineProps<{
	isFront?: boolean;
}>();

const state = reactive({
	preview: false,
});

</script>

<template>
	<div class="face-editor">

		<div class="editor-header">

			<div class="editor-summary">

				<div class="editor-title">
					<template v-if="props.isFront">
						Front
					</template>
					<template v-else>
						Back
					</template>

					face
				</div>

			</div>

			<CardThumbnail v-if="model" size="small" :face="model" @mouseover="state.preview = true" @mouseleave="state.preview = false" />

		</div>

		<hr />

		<div v-if="model" class="editor-canvas">

			<CardFaceEditorPreviewOverlay v-if="model && state.preview" :face="model" />

			<div class="editor-scroll-area">
				<CardFaceContentEditor v-model="model.content" :isFront="props.isFront" />
				<CardFaceThemeEditor v-model="model.theme" />
			</div>

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

		.editor-header {
			display: flex;
			flex-flow: row nowrap;
			align-items: start;
			justify-content: space-between;
			gap: 2rem;
			padding: 0 1rem;

			.editor-summary {
				display: flex;
				flex-direction: column;
				gap: 1rem;
	
				.editor-title {
					font-size: 1.125rem;
					font-weight: 300;
				}
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

		.editor-canvas {
			position: relative;
			width: 100%;
			height: 100%;
			min-height: 0;
			padding: 1rem;
			
			.editor-scroll-area {
				display: flex;
				flex-direction: column;
				gap: 2rem;
				min-height: 0;
				height: 100%;
				overflow: hidden auto;
				scrollbar-width: thin;
			}
		}
	}
</style>
