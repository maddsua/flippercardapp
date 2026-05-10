import type { CardDeck, CardDeckMetadata, Collection, CollectionMetadata } from "./api_models";

export interface Result <T> {
	data: T | null;
	error: Error | null;
};

export interface Page <T> {
	entries: T[];
	offset: number;
	limit: number;
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

export class ApiClient {

	readonly endpoint: URL;

	constructor(endpoint: string | URL) {
		if (typeof endpoint !== 'string' || endpoint.includes('://')) {
			this.endpoint = new URL(endpoint);
		} else {
			this.endpoint = new URL(window.location.href);
			this.endpoint.pathname = endpoint;
		}
	}

	private procedureURL = (name: string, params?: MethodParamsInit) => {

		const url = new URL(this.endpoint.href);

		const prefixNormalized = url.pathname.endsWith('/') ? url.pathname.slice(0, -1) : url.pathname;
		const methodNormalized = name.startsWith('/') ? name.slice(1) : name

		url.pathname = `${prefixNormalized}/${methodNormalized}`;

		serializeQueryParams(url.searchParams, params);

		return url;
	};

	private exec = async <T> (method: string, proc: string, params?: MethodParamsInit) => {

		const { response, fetchError } = await fetch(this.procedureURL(proc, params), {
			method: method,
			headers: {
				'Accept': 'application/json'
			}
		}).then(response => ({ response, fetchError: null }))
			.catch(err => ({ response: null, fetchError: unwrapError(err)}));

		if (fetchError) {
			return { data: null, error: new Error(`Fetch API: ${fetchError.message}`) };
		}

		const { result, parseError } = await response.json()
			.then(result => ({ result: result as Result<T>, parseError: null }))
			.catch(err => ({ result: null, parseError: unwrapError(err)}));

		if (parseError) {
			return { data: null, error: new Error(`Parse API response: ${parseError.message}; status code: ${response.status}`) };
		}

		return result;
	};

	listCollections = async (params?: { ids?: string[] | null } & Partial<Pagination>) => {
		return this.exec<Page<CollectionMetadata>>('GET', '/collections', params);
	};

	loadCollection = async (id: string) => {
		return this.exec<Collection>('GET', `/collections/${id}`);
	};

	listDecks = async (params?: { ids?: string[] | null, collection_id?: string } & Partial<Pagination>) => {
		return this.exec<Page<CardDeckMetadata>>('GET','/decks', params);
	};

	loadDeck = async (id: string) => {
		return this.exec<CardDeck>('GET',`/decks/${id}`);
	};
};

export const useClient = () => new ApiClient('/api');
