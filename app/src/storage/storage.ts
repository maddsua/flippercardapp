import type { CollectionPlayStats, DeckPlayStats } from "@/play";
import { useIDB } from "./idb";
import { GenericKVStore, GenericKVStoreWithDefault } from "./kv";



export const useStorage = () => {
	return {

		collections: {

			stats: {
				aggregated: async (collection_ids?: string | string[]) => {

					const queryFiltered = async () => {
						const store = await useIDB().then(db => db.deckPlayStats);
						switch (typeof collection_ids) {
							case 'string':
								return store.filter(item => item.deck_id === collection_ids);
							case 'object':
								const idSet = new Set(collection_ids);
								return store.filter(item => !!item.collection_id && idSet.has(item.collection_id));
							default:
								return store.all();
						}
					};

					const aggregate = new Map<string, CollectionPlayStats>();

					for (const entry of await queryFiltered()) {

						if (!entry.collection_id) {
							continue;
						}

						const acc = aggregate.get(entry.collection_id);

						aggregate.set(entry.collection_id, {
							collection_id: entry.collection_id,
							avg_score: acc ? (acc.avg_score + entry.score) / 2 : entry.score,
							decks_played: (acc?.decks_played ?? 0) + 1
						});
					}

					return Array.from(aggregate.entries());
				},
			},

			starred: {
				add: (id: string) => useIDB().then(store => store.starredCollections.add(id)),
				del: (id: string) => useIDB().then(store => store.starredCollections.del(id)),
				has: (id: string) => useIDB().then(store => store.starredCollections.has(id)),
				all: () => useIDB().then(store => store.starredCollections.entries()),
			},
		},

		decks: {

			stats: {
				store: (value: DeckPlayStats) => useIDB().then(db => db.deckPlayStats.store(value)),
				load: (id: string) => useIDB().then(db => db.deckPlayStats.load(id)),
				del: (key: string) => useIDB().then(db => db.deckPlayStats.del(key)),
				has: (key: string) => useIDB().then(db => db.deckPlayStats.has(key)),
				all: () => useIDB().then(db => db.deckPlayStats.all()),
				filter: async (predicate: (val: DeckPlayStats) => boolean) => useIDB().then(db => db.deckPlayStats.filter(predicate)),
			},

			starred: {
				add: (id: string) => useIDB().then(db => db.starredDecks.add(id)),
				del: (id: string) => useIDB().then(db => db.starredDecks.del(id)),
				has: (id: string) => useIDB().then(db => db.starredDecks.has(id)),
				all: () => useIDB().then(db => db.starredDecks.entries()),
			},

			editor: {
				snapshot: new GenericKVStore('deck_editor_state_snapshot'),
			},
		},

		preferences: {
			language: new GenericKVStore<string>('app_language'),
			playModeShowNavigation: new GenericKVStoreWithDefault('play_mode_show_navigation', true),
		},
	};
};
