import type { CardDeck, CardDeckMetadata, CollectionMetadata } from "./api_models";

export interface APIResult <T> {
	data: T | null;
	error: Error | null;
};

export interface ApiPage <T> {
	entries: T[];
	offset: number;
	limit: number;
};

export interface Pagination {
	offset?: number;
	limit?: number;
};

type MethodParamsInit = URLSearchParams | Record<string, string | number | null | undefined>;

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

	private procedureURL = (method: string, params?: MethodParamsInit) => {

		const url = new URL(this.endpoint.href);

		const prefixNormalized = url.pathname.endsWith('/') ? url.pathname.slice(0, -1) : url.pathname;
		const methodNormalized = method.startsWith('/') ? method.slice(1) : method

		url.pathname = `${prefixNormalized}/${methodNormalized}`;

		if (params) {
			const entries = params instanceof URLSearchParams ? params.entries() : Object.entries(params);
			for (const [name, value] of entries) {

				switch (typeof value) {
					case 'string':
						url.searchParams.append(name, value);
						break;
					case 'number':
						url.searchParams.append(name, value.toString());
						break;
				}
			}
		}

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
			.then(result => ({ result: result as APIResult<T>, parseError: null }))
			.catch(err => ({ result: null, parseError: unwrapError(err)}));

		if (parseError) {
			return { data: null, error: new Error(`Parse API response: ${parseError.message}; status code: ${response.status}`) };
		}

		return result;
	};

	listCollections = async (params?: { id?: string } & Partial<Pagination>) => {
		return this.exec<ApiPage<CollectionMetadata>>('GET', '/collections', params);
	};

	listDecks = async (params?: { id?: string, collection_id?: string } & Partial<Pagination>) => {
		return this.exec<ApiPage<CardDeckMetadata>>('GET','/decks', params);
	};

	loadDeck = async (id: string) => {
		return this.exec<CardDeck>('GET',`/decks/${id}`);
	};
};

const unwrapError = (error: any): Error => {
	return error instanceof Error ? error : new Error(`${error}`);
};

export const useClient = () => new ApiClient('/api');
