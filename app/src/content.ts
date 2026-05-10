
export interface CardNode {
	id: string;
	front: CardFace;
	back: CardFace;
};

export interface CardFace {
	theme?: CardFaceTheme;
	content: CardContentNode[];
};

export interface CardFaceTheme {
	card: CardCanvasTheme;
	interactives?: ElementTheme;
}

export interface CardCanvasTheme extends ElementTheme {
	outline_color?: string;
};

export interface ElementTheme {
	fill_color?: string;
	mask_color?: string;
}

//	todo: add images

export type CardContentNode = CardTitleNode | CardTextBoxNode | CardPollNode;

interface BaseContentNode {
	type: string;
};

export interface CardTitleNode extends BaseContentNode {
	type: 'title';
	content: string;
};

export interface CardTextBoxNode extends BaseContentNode {
	type: 'textbox';
	content: Array<TextNote | NewlineNode>;
};

export interface TextNote extends BaseContentNode {
	type: 'text';
	content: string;
	theme?: TextNodeTheme;
};

export interface TextNodeTheme {
	highlight?: TextHighlight;
	bold?: boolean;
	italic?: boolean;
	decoration?: 'underline' | 'strikethrough';
}

export interface TextHighlight {
	text_color: string;
	fill_color: string;
}

export interface NewlineNode extends BaseContentNode {
	type: 'newline';
};

export interface CardPollNode extends BaseContentNode {
	type: 'poll';
	is_quiz?: boolean;
	content: PollOption[];
};

export interface PollOption {
	value: string;
	is_answer?: boolean;
}
