import type { CardContentNode } from "./content";

export interface CollectionMetadata {
	id: string;
	name: string;
	description?: string;
	created: string;
	updated: string;
	size: number;
};

export interface CollectionSearchResult extends CollectionMetadata {
	rank: number;
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

export interface Card extends CardContentNode {
	created: string;
	updated: string;
};

export interface Collection extends CollectionMetadata {
	decks: CardDeckMetadata[];
};

export interface AuthState {
	actor?: AuthActor | null;
	session?: AuthSession | null;
};

export interface AuthActor {
	id: string;
	name: string;
	permissions: UserPermissions;
};

export interface AuthSession {
	id: string;
	expires: string;
};

export interface UserPermissions {
	administrative: boolean;
	content_edit: boolean;
};

export interface SignInParams {
	username: string;
	password: string;
};

export interface CollectionPatch {
	name: string;
	description?: string;
};

export interface CardDeckPatch {
	collection_id?: string | null;
	details?: CardDeckDetailsPatch | null;
	content?: CardDeckContentPatch | null;
};

export interface CardDeckDetailsPatch {
	name: string;
	description?: string | null;
};

export interface CardDeckContentPatch  {
	cards: CardPatch[];
};

export interface CardPatch extends Omit<CardContentNode, 'id'> {
	id?: string | null;
};
