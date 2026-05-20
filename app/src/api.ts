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

const unwrapError = (error: any): Error => {
	return error instanceof Error ? error : new Error(`${error}`);
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

	private exec = async <T, R extends {} = {}> (method: string, proc: string, params?: MethodParamsInit, body?: R | null): Promise<Result<T>> => {

		const headers = new Headers({ 'Accept': 'application/json' });

		if (body) {
			headers.set('Content-Type', 'application/json');
		}

		const { response, fetchError } = await fetch(this.procURL(proc, params), {
			method,
			headers,
			body: body ? JSON.stringify(body) : null,
		}).then(response => ({ response, fetchError: null }))
			.catch(err => ({ response: null, fetchError: unwrapError(err)}));

		if (fetchError) {
			return { data: null, error: new Error(`Fetch API: ${fetchError.message}`) };
		}

		if (response.status === 204) {
			return { data: null, error: null };
		}

		const { result, parseError } = await response.json()
			.then(result => ({ result: result as Result<T>, parseError: null }))
			.catch(err => ({ result: null, parseError: unwrapError(err)}));

		if (parseError) {
			return { data: null, error: new Error(`Parse API response: ${parseError.message}; status code: ${response.status}`) };
		}

		return result;
	};

	auth = {

		whoami: async (opts?: { cached?: boolean }): Promise<Result<AuthState>> => {

			if (opts?.cached && this.authCache.valid()) {
				return { data: this.authCache.state, error: null };
			}

			const result = await this.exec<AuthState>('GET', '/auth/whoami');
			this.authCache.store(result.data);
			return result;
		},

		signin: async (params: SignInParams): Promise<Result<AuthState>> => {
			const result = await this.exec<AuthState>('POST', '/auth/signin', {}, params)
			this.authCache.store(result.data);
			return result;
		},

		signout: async (): Promise<Result<AuthState>> => {
			const result = await this.exec<AuthState>('POST', '/auth/signout')
			this.authCache.store(result.data);
			return result;
		},
	};

	collections = {

		list: async (params?: { ids?: string[] | null } & Partial<Pagination>) =>
			this.exec<Page<CollectionMetadata>>('GET', '/collections', params),

		search: async (term: string) =>
			this.exec<Page<CollectionSearchResult>>('GET', '/collections/search', { term }),

		load: async (id: string) =>
			this.exec<Collection>('GET', `/collections/${id}`),

		create: async (patch: CollectionPatch) =>
			this.exec<CollectionMetadata>('PUT', '/manage/content/collection', {}, patch),

		update: async (id: string, patch: CollectionPatch) =>
			this.exec<CollectionMetadata>('PATCH', `/manage/content/collection/${id}/metadata`, {}, patch),

		remove: async (id: string) =>
			this.exec<null>('DELETE', `/manage/content/collection/${id}`),
	};

	decks = {

		list: async (params?: { ids?: string[] | null, collection_id?: string } & Partial<Pagination>) =>
			this.exec<Page<CardDeckMetadata>>('GET','/decks', params),

		load: async (id: string) =>
			this.exec<CardDeck>('GET',`/decks/${id}`),

		create: async (patch: CardDeckPatch) =>
			this.exec<CardDeckMetadata>('PUT', '/manage/content/deck', {}, patch),

		update: async (id: string, patch: CardDeckPatch) =>
			this.exec<CardDeckMetadata>('PATCH', `/manage/content/deck/${id}`, {}, patch),

		remove: async (id: string) =>
			this.exec<null>('DELETE', `/manage/content/deck/${id}`),
	};
};

export const useClient = () => new ApiClient('/api');
