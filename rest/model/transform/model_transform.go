package transform

import (
	db_gen "github.com/maddsua/flippercardapp/db/generated"
	"github.com/maddsua/flippercardapp/rest/model"
)

func ToPtr[T any](val T) *T {
	return &val
}

func CardDeckMetadataFromBatchRow(row db_gen.GetDecksBatchRow) model.CardDeckMetadata {
	return model.CardDeckMetadata{
		ID:           row.ID,
		CollectionID: row.CollectionID,
		Name:         row.Name,
		Description:  row.Description.String,
		Created:      row.CreatedAt.Time,
		Updated:      row.UpdatedAt.Time,
		Size:         int(row.Size),
	}
}

func CollectionMetadataFromBatchRow(row db_gen.GetCollectionBatchRow) model.CollectionMetadata {
	return model.CollectionMetadata{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description.String,
		Created:     row.CreatedAt.Time,
		Updated:     row.UpdatedAt.Time,
		Size:        int(row.Size),
	}
}

func CardFromRow(row db_gen.Card) model.Card {
	return model.Card{
		ID:      row.ID,
		Created: row.CreatedAt.Time,
		Updated: row.UpdatedAt.Time,
		Content: row.Content,
	}
}

func CardDeckMetadataFromRow(row db_gen.Deck) model.CardDeckMetadata {
	return model.CardDeckMetadata{
		ID:           row.ID,
		CollectionID: row.CollectionID,
		Name:         row.Name,
		Description:  row.Description.String,
		Created:      row.CreatedAt.Time,
		Updated:      row.UpdatedAt.Time,
	}
}

func CollectionMetadataFromRow(row db_gen.Collection) model.CollectionMetadata {
	return model.CollectionMetadata{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description.String,
		Created:     row.CreatedAt.Time,
		Updated:     row.UpdatedAt.Time,
	}
}
