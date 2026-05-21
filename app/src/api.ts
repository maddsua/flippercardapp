import type {
	AuthState,
	CardDeck,
	CardDeckMetadata,
	CardDeckPatch,
	Collection,
	CollectionMetadata,
	CollectionPatch,
	CollectionSearchResult,
	SignInParams,
} from "./api_models";

export interface Result <T> {
	data: T | null;
	error: Error | null;
};

export interface BlobResult {
	blob: Blob | null;
	error: Error | null;
};

export interface Page <T> {
	entries: T[];
	offset: number;
	limit: number;
	has_next: boolean;
};

export interface Pagination {
	offset?: number;
	limit?: number;
};

type OneOrMore<T> = T | T[];
type ParamValue = OneOrMore<string> | OneOrMore<number>;
type MethodParamsInit = Record<string, ParamValue | null | undefined>;

type Method = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';

const unwrapError = (error: any): Error => {
	return error instanceof Error ? error : new Error(unwrapErrorMessage(error));
};

const unwrapErrorMessage = (error: any): string => {
	return typeof error === 'object' ? 'message' in error ? error.message : 'Unknown error' : `${error}`;
};

const serializeQueryParams = (target: URLSearchParams, params?: MethodParamsInit) => {

	if (!params) {
		return;
	}

	const entries = Object.entries(params);

	for (const [name, value] of entries) {

		if (!value) {
			continue
		}

		switch (typeof value) {
			case 'string':
				target.append(name, value);
				break;
			case 'number':
				target.append(name, value.toString());
				break;
			case 'object':
				serializeArrayQueryParams(target,name, value);
				break
			default: break;
		}
	}
};

const serializeArrayQueryParams = (target: URLSearchParams, name: string, value: ParamValue) => {
	if (!Array.isArray(value) || !value.length) {
		return
	}
	target.append(name, value.join(','));
};

class CachedAuthState {

	state: AuthState = {};

	private epoch: number = 0;

	readonly ttl = 900_000; // 900s or 15m

	store = (val: AuthState | null) => {
		this.state = val || {};
		this.epoch = new Date().getTime();
	};

	valid = () => (new Date().getTime() - this.epoch) < this.ttl;
};

export class ApiClient {

	readonly endpoint: URL;

	private authCache: CachedAuthState;

	constructor(endpoint: string | URL) {

		if (typeof endpoint !== 'string' || endpoint.includes('://')) {
			this.endpoint = new URL(endpoint);
		} else {
			this.endpoint = new URL(window.location.href);
			this.endpoint.pathname = endpoint;
		}

		this.authCache = new CachedAuthState();
	}

	private procURL = (name: string, params?: MethodParamsInit) => {

		const url = new URL(this.endpoint.href);

		const prefixNormalized = url.pathname.endsWith('/') ? url.pathname.slice(0, -1) : url.pathname;
		const methodNormalized = name.startsWith('/') ? name.slice(1) : name

		url.pathname = `${prefixNormalized}/${methodNormalized}`;

		serializeQueryParams(url.searchParams, params);

		return url;
	};

	private wrapPayload = <T extends {} = {}>(payload?: T | null) => {

		const headers = new Headers();

		if (payload instanceof Blob) {
			headers.set('Content-Type', payload.type || 'application/octet-stream');
			return { headers, body: payload, };
		} else if (payload instanceof File) {
			headers.set('Content-Type', payload.type || 'application/octet-stream');
			return { headers, body: payload, };
		} else if (payload) {
			headers.set('Content-Type', 'application/json');
			return { headers, body: JSON.stringify(payload), };
		}

		return { headers, body: null };
	};

	private fetch = async (url: URL | string, method: Method, headers: Headers, body: BodyInit | null) => {

		const { response, error } = await fetch(url, { method, headers, body, })
			.then(response => ({ response, error: null }))
			.catch(err => ({ response: null, error: unwrapError(err) }));

		if (error) {
			return { response: null, error: new Error(`Fetch API: ${error.message}`) };
		}

		return { response, error: null };
	};

	private unwrapJSON = async <T extends any = object> (response: Response): Promise<Result<T>> => {

		const { result, error: parseError } = await response.json()
			.then(result => ({ result: result as Result<T>, error: null }))
			.catch(err => ({ result: null, error: unwrapError(err)}));

		if (parseError) {
			return { data: null, error: new Error(`Parse API response: ${parseError.message}; status code: ${response.status}`) };
		}

		const { data, error } = result;
		if (!data && !error) {
			return { data: null, error: new Error('Received an empty data response') };
		}

		return { data, error };
	};

