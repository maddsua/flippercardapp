import type { DeckPlayStats } from "@/play";
import { GenericKVStore } from "./kv";
import type { DeckEditorHistoryMetaEntry } from "./storage";

declare let window: Window & {
	appUserDB?: IDBDatabase;
};

interface dbOpenOptions {
	onUpgrade?: (db: IDBDatabase, event: IDBVersionChangeEvent) => void;
};

const openDB = (name: string, version?: number, opts?: dbOpenOptions): Promise<IDBDatabase> => {

	const req = window.indexedDB.open(name, version);

	const { onUpgrade } = opts || {};

	if (onUpgrade) {
		req.onupgradeneeded = (event) => onUpgrade(req.result, event);
	}

	req.onblocked = (event) => {
		console.warn(`IDB upgrade has been requested for '${name}'; version ${event.oldVersion}->${event.newVersion}`);
		window.location.reload();
	};

	return unwrapRequest(req);
};

const initDB = (name: string, migrations: Migration[]): Promise<IDBDatabase> => {

	try {

		migrations = migrations.sort((a, b) => a.version - b.version);
		const latest = migrations[migrations.length - 1];
		if (!latest) {
			throw new Error('Unable to find latest IDB migration');
		}

		return openDB(name, latest.version, {
			onUpgrade: (db, event) => {

				const currentMigration = migrations.findIndex(item => item.version === event.oldVersion);
				const steps = migrations.slice(currentMigration + 1);

				console.debug('IDB executing migrations:', steps.map(item => item.version));

				for (const step of steps) {
					step.onUpgrade(db);
				}

				console.debug(`IDB version upgraded: version ${event.oldVersion}->${event.newVersion}`);
			}
		});

	} catch (error) {
		console.error('IDB INIT:', name, error instanceof Error ? error.message : `${error}`);
		throw error;
	}
};

const unwrapRequest = <T>(req: IDBRequest<T>) => {
	return new Promise<T>((resolve, reject) => {
		req.onsuccess = () => resolve(req.result);
		req.onerror = () => {
			console.error('IDB request rejected:', req.error?.message);
			reject(req.error);
		};
	});
};

const getStore = (db: IDBDatabase, storeName: string, mode: IDBTransactionMode): IDBObjectStore => {
	return db.transaction(storeName, mode).objectStore(storeName)
};

interface IDSetStoreValue <T> {
	id: T;
};

type Maybe <T> = T | null | undefined;

class IDSetStore <T extends string> {

	private readonly db: IDBDatabase;
	private readonly storeName: string;

	constructor(db: IDBDatabase, storeName: string) {
		this.db = db;
		this.storeName = storeName;
	}

	static create = (db: IDBDatabase, name: string, entriesInit?: string[]) => {
		const store = db.createObjectStore(name, { keyPath: 'id' satisfies keyof IDSetStoreValue<string> });
		store.createIndex('set_idx', 'id', { unique: true });
		entriesInit?.forEach(entry => store.put({ id: entry } satisfies IDSetStoreValue<string>));
	};

	private tx = (mode: IDBTransactionMode) =>
		getStore(this.db, this.storeName, mode);

	add = async (value: T) =>
		unwrapRequest(this.tx('readwrite').put({ id: value } satisfies IDSetStoreValue<T>));

	del = async (value: T) =>
		unwrapRequest(this.tx('readwrite').delete(value));

	has = async (value: T) =>
		unwrapRequest<Maybe<IDSetStoreValue<T>>>(this.tx('readonly').get(value))
			.then(entry => !!entry?.id);

	entries = async () =>
		unwrapRequest<Array<IDSetStoreValue<T>>>(this.tx('readonly').getAll())
			.then(entries => entries.map(item => item.id));
};

class UniqueCollectionStore <T extends {}> {

	private readonly db: IDBDatabase;
	private readonly storeName: string;

	constructor(db: IDBDatabase, storeName: string) {
		this.db = db;
		this.storeName = storeName;
	}

	static create = <T extends {}>(db: IDBDatabase, name: string, keyPathInit: keyof T, entriesInit?: T[]) => {
		const store = db.createObjectStore(name, { keyPath: String(keyPathInit) });
		store.createIndex(`${String(keyPathInit)}_idx`, 'id', { unique: true });
		entriesInit?.forEach(entry => store.put(entry));
	};

	private tx = (mode: IDBTransactionMode) =>
		getStore(this.db, this.storeName, mode);

	store = async (value: T) =>
		unwrapRequest(this.tx('readwrite').put(value));

	del = async (key: string) =>
		unwrapRequest(this.tx('readwrite').delete(key));

