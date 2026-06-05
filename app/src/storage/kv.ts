
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
			console.error('GenericKVStore.load:' ,error);
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
			console.error('GenericKVStore.store:' ,error);
		}
	};

	clear = () => {
		try {
			localStorage.removeItem(this.key);
		} catch (error) {
			console.error('GenericKVStore.clear:' ,error);
		}
	};
};

export class GenericKVStoreWithDefault <T> {

	private readonly defaultValue: T;

	//	doing this because fuck the way java/type-script handles classes
	private base: GenericKVStore<T>;

	constructor(key: string, defaultValue: T) {
		this.base = new GenericKVStore<T>(key);
		this.defaultValue = defaultValue;
	}

	load = (): T => this.base.load() ?? this.defaultValue;

	store = (val: T | null) => this.base.store(val);

	clear = () => this.base.clear();
};
