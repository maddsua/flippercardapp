<script setup lang="ts">

import { computed } from 'vue';
import type { CardFaceTheme } from '@/content';
import CardColorSwatch from './CardColorSwatch.vue';
import EditorGroup from '../EditorGroup.vue';

const model = defineModel<CardFaceTheme>();

const cardTheme = computed(() => model.value?.card || {});
const interactivesTheme = computed(() => model.value?.interactives || {});

</script>

<template>

	<EditorGroup title="Card theme">
		<CardColorSwatch label="Card fill color" v-model="cardTheme.fill_color" />
		<CardColorSwatch label="Text color" v-model="cardTheme.mask_color" />
		<CardColorSwatch label="Card outline color" v-model="cardTheme.outline_color" />
	</EditorGroup>

	<EditorGroup title="Interactive elements' theme">
		<CardColorSwatch label="Element fill" v-model="interactivesTheme.fill_color" />
		<CardColorSwatch label="Element text" v-model="interactivesTheme.mask_color" />
	</EditorGroup>

</template>