	private execJSON = async <R extends any = object, P extends any = object> (
		method: Method,
		proc: string,
		params?: MethodParamsInit,
		payload?: P | null
	): Promise<Result<R>> => {

		const url = this.procURL(proc, params);
		const { headers, body } = this.wrapPayload(payload);

		headers.set('Accept', 'application/json');

		const { response, error: fetchError } = await this.fetch(url, method, headers, body);
		if (fetchError) {
			return { data: null, error: fetchError };
		}

		if (response.status === 204) {
			return { data: null, error: null };
		}

		return await this.unwrapJSON<R>(response);
	};

	private execBlob = async <P extends any = object> (
		method: Method,
		proc: string,
		params?: MethodParamsInit,
		payload?: P | null
	): Promise<BlobResult> => {

		const url = this.procURL(proc, params);
		const { headers, body } = this.wrapPayload(payload);

		headers.set('Accept', '*');

		const { response, error: fetchError } = await this.fetch(url, method, headers, body);
		if (fetchError) {
			return { blob: null, error: fetchError };
		}

		if (response.status === 204) {
			return { blob: null, error: null };
		} else if (!response.ok) {
			const { error } = await this.unwrapJSON<any>(response);
			return { blob: null, error };
		}

		const { blob, error: blobError } = await response.blob()
			.then(blob => ({ blob, error: null }))
			.catch(err => ({ blob: null, error: unwrapError(err) }));

		if (blobError) {
			return { blob: null, error: new Error(`Retreive blob: ${blobError}`) };
		}

		return { blob, error: null };
	};

	auth = {

		whoami: async (opts?: { cached?: boolean }): Promise<Result<AuthState>> => {

			if (opts?.cached && this.authCache.valid()) {
				return { data: this.authCache.state, error: null };
			}

			const result = await this.execJSON<AuthState>('GET', '/auth/whoami');
			this.authCache.store(result.data);
			return result;
		},

		signin: async (params: SignInParams): Promise<Result<AuthState>> => {
			const result = await this.execJSON<AuthState>('POST', '/auth/signin', {}, params);
			this.authCache.store(result.data);
			return result;
		},

		signout: async (): Promise<Result<AuthState>> => {
			const result = await this.execJSON<AuthState>('POST', '/auth/signout');
			this.authCache.store(result.data);
			return result;
		},
	};

	collections = {

		list: async (params?: { ids?: string[] | null } & Partial<Pagination>) =>
			this.execJSON<Page<CollectionMetadata>>('GET', '/collections', params),

		search: async (term: string) =>
			this.execJSON<Page<CollectionSearchResult>>('GET', '/collections/search', { term }),

		load: async (id: string) =>
			this.execJSON<Collection>('GET', `/collections/${id}`),

		create: async (patch: CollectionPatch) =>
			this.execJSON<CollectionMetadata>('PUT', '/manage/content/collection', {}, patch),

		update: async (id: string, patch: CollectionPatch) =>
			this.execJSON<CollectionMetadata>('PATCH', `/manage/content/collection/${id}/metadata`, {}, patch),

		remove: async (id: string) =>
			this.execJSON<null>('DELETE', `/manage/content/collection/${id}`),

		exportBundle: async (id: string) =>
			this.execBlob('POST', `/manage/content/collection/${id}/export`),

		importBundle: async (file: File) =>
			this.execJSON<CollectionMetadata>('POST', `/manage/content/collections/import`, {}, file),
	};

	decks = {

		list: async (params?: { ids?: string[] | null, collection_id?: string } & Partial<Pagination>) =>
			this.execJSON<Page<CardDeckMetadata>>('GET','/decks', params),

		load: async (id: string) =>
			this.execJSON<CardDeck>('GET',`/decks/${id}`),

		create: async (patch: CardDeckPatch) =>
			this.execJSON<CardDeckMetadata>('PUT', '/manage/content/deck', {}, patch),

		update: async (id: string, patch: CardDeckPatch) =>
			this.execJSON<CardDeckMetadata>('PATCH', `/manage/content/deck/${id}`, {}, patch),

		remove: async (id: string) =>
			this.execJSON<null>('DELETE', `/manage/content/deck/${id}`),
	};
};

export const useClient = () => new ApiClient('/api');
