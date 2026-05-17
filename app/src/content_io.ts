import type { Card, CardDeckMetadata, CollectionMetadata } from "./api_models";

export const pickUploadFiles = async (opts?: { multiple?: boolean, accept?: string[] }) => {

	const input = document.createElement('input');
	input.type = 'file';
	input.accept = opts?.accept?.join(',') || '';
	input.multiple = opts?.multiple || false;

	const filePromise = new Promise<FileList | null>((resolve) => {
		input.addEventListener('change', () => resolve(input.files));
		input.addEventListener('cancel', () => resolve(null));
	});

	input.click();
	const files = await filePromise;
	input.remove();

	return files;
};

export const downloadFile = (content: BlobPart, filename: string, opts?: { mimetype?: string }) => {

	const blob = new Blob([content], {
		type: opts?.mimetype || 'application/json'
	});

	const url = window.URL.createObjectURL(blob);
	const link = document.createElement('a');

	link.href = url;
	link.download = filename;

	link.click();

	window.URL.revokeObjectURL(url);
};

export interface CollectionBundle {
	type: 'collection_bundle';
	content: CollectionBundleContent[];
};

export interface CollectionBundleContent {
	meta: Omit<CollectionMetadata, 'id' | 'size'>;
	decks: CollectionBundleDeckContent[];
};

export interface CollectionBundleDeckContent {
	meta: Omit<CardDeckMetadata, 'id' | 'created' | 'updated' | 'size' | 'collection_id'>;
	cards: Array<Omit<Card, 'id' | 'created' | 'updated'>>;
};

export interface DeckBundle {
	type: 'deck_bundle';
	content: CollectionBundleDeckContent[];
};

export type ContentBundle = CollectionBundle | DeckBundle;
