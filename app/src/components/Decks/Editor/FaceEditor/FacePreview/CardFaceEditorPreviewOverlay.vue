<script setup lang="ts">
import type { CardContentFace } from '@/content';
import CardFace from '@/components/Cards/CardFace.vue';

const props = defineProps<{
	face: CardContentFace;
}>();

</script>

<template>
	<div class="editor-preview-overlay">
		<div class="preview-slot">
			<CardFace :entry="face" />
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.editor-preview-overlay {
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		border-radius: 1rem;
		padding: 1rem;
		z-index: 20;
		display: flex;
		align-items: center;
		justify-content: center;
		background-color: var(--app-theme-carbon);

		.preview-slot {
			position: relative;
			width: 70%;
			container-type: inline-size;
			aspect-ratio: 0.6667;
		}

		&:hover {
			cursor: pointer;
		}
	}
</style>
