
class GenericStore<T> {

	private readonly key: string;

	constructor(key: string) {
		this.key = key;
	}

	load = async (): Promise<T | null> => {
		const val = localStorage.getItem(this.key);
		return val ? JSON.parse(val) : null;
	};

	store = async (val: T | null) => {
		if (val === null) {
			localStorage.removeItem(this.key);
			return;
		}
		localStorage.setItem(this.key, JSON.stringify(val));
	};
};

class SetStore<T> {

	private readonly store: GenericStore<T[]>;

	constructor(key: string) {
		this.store = new GenericStore(key)
	}

	entries = async (): Promise<T[]> => {
		return (await this.store.load()) || [];
	};

	fromEntries = async (entries: T[]) => {
		await this.store.store(entries.length ? entries : null);
	};

	contains = async (value: T): Promise<boolean> => {
		return await this.entries().then(entries => entries.some((item) => item === value));
	};

	add = async (val: T) => {

		const entries = await this.entries();
		if (new Set(entries).has(val)) {
			return false;
		}

		entries.push(val);

		await this.store.store(entries);

		return true;
	};

	remove = async (val: T) => {

		const entries = await this.entries();
		const newEntries = entries.filter(item => item !== val);

		if (entries.length === newEntries.length) {
			return false;
		}

		await this.store.store(newEntries);

		return true;
	};
};

class MapStore <T> {

	private readonly _store: GenericStore<Record<string, T>>;

	constructor(key: string) {
		this._store = new GenericStore(key)
	}

	entries = async () => Object.entries(await this._store.load() || {});

	store = async (key: string, val: T | null) => {

		const fields = await this._store.load() || {};

		if (val === null) {
			delete fields[key];
		} else {
			fields[key]=val;
		}

		await this._store.store(fields);
	};

	load = async (key: string): Promise<T | null> => {

		const fields = await this._store.load();
		if (!fields) {
			return null;
		}

		return fields[key] || null;
	};

	remove = async (key: string) => {

		const fields = await this._store.load();
		if (!fields || !fields[key]) {
			return;
		}

		delete fields[key];

		await this._store.store(fields);
	};

};

class PlayStatsStore extends MapStore<PlayStats> {

	constructor() {
		super('play_stats')
	}

	collectionScores = async (): Promise<Map<string, number>> => {

		const scoreMap = new Map<string, number>();

		for (const [_, entry] of await this.entries()) {

			if (!entry.collection_id) {
				continue;
			}

			const acc = scoreMap.get(entry.collection_id);
			const avg = typeof acc === 'number' ? (acc + entry.score) / 2 : entry.score;

			scoreMap.set(entry.collection_id, avg);
		}

		return scoreMap;
	};
};

class Storage {

	constructor() {
		this.collections = new SetStore('saved_collections');
		this.starredDecks = new SetStore('starred_decks');
		this.deckEditor = new GenericStore('deck_editor_state_snapshot');
		this.playStats = new PlayStatsStore();
	}

	collections: SetStore<string>;
	starredDecks: SetStore<string>;
	deckEditor: GenericStore<object>;
	playStats: PlayStatsStore;
}

export interface PlayStats {
	deck_id: string;
	collection_id: string | null;
	score: number;
};

export const useStorage = () => new Storage();
