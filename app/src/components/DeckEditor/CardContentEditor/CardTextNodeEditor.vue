<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';
import type { CardTextboxElementNewlineNode, CardTextboxElementTextNode } from '../../../content';
import CardNodeHarness from './CardNodeHarness.vue';

const model = defineModel<Array<CardTextboxElementTextNode | CardTextboxElementNewlineNode>>();

//	todo: make it the proper way
//	todo: implement text formatting

const rawValue = ref('');

const modelToRaw = () => {
	return model.value?.map(item => item.type === 'text' ? item.content : '').join(' ') || '';
};

const rawToModel = (): CardTextboxElementTextNode[] => {
	return rawValue.value.length ? [{ type: 'text', content: rawValue.value }] : [];
};

onMounted(() => {
	rawValue.value = modelToRaw();
	watch(() => model.value, () => rawValue.value = modelToRaw());
	watch(() => rawValue.value, () => model.value = rawToModel());
});

const emit = defineEmits<{
	(e: 'remove'): void;
}>();

</script>

<template>
	<CardNodeHarness @remove="emit('remove')">
		<template v-slot:title>
			Text
		</template>

		<template v-slot:content>
			<textarea type="text" placeholder="Text" v-model="rawValue"></textarea>
		</template>
	</CardNodeHarness>
</template>

<style lang="scss" scoped>
	textarea {
		display: block;
		max-width: 100%;
		width: 25rem;
		height: 10rem;
		resize: none;
		background-color: unset;
		border: none;
		outline: none;
		color: var(--app-theme-snow-white);
		font-size: 0.85rem;

		&::placeholder {
			color: var(--app-theme-mysterious-white);
		}
	}
</style>
