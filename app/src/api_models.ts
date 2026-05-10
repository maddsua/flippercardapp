import type { CardNode } from "./content";

export interface CollectionMetadata {
	id: string;
	name: string;
	description?: string;
	created: string;
	updated: string;
	size: number;
};

export interface CardDeckMetadata {
	id: string;
	collection_id: string;
	name: string;
	description?: string;
	created: string;
	updated: string;
	size: number;
};

export interface CardDeck extends CardDeckMetadata {
	labels: string[];
	cards: Card[];
};

export interface Card {
	id: string;
	created: string;
	updated: string;
	content: Omit<CardNode, 'id'>;
};
