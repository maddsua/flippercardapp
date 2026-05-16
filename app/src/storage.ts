
const loadTyped = <T>(key: string): T | null => {
	const val = localStorage.getItem(key);
	return val ? JSON.parse(val) : null;
};

const storeTyped = <T>(key: string, val: T | null) => {
	if (val === null) {
		localStorage.removeItem(key);
		return;
	}
	localStorage.setItem(key, JSON.stringify(val));
};

class Storage {

	//	todo: the rest of the nested methods, when moving to indexeddb

	//	todo: return as a set when needed
	collections = async () => {
		return loadTyped<string[]>('collections') || [];
	};

	addCollection = async (id: string) => {

		const entries = await this.collections();

		if (new Set(entries).has(id)) {
			return false;
		}

		entries.push(id);

		storeTyped('collections', entries);

		return true;
	};

	removeCollection = async (id: string) => {

		const entries = await this.collections();
		const newEntries = entries.filter(item => item !== id);

		if (entries.length === newEntries.length) {
			return false;
		}

		storeTyped('collections', newEntries);
		return true;
	};

	starred = async () => {
		return loadTyped<string[]>('starred') || [];
	};

	addStar = async (deckID: string) => {

		const entries = await this.starred();

		if (new Set(entries).has(deckID)) {
			return false;
		}

		entries.push(deckID);

		storeTyped('starred', entries);

		return true;
	};

	removeStar = async (deckID: string) => {

		const entries = await this.starred();
		const newEntries = entries.filter(item => item !== deckID);

		if (entries.length === newEntries.length) {
			return false;
		}

		storeTyped('starred', newEntries);
		return true;
	};

	deckEditor = {
		loadSnapshot: async () => loadTyped<any>('deck_editor_state_snapshot'),
		storeSnapshot: async (val: any | null) => storeTyped('deck_editor_state_snapshot', val),
	};
}

export const useStorage = () => new Storage();