	load = async (key: string) =>
		unwrapRequest<Maybe<T>>(this.tx('readonly').get(key))
			.then(entry => entry || null);

	has = async (key: string) =>
		unwrapRequest(this.tx('readonly').get(key))
			.then(entry => !!entry);

	entries = async () =>
		unwrapRequest<Array<T>>(this.tx('readonly').getAll());

	filter = async (predicate: (val: T) => boolean) => {
		return new Promise<T[]>((resolve, reject) => {

			const cursor = this.tx('readonly').openCursor();

			const entries: T[] = [];

			cursor.onsuccess = () => {

				const { result: next } = cursor;
				if (!next) {
					resolve(entries);
					return;
				}

				const { value } = next;
				if (typeof value === 'object' && value && predicate(value)) {
					entries.push(structuredClone(value));
				}

				next.continue();
			};

			cursor.onerror = () => {
				console.error('IDB cursor rejected:', cursor.error?.message);
				reject(cursor.error);
			};
		});
	};
};

interface VersionedStoreEntryBase {
	timestamp: Date;
};

class VersionedCollectionStore <E extends VersionedStoreEntryBase> {

	private readonly db: IDBDatabase;
	private readonly storeName: string;

	constructor(db: IDBDatabase, storeName: string) {
		this.db = db;
		this.storeName = storeName;
	}

	static create = <T extends {}>(db: IDBDatabase, name: string, primaryKey: Exclude<keyof T, 'timestamp'>, entriesInit?: T[]) => {

		const keyPath = [String(primaryKey), 'timestamp' satisfies keyof VersionedStoreEntryBase];

		const store = db.createObjectStore(name, { keyPath: keyPath });
		keyPath.forEach(key => store.createIndex(`${key}_idx`, key, { unique: false }))

		entriesInit?.forEach(entry => store.put(entry));
	};

	private tx = (mode: IDBTransactionMode) =>
		getStore(this.db, this.storeName, mode);

	private keyRange = (key: string) => IDBKeyRange.bound([key, new Date(0)], [key, new Date()]);

	push = async <V extends E>(val: V) =>
		unwrapRequest(this.tx('readwrite').put(val));

	latest = async <V extends E>(key: string) => {
		return new Promise<V | null>((resolve, reject) => {

			const cursor = this.tx('readonly').openCursor(this.keyRange(key), 'prev');

			cursor.onsuccess = () => {

				const { result: next } = cursor;
				if (!next) {
					resolve(null);
					return;
				}

				resolve(structuredClone(next.value));
			};

			cursor.onerror = () => {
				console.error('IDB cursor rejected:', cursor.error?.message);
				reject(cursor.error);
			};
		});
	};

	versions = async (key: string) =>
		unwrapRequest<Array<E>>(this.tx('readonly').getAll(this.keyRange(key)));

	remove = async (key: string) =>
		unwrapRequest(this.tx('readwrite').delete(this.keyRange(key)));
};

interface Migration {
	version: number;
	onUpgrade: (db: IDBDatabase) => void;
};

const migrationList: Migration[] = [
	{
		version: 1,
		onUpgrade: (db) => {

			const kvIdList = (key: string): string[] => {
				const values: string[] = new GenericKVStore<string[]>(key).load() || [];
				if (!Array.isArray(values)) {
					return [];
				}
				return values.filter(item => typeof item === 'string');
			};

			IDSetStore.create(db, 'starred_collections', kvIdList('saved_collections'));
			IDSetStore.create(db, 'starred_decks', kvIdList('starred_decks'));

			const statEntries = Object.values(new GenericKVStore('play_stats').load() || {})
				.filter(item => typeof item === 'object' && !!item && 'deck_id' in item);

			UniqueCollectionStore.create<DeckPlayStats>(db, 'deck_play_stats', 'deck_id', statEntries);
		},
	},
	{
		version: 2,
		onUpgrade: (db) => {
			VersionedCollectionStore.create<DeckEditorHistoryMetaEntry>(db, 'deck_editor_history', 'deck_id');
		},
	}
];

export const useIDB = async () => {

	if (!window.appUserDB) {
		window.appUserDB = await initDB('usercontent', migrationList);
	}

	return {
		starredDecks: new IDSetStore(window.appUserDB, 'starred_decks'),
		starredCollections: new IDSetStore(window.appUserDB, 'starred_collections'),
		deckPlayStats: new UniqueCollectionStore<DeckPlayStats>(window.appUserDB, 'deck_play_stats'),
		deckEditorHistory: new VersionedCollectionStore<DeckEditorHistoryMetaEntry>(window.appUserDB, 'deck_editor_history'),
	};
};
