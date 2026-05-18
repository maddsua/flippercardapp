package model

import (
	"errors"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	db_model "github.com/maddsua/flippercardapp/db/model"
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
	db_model.CardNodeContent
	ID      uuid.UUID `json:"id"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type CollectionDetailsPatch struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func (patch *CollectionDetailsPatch) Valid() error {

	if patch.Name = strings.TrimSpace(patch.Name); patch.Name == "" {
		return errors.New("name field is empty")
	} else if len(patch.Name) > math.MaxUint8 {
		return errors.New("'name' field too long")
	} else if len(patch.Description) > math.MaxUint8 {
		return errors.New("'description' field too long")
	}

	return nil
}

type CardDeckPatch struct {
	CollectionID uuid.NullUUID         `json:"collection_id"`
	Details      *CardDeckDetailsPatch `json:"details,omitempty"`
	Content      *CardDeckContentPatch `json:"content,omitempty"`
}

type CardDeckDetailsPatch struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func (patch *CardDeckDetailsPatch) Valid() error {

	if patch.Name = strings.TrimSpace(patch.Name); patch.Name == "" {
		return errors.New("name field is empty")
	} else if len(patch.Name) > math.MaxUint8 {
		return errors.New("'name' field too long")
	} else if len(patch.Description) > math.MaxUint8 {
		return errors.New("'description' field too long")
	}

	return nil
}

type CardDeckContentPatch struct {
	Cards []CardPatch `json:"cards"`
}

type CardPatch struct {
	db_model.CardNodeContent
	ID uuid.NullUUID `json:"id"`
}

type SignInParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
