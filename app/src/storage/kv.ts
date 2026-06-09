
export class GenericKVStore <T> {

	private readonly key: string;

	constructor(key: string) {
		this.key = key;
	}

	load = (): T | null => {
		const val = localStorage.getItem(this.key);
		try {
			return val ? JSON.parse(val) : null;
		} catch (error) {
			console.error(`GenericKVStore.load '${this.key}':`, error);
			return null;
		}
	};

	store = (val: T | null) => {
		try {
			if (val === null) {
				localStorage.removeItem(this.key);
			} else {
				localStorage.setItem(this.key, JSON.stringify(val));
			}
		} catch (error) {
			console.error(`GenericKVStore.store '${this.key}':`, error);
		}
	};

	clear = () => {
		try {
			localStorage.removeItem(this.key);
		} catch (error) {
			console.error(`GenericKVStore.clear '${this.key}':`, error);
		}
	};
};

export class KVFlagStore {

	private readonly key: string;
	private readonly defaultValue: boolean;

	constructor(key: string, defaultValue: boolean) {
		this.key = key;
		this.defaultValue = defaultValue;
	}

	load = (): boolean => {

		const value = localStorage.getItem(this.key);
		if (value === 'true') {
			return true;
		} else if (value === 'false') {
			return false;
		}

		return this.defaultValue;
	};

	store = (value: boolean) => {

		if (value === this.defaultValue) {
			this.clear();
			return;
		}

		localStorage.setItem(this.key, value ? 'true' : 'false');
	};

	clear = () => localStorage.removeItem(this.key);
};

export class KVStringStore {

	private readonly key: string;

	constructor(key: string) {
		this.key = key;
	}

	load = () => localStorage.getItem(this.key);
	store = (value: string) => localStorage.setItem(this.key, value);
	clear = () => localStorage.removeItem(this.key);
};
