<script setup lang="ts">
import { onMounted, reactive, watch } from 'vue';
import { parseTextBoxContent, stringifyTextBoxContent, type CardTextBoxElementNode } from '@/content';
import EditableNodeHarness from './EditableNodeHarness.vue';
import { blurInteractive } from '@/dom';

const model = defineModel<CardTextBoxElementNode[]>();

const state = reactive({
	raw: '',
	editing: false,
});

const serializeModel = () => state.raw = model.value ? stringifyTextBoxContent(model.value) : '';

onMounted(() => {
	serializeModel();
	watch(() => model.value, () => !state.editing ? serializeModel() : void 0);
	watch(() => state.raw, (val) => state.editing ? model.value = parseTextBoxContent(val) : void 0);
});

</script>

<template>
	<EditableNodeHarness>

		<template v-slot:title>
			Text
		</template>

		<textarea type="text"
			placeholder="Text"
			v-model="state.raw"
			@focus="state.editing = true"
			@blur="state.editing = false"
			@keydown.escape.stop="blurInteractive">
		</textarea>

	</EditableNodeHarness>
</template>

<style lang="scss" scoped>
	textarea {
		display: block;
		max-width: 100%;
		width: 25rem;
		height: 13rem;
		resize: none;
		background-color: unset;
		border: none;
		outline: none;
		color: var(--app-theme-snow-white);
		font-size: 0.8rem;
		line-height: 1.5em;

		&::placeholder {
			color: var(--app-theme-mysterious-white);
		}
	}
</style>
