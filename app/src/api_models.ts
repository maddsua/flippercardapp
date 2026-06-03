import type { CardNode } from "./content";

export interface ContentEntryMetaBase {
	name: string;
	description?: string | null;
	visibility: ResourceVisibility;
};

export type ResourceVisibility = 'PRIVATE' | 'HIDDEN' | 'PUBLIC';

export interface CollectionMetadata extends ContentEntryMetaBase{
	id: string;
	created: string;
	updated: string;
	size: number;
};

export interface CollectionSearchResult extends CollectionMetadata {
	rank: number;
};

export interface CardDeckMetadata extends ContentEntryMetaBase{
	id: string;
	collection_id: string;
	created: string;
	updated: string;
	size: number;
};

export interface CardDeck extends CardDeckMetadata {
	labels: string[];
	cards: CardNode[];
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
	team_member: boolean;
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
	meta?: ContentEntryMetaBase | null;
	content?: CardDeckContentPatch | null;
};

export interface CardDeckContentPatch  {
	cards: CardPatch[];
};

export interface CardPatch extends Omit<CardNode, 'id'> {
	id?: string | null;
};

export interface ImageMetadata {
	id: string;
	created: string;
	mimetype: string;
	source_name: string;
	source_sha512_hash: string;
	data_sha512_hash: string;
	data_size: number;
};

export interface CardDeckVersionMetadata {
	id: string;
	created: string;
	deck_id: string;
	card_count: number;
	is_latest: boolean;
	label?: string | null;
};
