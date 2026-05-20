
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

class ListStore<T> {

	private readonly store: GenericStore<T[]>;

	constructor(key: string) {
		this.store = new GenericStore(key)
	}

	entries = async (): Promise<T[]> => {
		return (await this.store.load()) || []
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

class Storage {

	constructor() {
		this.collections = new ListStore<string>('saved_collections');
		this.starredDecks = new ListStore<string>('starred_decks');
		this.deckEditor = new GenericStore<any>('deck_editor_state_snapshot');
	}

	collections: ListStore<string>;
	starredDecks: ListStore<string>;
	deckEditor: GenericStore<object>;
}

export const useStorage = () => new Storage();
