
interface BaseNode {
	type: string;
};

export interface TitleNode extends BaseNode {
	type: 'title';
	content: string;
};

export interface TextBoxNode extends BaseNode {
	type: 'textbox';
	content: Array<TextNote | NewlineNode>;
};

export interface TextNote extends BaseNode {
	type: 'text';
	content: string;
	//	todo: add styling
};

export interface NewlineNode extends BaseNode {
	type: 'newline';
};

export interface PollNode extends BaseNode {
	type: 'poll';
	is_quiz?: boolean;
	content: PollOption[];
};

export interface PollOption {
	value: string;
	is_answer?: boolean;
}

//	todo: add images

export type ContentNode = TitleNode | TextBoxNode | PollNode;

export interface CardSideNode {
	content: ContentNode[];
	theme?: CardNodeTheme;
};

export interface ElementTheme {
	fill_color?: string;
	mask_color?: string;
}

export interface CardFaceTheme extends ElementTheme {
	outline_color?: string;
};

export interface CardNodeTheme {
	card: CardFaceTheme;
	interactives?: ElementTheme;
}

export interface CardNode {
	id: string;
	front: CardSideNode;
	back: CardSideNode;
};
