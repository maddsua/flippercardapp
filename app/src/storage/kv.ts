
export class GenericKVStore<T> {

	private readonly key: string;

	constructor(key: string) {
		this.key = key;
	}

	load = (): T | null => {
		const val = localStorage.getItem(this.key);
		try { return val ? JSON.parse(val) : null; }
			catch (_) { return null }
	};

	store = (val: T | null) => {
		if (val === null) {
			localStorage.removeItem(this.key);
			return;
		}
		localStorage.setItem(this.key, JSON.stringify(val));
	};

	clear = () => {
		localStorage.removeItem(this.key);
	};
};
