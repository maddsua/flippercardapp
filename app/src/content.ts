
export interface CardNode {
	id: string;
	front: CardContentFace;
	back: CardContentFace;
};

export interface CardContentFace {
	theme?: CardFaceTheme;
	content: CardContentNode[];
};

export interface CardFaceTheme {
	card?: CardCanvasTheme | null;
	interactives?: CardContentElementTheme | null;
};

export interface CardCanvasTheme extends CardContentElementTheme {
	outline_color?: string | null;
};

export interface CardContentElementTheme {
	fill_color?: string | null;
	mask_color?: string | null;
}

export type CardContentNode = CardTitleNode | CardTextBoxNode | CardPollNode | CardImageNode;

interface CardContentNodeBase {
	type: string;
};

export interface CardTitleNode extends CardContentNodeBase {
	readonly type: 'title';
	content: string;
};

export interface CardTextBoxNode extends CardContentNodeBase {
	readonly type: 'textbox';
	content: CardTextBoxElementNode[];
};

export type CardTextBoxElementNode = CardTextboxTextNode | CardTextboxNewlineNode;

export interface CardTextboxTextNode extends CardContentNodeBase {
	readonly type: 'text';
	content: string;
	theme?: CardTextboxElementTheme;
};

export interface CardTextboxElementTheme {
	highlight?: CardTextboxElementTextHighlight;
	bold?: boolean | null;
	italic?: boolean | null;
	decoration?: 'underline' | 'strikethrough' | null;
	size?: TextSize;
}

type TextSize = 'xs' | 's' | 'm' | 'l' | 'xl';

const textSizeSet = new Set<string>(Object.keys({
	xs: null,
	s: null,
	m: null,
	l: null,
	xl: null,
} satisfies Record<TextSize, null>));

export interface CardTextboxElementTextHighlight {
	text_color?: string | null;
	fill_color?: string | null;
}

export interface CardTextboxNewlineNode extends CardContentNodeBase {
	readonly type: 'newline';
};

export interface CardPollNode extends CardContentNodeBase {
	readonly type: 'poll';
	is_quiz?: boolean | null;
	content: CardPollNodeOption[];
};

export interface CardPollNodeOption {
	value: string;
	is_answer?: boolean | null;
};

export interface CardImageNode {
	readonly type: 'image';
	media_id?: string | null;
	media_url?: string | null;
	state?: UploadReadyState | null;
};

export enum UploadReadyState {
	Idle,
	Uploading,
	Done
};

export interface ContentBundle {
	decks?: DeckContentBundle[];
	image_blobs?: ImageBlobBundle[];
};

export interface DeckContentBundle {
	deck_id: string;
	collection_id: string;
	name: string;
	description: string | null;
	cards: CardNode[];
};

export interface ImageBlobBundle {
	media_id: string;
	source_name: string;
	data_url: string;
};

export interface CardContentCSVRow {
	front_title: string;
	front_image: string;
	front_textarea: string;
	front_quiz: string;
	back_title: string;
	back_image: string;
	back_textarea: string;
};

export const stringifyTextBoxContent = (nodes: CardTextBoxElementNode[]): string => {

	let result = '';

	for (const node of nodes) {
		switch (node.type) {
			case 'newline':
				result += '\n';
				continue;
			case 'text':
				if (node.content.trim().length) {
					result += stringifyTextBoxNode(node);
				}
				continue;
		}
	}

	return result;
};

const escapeTextBoxNodeContent = (content: string): string => {
	if (/[\|\`\s]/.test(content)) {
		return `\`${content.replaceAll('`', `'`)}\``;
	}
	return content;
};

const stringifyTextBoxNode = (node: CardTextboxTextNode): string => {

	const modifiers: string[] = [];

	if (node.theme?.size) {
		modifiers.push(`size-${node.theme?.size}`);
	}

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

	return `{{ ${escapeTextBoxNodeContent(node.content)} | ${modifiers.join(' | ')} }}`;
};

const parseTextBoxNode = (content: string, modifiers?: string | null): CardTextboxTextNode => {

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

		const textColorAttr = prefixedAttributeValue(attr, 'text-');
		if (textColorAttr) {
			if (!theme.highlight) {
				theme.highlight = {};
			}
			theme.highlight.text_color = textColorAttr;
			continue;
		}

		const fillColorAttr = prefixedAttributeValue(attr, 'fill-');
		if (fillColorAttr) {
			if (!theme.highlight) {
				theme.highlight = {};
			}
			theme.highlight.fill_color = fillColorAttr;
			continue;
		}

		const sizeAttr = prefixedAttributeValue(attr, 'size-');
		if (sizeAttr && textSizeSet.has(sizeAttr)) {
			theme.size = sizeAttr as TextSize;
			continue;
		}
	}

	return { type: 'text', content, theme };
};

const prefixedAttributeValue = (attr: string, prefix: string): string | null => {
	if (!attr.toLowerCase().startsWith(prefix)) {
		return null;
	}
	return attr.slice(prefix.length) || null;
};

export const parseTextBoxContent = (value: string): CardTextBoxElementNode[] => {

	if (!value.trim().length) {
		return [];
	}

	const nodes: CardTextBoxElementNode[] = [];

	const markExpr = /\{{2}\s*(([^\{\}\|\`]+)|(\`[^\`]+\`))((\s*\|\s*[\(\)a-z-0-9\_\-]*)*)\s*\}{2}/gi;
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

			if (expr.length < 5) {
				throw new Error('Invalid regexp result');
			}

			nodes.push({ type: 'text', content: line.slice(lastExprIdx, expr.index) });
			lastExprIdx = expr.index + expr[0].length;

			const textContent = expr.at(3)?.slice(1, -1) || expr.at(2)?.trim() || null;
			if (!textContent?.length) {
				continue;
			}

			nodes.push(parseTextBoxNode(textContent, expr[4]));
		}

		if (lastExprIdx < line.length) {
			nodes.push({ type: 'text', content: line.slice(lastExprIdx) });
		}
	}

	return nodes;
};

export const parseQuizOptions = (value: string): CardPollNodeOption[] => {
	const tokens = value.split(value.includes('|') ? '|' : ',').map(item => item.trim()).filter(item => item.length);
	return tokens.map((item, idx) => ({ is_answer: idx === 0, value: item }))
};

type NodeType = CardContentNode['type'];
const nodeSortOrder: NodeType[] = ['title', 'image', 'textbox', 'poll'];
const nodeSortOrderMap = new Map<NodeType, number>(nodeSortOrder.map((item, idx) => ([item, idx])));

const createNode = (type: NodeType): CardContentNode => {
	switch (type) {
		case 'title':
			return { type: 'title', content: '' };
		case 'image':
			return { type: 'image', state: UploadReadyState.Idle };
		case 'textbox':
			return { type: 'textbox', content: [] };
		case 'poll':
			return { type: 'poll', content: [] };
		default:
			throw new Error(`Unknown node type: ${type}`);
	}
};

export const addModelNode = (target: CardContentNode[] | null | undefined, type: NodeType) => {

	if (!target) {
		throw new Error('Unable to insert node into a non-existent target');
	}

	const next = createNode(type);
	const nextPosition = nodeSortOrderMap.get(type) ?? nodeSortOrderMap.size;

	for (let idx = 0; idx < target.length; idx++) {
		const node = target[idx];
		const nodePosition = nodeSortOrderMap.get(node.type) ?? 0;
		if (nodePosition > nextPosition) {
			target.splice(idx, 0, next);
			return;
		}
	}

	target.push(next);
};
