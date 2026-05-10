import { useClient, type ApiClient } from "./api";
import type { CardDeckMetadata, CollectionMetadata } from "./api_models";
import type { CardCollection, CardDeck } from "./content";

export interface MethodResult <T> {
	data: T | null;
	error: Error | null;
};

export const useContent = (): ContentProvider => {
	return new ContentProvider(useClient());
};

class ContentProvider {

	private client: ApiClient;

	constructor(client: ApiClient) {
		this.client = client;
	}

	collections = async (id?: string): Promise<MethodResult<CardCollection[]>> => {

		const { data, error } = await this.client.listCollections({ ids: id ? [id] : null });
		if (!data || error) {
			return { data: null, error: error };
		}

		return { data: data.entries.map(item => new CollectionProvider(item, this.client)), error: null };
	};

	decks = async (id?: string) => {

		const { data, error } = await this.client.listDecks({ ids: id ? [id] : null });
		if (!data || error) {
			return { data: null, error: error };
		}

		return { data: data.entries.map(item => new DeckProvider(item, this.client)), error: null };
	};

	starred = async (): Promise<MethodResult<CardDeck[]>> => {
		//	todo: implement
		throw new Error('not implemented yet');
	};
};

class CollectionProvider implements CardCollection {

	private client: ApiClient;

	readonly id: string;
	readonly name: string;
	readonly size: number;

	constructor(collection: CollectionMetadata, client: ApiClient) {
		this.client = client;
		this.id = collection.id;
		this.name = collection.name;
		this.size = collection.size;
	}

	decks = async () => {

		const { data, error } = await this.client.listDecks({ collection_id: this.id });
		if (!data || error) {
			return { data: null, error: error };
		}

		return { data: data.entries.map(item => new DeckProvider(item, this.client)), error: null };
	};
};

class DeckProvider implements CardDeck {

	private client: ApiClient;
	private collectionID: string;

	readonly id: string;
	readonly name: string;
	readonly size: number;

	constructor(deck: CardDeckMetadata, client: ApiClient) {
		this.client = client;
		this.collectionID = deck.collection_id;
		this.id = deck.id;
		this.name = deck.name;
		this.size = deck.size;
	}

	collection = async () => {

		const { data, error } = await this.client.listCollections({ ids: [this.collectionID] });
		if (!data || error) {
			return { data: null, error: error };
		}

		const collection = data.entries.at(0);
		if (!collection) {
			return { data: null, error: new Error('collection not found') };
		}

		return { data: new CollectionProvider(collection, this.client), error: null };
	};

	cards = async () => {

		const { data, error } = await this.client.loadDeck(this.id);
		if (!data || error) {
			return { data: null, error: error };
		}

		return { data: data.cards.map(item => ({ id: item.id, ... item.content })), error: null };
	};
};
