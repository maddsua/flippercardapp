// this is a sample data providerd. it is to be deleted later on during development

import type { MethodResult } from "../../api";
import type { CardCollection, CardDeck, CardNode, CollectionProvider } from "../../components/Cards/content";

import mixed_deck from './decks/mixed_deck';

const wrapResult = <T>(data: T): MethodResult<T> => ({ data, error: null });

class SampleCollectionProvider implements CollectionProvider {

	readonly entries: SampleCollection[];

	constructor(collections: SampleCollection[]) {
		this.entries = collections;
	}

	collections = async (id?: string) => wrapResult(this.entries.filter(item => item.id === id || !id));
	decks = async (id?: string) => wrapResult(this.entries.map(item => item.entries).flat().filter(item => item.id === id || !id));
};

class SampleCollection implements CardCollection {

	readonly id: string;
	readonly name: string;
	readonly size: number;

	readonly entries: SampleDeck[];

	constructor(id: string, name: string, decks: SampleDeck[]) {
		this.id = id;
		this.name = name;
		this.size = decks.length;
		this.entries = decks.map(item => item.linkWithCollection(this));
	}

	decks = async (id?: string) => wrapResult(this.entries.filter(item => item.id === id || !id));
};

class SampleDeck implements CardDeck {

	readonly id: string;
	readonly name: string;
	readonly size: number;

	readonly entries: CardNode[];

	private _parent: SampleCollection | null;

	constructor(id: string, name: string, cards: CardNode[]) {
		this.id = id;
		this.name = name;
		this.size = cards.length;
		this.entries = cards;
		this._parent = null;
	}

	linkWithCollection = (collection: SampleCollection) => {
		this._parent = collection
		return this;
	};

	collection = async () => {
		if (!this._parent) {
			throw new Error('Not initialized');
		}

		return { data: this._parent, error: null };
	};

	cards = async (id?: string) => wrapResult(this.entries.filter(item => item.id === id || !id));
};

export const sampleProvider = new SampleCollectionProvider([
	new SampleCollection( '630cfd08-924c-49fa-b5b2-2c81e979829d', 'Sample cards', [
		new SampleDeck('33511bb5-f9d7-4795-9b3a-c1479378c27b', 'A quiz about everything', mixed_deck)
	]),
]);
