<script setup lang="ts">
import { onMounted, reactive, watch } from 'vue';
import type { CardTextBoxElementNode, CardTextboxElementTheme, CardTextboxTextNode } from '../../../content';
import CardNodeHarness from './CardNodeHarness.vue';

const model = defineModel<CardTextBoxElementNode[]>();

const state = reactive({
	raw: '',
	editing: false,
});

const serializeFormat = (nodes: CardTextBoxElementNode[]) => {

	let result = '';

	for (const node of nodes) {
		switch (node.type) {
			case 'newline':
				result += '\n';
				continue;
			case 'text':
				result += serializeTextNode(node);
				continue;
		}
	}

	return result;
};

const escapeTextNodeContent = (content: string): string => {
	if (/[\|\`\s]/.test(content)) {
		return `\`${content.replaceAll('`', `'`)}\``;
	}
	return content;
};

const serializeTextNode = (node: CardTextboxTextNode): string => {

	const modifiers: string[] = [];

	if (node.theme?.bold) {
		modifiers.push('bold');
	}

	if (node.theme?.italic) {
		modifiers.push('italic');
	}

	if (node.theme?.decoration) {
		modifiers.push(node.theme.decoration);
	}

	if (node.theme?.highlight?.fill_color) {
		modifiers.push(`fill-${node.theme?.highlight.fill_color}`);
	}

	if (node.theme?.highlight?.text_color) {
		modifiers.push(`text-${node.theme?.highlight.text_color}`);
	}

	if (!modifiers.length) {
		return node.content;
	}

	return `{{ ${escapeTextNodeContent(node.content)} | ${modifiers.join(' | ')} }}`;
};

const parseTextNodeFormat = (content: string, modifiers?: string | null): CardTextboxTextNode => {

	const textColorPrefix = 'text-';
	const fillColorPrefix = 'fill-';

	const attributes = modifiers?.split('|').map(item => item.trim()).filter(item => item.length) || [];

	const theme: CardTextboxElementTheme = {};

	for (const attr of attributes) {

		switch (attr.toLowerCase()) {
			case 'bold':
				theme.bold = true;
				continue;
			case 'italic':
				theme.italic = true;
				continue;
			case 'underline':
				theme.decoration = 'underline';
				continue;
			case 'strikethrough':
				theme.decoration = 'strikethrough';
				continue;
		}

		if (attr.toLowerCase().startsWith(textColorPrefix)) {
			if (!theme.highlight) {
				theme.highlight = {};
			}
			theme.highlight.text_color = attr.slice(textColorPrefix.length);
			continue;
		}

		if (attr.toLowerCase().startsWith(fillColorPrefix)) {
			if (!theme.highlight) {
				theme.highlight = {};
			}
			theme.highlight.fill_color = attr.slice(fillColorPrefix.length);
			continue;
		}
	}

	return { type: 'text', content, theme };
};

const parseTextFormat = (value: string) => {

	const markExpr = /\{{2}\s*(([^\{\}\|\`]+)|(\`[^\`]+\`))((\s*\|\s*[\(\)a-z-0-9\_\-]*)*)\s*\}{2}/gi;

	const nodes: CardTextBoxElementNode[] = [];

	const lines = value.replaceAll('\r\n', '\n').replaceAll('\r', '\n').split('\n');

	for (const line of lines) {

		if (!line) {
			nodes.push({ type: 'newline' });
			continue;
		}

		if (nodes.length) {
			nodes.push({ type: 'newline' });
		}

		const exprs = line.matchAll(markExpr);

		let lastExprIdx = 0;

		for (const expr of exprs) {

			if (!expr.length) {
				throw new Error('Invalid regexp result');
			}

			nodes.push({ type: 'text', content: line.slice(lastExprIdx, expr.index) });
			lastExprIdx = expr.index + expr[0].length;

			const textContent = expr.at(3)?.slice(1, -1) || expr.at(2)?.trim() || null;
			if (!textContent?.length) {
				continue;
			}

			nodes.push(parseTextNodeFormat(textContent, expr[4]));
		}

		if (lastExprIdx < line.length) {
			nodes.push({ type: 'text', content: line.slice(lastExprIdx) });
		}
	}

	return nodes;
};

const serializeModel = () => state.raw = model.value ? serializeFormat(model.value) : '';

onMounted(() => {
	serializeModel();
	watch(() => model.value, () => !state.editing ? serializeModel() : void 0);
	watch(() => state.raw, (val) => state.editing ? model.value = parseTextFormat(val) : void 0);
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
			<textarea type="text" placeholder="Text" v-model="state.raw" @focus="state.editing = true" @blur="state.editing = false"></textarea>
		</template>
	</CardNodeHarness>
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
