import type { CardCollection, CardDeck } from "./content";
import { sampleProvider } from "./data/sample";

export interface MethodResult <T> {
	data: T | null;
	error: Error | null;
};

export interface CollectionProvider {
	collections: (id?: string) => Promise<MethodResult<CardCollection[]>>;
	decks: (id?: string) => Promise<MethodResult<CardDeck[]>>;
};

//	todo: hook up API and Cache providers instead, when ready
export const useCollectionProvider = (): CollectionProvider => {
	return sampleProvider;
};
