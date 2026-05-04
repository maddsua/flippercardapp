// this is a sample data providerd. it is to be deleted later on during development

import type { CardCollection, CardDeck, CardNode, CollectionProvider } from "../../components/Cards/content";

import mixed_deck from './decks/mixed_deck';

const wrapResult = <T>(data: T) => ({ data, error: null });

class SampleCollectionProvider implements CollectionProvider {

	private readonly _collections: SampleCollection[];

	constructor(collections: SampleCollection[]) {
		this._collections = collections;
	}

	collections = async (id?: string) => wrapResult(this._collections.filter(item => item.id === id || !id));
};

class SampleCollection implements CardCollection {

	readonly id: string;
	readonly name: string;
	readonly size: number;

	private readonly _decks: SampleDeck[];

	constructor(name: string, decks: SampleDeck[]) {
		this.id = crypto.randomUUID();
		this.name = name;
		this.size = decks.length;
		this._decks = decks.map(item => item.linkWithCollection(this));
	}

	decks = async (id?: string) => wrapResult(this._decks.filter(item => item.id === id || !id));
};

class SampleDeck implements CardDeck {

	readonly id: string;
	readonly name: string;
	readonly size: number;

	private _collection: SampleCollection | null;
	private readonly _cards: CardNode[];

	constructor(name: string, cards: CardNode[]) {
		this.id = crypto.randomUUID();
		this.name = name;
		this.size = cards.length;
		this._cards = cards;
		this._collection = null;
	}

	linkWithCollection = (collection: SampleCollection) => {
		this._collection = collection
		return this;
	};

	collection = async () => {
		if (!this._collection) {
			throw new Error('Not initialized');
		}

		return { data: this._collection, error: null };
	};

	cards = async (id?: string) => wrapResult(this._cards.filter(item => item.id === id || !id));
};

export const sampleProvider = new SampleCollectionProvider([
	new SampleCollection('Sample cards', [
		new SampleDeck('A quiz about everything', mixed_deck)
	]),
]);
