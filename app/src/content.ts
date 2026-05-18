
export interface CardContentNode {
	id: string;
	front: CardContentFace;
	back: CardContentFace;
};

export interface CardContentFace {
	theme?: CardFaceTheme;
	content: CardContentElement[];
};

export interface CardFaceTheme {
	card?: CardCanvasTheme;
	interactives?: CardContentElementTheme;
};

export interface CardCanvasTheme extends CardContentElementTheme {
	outline_color?: string;
};

export interface CardContentElementTheme {
	fill_color?: string;
	mask_color?: string;
}

//	todo: add images

export type CardContentElement = CardTitleElement | CardTextBoxElement | CardPollElement;

interface BaseCardContentElement {
	type: string;
};

export interface CardTitleElement extends BaseCardContentElement {
	type: 'title';
	content: string;
};

export interface CardTextBoxElement extends BaseCardContentElement {
	type: 'textbox';
	content: Array<CardTextboxElementTextNode | CardTextboxElementNewlineNode>;
};

export interface BaseCardTextboxElement extends BaseCardContentElement {};

export interface CardTextboxElementTextNode extends BaseCardTextboxElement {
	type: 'text';
	content: string;
	theme?: CardTextboxElementTextNodeTheme;
};

export interface CardTextboxElementTextNodeTheme {
	highlight?: CardTextboxElementTextNodeHighlight;
	bold?: boolean;
	italic?: boolean;
	decoration?: 'underline' | 'strikethrough';
}

export interface CardTextboxElementTextNodeHighlight {
	text_color: string;
	fill_color: string;
}

export interface CardTextboxElementNewlineNode extends BaseCardTextboxElement {
	type: 'newline';
};

export interface CardPollElement extends BaseCardContentElement {
	type: 'poll';
	is_quiz?: boolean;
	content: CardPollElementOptionNode[];
};

export interface CardPollElementOptionNode {
	value: string;
	is_answer?: boolean;
};
