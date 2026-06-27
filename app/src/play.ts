import type { CollectionMetadata } from "./api_models";

export interface DeckPlayStats {
	deck_id: string;
	collection_id: string | null;
	score: number;
};

export interface CollectionPlayStats {
	collection_id: string | null;
	avg_score: number;
	decks_played: number;
};

export const collectionCompletionMetric = (stats: Map<string, CollectionPlayStats>, meta: CollectionMetadata)  =>{
	const stat = stats.get(meta.id);
	if (!stat) {
		return 0;
	}
	return stat.avg_score * (stat.decks_played / (meta.size ?? 1));
};
