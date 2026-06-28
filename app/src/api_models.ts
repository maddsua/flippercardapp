import type { CardNode, ContentSummary } from "./content";

interface ContentEntryMeta extends ContentSummary {
	name: string;
	description?: string | null;
	visibility: ResourceVisibility;
	id: string;
	created: string;
	updated: string;
};

export type ResourceVisibility = 'PRIVATE' | 'HIDDEN' | 'PUBLIC';

export interface CollectionMeta extends ContentEntryMeta {
	size: number;
};

export interface CollectionSearchResult extends CollectionMeta {
	rank: number;
};

export interface CardDeckMeta extends ContentEntryMeta {
	collection_id: string;
	version_id?: string | null;
	size: number;
};

export interface CardDeckContent {
	cards: CardNode[];
};

export interface CardDeck extends CardDeckMeta, CardDeckContent {
	labels: string[];
};

export interface Collection extends CollectionMeta {
	decks: CardDeckMeta[];
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

export interface CollectionPatch extends ContentSummary {
	visibility: ResourceVisibility;
};

export interface CardDeckPatch {
	collection_id?: string | null;
	label?: string | null;
	summary?: ContentSummary | null;
	visibility?: ResourceVisibility | null;
	content?: CardDeckContentPatch | null;
};

export interface CardDeckContentPatch {
	cards: CardPatch[];
};

export interface CardPatch extends Omit<CardNode, 'id'> {
	id?: string | null;
};

export interface ImageMeta {
	id: string;
	created: string;
	mimetype: string;
	source_name: string;
	source_sha512_hash: string;
	data_sha512_hash: string;
	data_size: number;
};

interface CardDeckVersionMetaBase {
	id: string;
	created: string;
	deck_id: string;
	card_count: number;
	is_latest: boolean;
	label?: string | null;
};

export interface CardDeckVersionMeta extends CardDeckVersionMetaBase {
	summary?: ContentSummary | null;
};

export interface CardDeckVersionContent extends CardDeckContent {
	summary: ContentSummary;
};

export interface CardDeckVersion extends CardDeckVersionMetaBase {
	content: CardDeckVersionContent;
};
