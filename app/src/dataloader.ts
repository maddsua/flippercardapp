import type { Page, Pagination, Result } from "./api";

export interface PaginatorState<T> {
	ready: boolean;
	locked: boolean;
	has_next: boolean;
	entries: T[];
	offset: number;
	error: string | null;
};

export type PageLoaderFn<T> = (pagination: Pagination) => Promise<Result<Page<T>>>;

export const genericPageState = <T>(): PaginatorState<T> => {
	return { entries: [], has_next: true, offset: 0, error: null, locked: false, ready: false };
};

const resetPage = <T>(state: PaginatorState<T>) => {
	if (state.locked) {
		return;
	}

	state.ready = false;
	state.error = null;
	state.has_next = true;
	state.offset = 0;
	state.entries = [];
};

const loadMore = async <T>(state: PaginatorState<T>, loader: PageLoaderFn<T>) => {
	if (state.locked) {
		return;
	}

	state.error = null;
	state.locked = true;

	await (async () => {

		const { data, error } = await loader({ offset: state.offset });
		if (!data || error) {
			state.error = error ? error.message : 'Empty page';
			return;
		}

		if (state.offset > 0) {
			state.offset += data.entries.length;
			state.entries.push(...data.entries);
		} else {
			state.offset = data.entries.length;
			state.entries = data.entries;
		}

		state.has_next = data.has_next;

	})();

	state.locked = false;
	state.ready = true;
};

export interface PageControls {
	more: () => Promise<void>;
	reload: () => Promise<void>;
}

export const pageControls = <T>(state: PaginatorState<T>, fn: PageLoaderFn<T>): PageControls => {
	return {
		more: () => loadMore(state, fn),
		reload: () => {
			resetPage(state);
			return loadMore(state, fn);
		},
	};
};
