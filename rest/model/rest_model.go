package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CollectionMetadata struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Size        int       `json:"size"`
}

type CollectionSearchResult struct {
	CollectionMetadata
	Rank int `json:"rank"`
}

type Collection struct {
	CollectionMetadata
	Decks []CardDeckMetadata `json:"decks"`
}

type CardDeckMetadata struct {
	ID           uuid.UUID `json:"id"`
	CollectionID uuid.UUID `json:"collection_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
	Size         int       `json:"size"`
}

type CardDeck struct {
	CardDeckMetadata
	Labels []string `json:"labels"`
	Cards  []Card   `json:"cards"`
}

type Card struct {
	ID      uuid.UUID       `json:"id"`
	Created time.Time       `json:"created"`
	Updated time.Time       `json:"updated"`
	Content json.RawMessage `json:"content"`
}
