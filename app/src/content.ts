
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

export type CardContentNode = CardTitleElement | CardTextBoxElement | CardPollElement | CardImageElement;

interface BaseCardContentElement {
	type: string;
};

export interface CardTitleElement extends BaseCardContentElement {
	readonly type: 'title';
	content: string;
};

export interface CardTextBoxElement extends BaseCardContentElement {
	readonly type: 'textbox';
	content: CardTextBoxElementContentNode[];
};

export type CardTextBoxElementContentNode = CardTextboxElementTextNode | CardTextboxElementNewlineNode;

export interface BaseCardTextboxElement extends BaseCardContentElement {};

export interface CardTextboxElementTextNode extends BaseCardTextboxElement {
	readonly type: 'text';
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
	readonly type: 'newline';
};

export interface CardPollElement extends BaseCardContentElement {
	readonly type: 'poll';
	is_quiz?: boolean;
	content: CardPollElementOptionNode[];
};

export interface CardPollElementOptionNode {
	value: string;
	is_answer?: boolean;
};

export interface CardImageElement {
	readonly type: 'image';
	media_id?: string | null;
	state?: UploadReadyState | null;
};

export enum UploadReadyState {
	Idle,
	Uploading,
	Done
};
