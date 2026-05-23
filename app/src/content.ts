
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
	bold?: boolean;
	italic?: boolean;
	decoration?: 'underline' | 'strikethrough';
}

export interface CardTextboxElementTextHighlight {
	text_color: string;
	fill_color: string;
}

export interface CardTextboxNewlineNode extends CardContentNodeBase {
	readonly type: 'newline';
};

export interface CardPollNode extends CardContentNodeBase {
	readonly type: 'poll';
	is_quiz?: boolean;
	content: CardPollNodeOption[];
};

export interface CardPollNodeOption {
	value: string;
	is_answer?: boolean;
};

export interface CardImageNode {
	readonly type: 'image';
	media_id?: string | null;
	state?: UploadReadyState | null;
};

export enum UploadReadyState {
	Idle,
	Uploading,
	Done
};
