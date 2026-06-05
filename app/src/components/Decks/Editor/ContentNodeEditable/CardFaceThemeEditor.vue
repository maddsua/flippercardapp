<script setup lang="ts">

import { computed, onMounted } from 'vue';
import type { CardFaceTheme } from '@/content';
import CardColorSwatch from './CardColorSwatch.vue';

const model = defineModel<CardFaceTheme>();

onMounted(() => {
	if (!model.value) {
		throw new Error('Editor model binding may not be omitted');
	}
});

const cardTheme = computed(() => model.value?.card || {});
const interactivesTheme = computed(() => model.value?.interactives || {});

</script>

<template>
	<div class="card-face-theme-editor">

		<div class="title">
			Face theme
		</div>

		<div class="swatches">
			<CardColorSwatch label="Card fill" v-model="cardTheme.fill_color" />
			<CardColorSwatch label="Card masks (text and elements)" v-model="cardTheme.mask_color" />
			<CardColorSwatch label="Card outlines" v-model="cardTheme.outline_color" />
		</div>
		
		<div class="swatches">
			<CardColorSwatch label="Element fill" v-model="interactivesTheme.fill_color" />
			<CardColorSwatch label="Element fill (text)" v-model="interactivesTheme.mask_color" />
		</div>

	</div>
</template>

<style lang="scss" scoped>
	.card-face-theme-editor {
		display: flex;
		flex-direction: column;
		gap: 2rem;
		height: 100%;
		min-height: 0;

		.title {
			font-size: 0.75rem;
			font-weight: 600;
		}

		.swatches {
			display: flex;
			flex-direction: column;
			gap: 1rem;
		}
	}
</style>
