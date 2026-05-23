package model

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	db_gen "github.com/maddsua/flippercardapp/db/generated"
	db_model "github.com/maddsua/flippercardapp/db/model"
)

type Response[T any] struct {
	Data  *T     `json:"data"`
	Error *Error `json:"error"`
}

func (resp *Response[T]) Write(wrt http.ResponseWriter) {

	if resp.Data == nil && resp.Error == nil {
		wrt.WriteHeader(http.StatusNoContent)
		return
	}

	wrt.Header().Set("Content-Type", "application/json")

	if resp.Error != nil {
		wrt.WriteHeader(resp.Error.StatusCode())
	}

	json.NewEncoder(wrt).Encode(resp)
}

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"-"`
}

func (err *Error) StatusCode() int {
	// min-max the error code to avoid whoopsie-daisies with invalid statuses
	return min(max(http.StatusBadRequest, err.Code), http.StatusNetworkAuthenticationRequired)
}

func (err *Error) Error() string {
	return err.Message
}

type ContentEntryMetaBase struct {
	Name        string                      `json:"name"`
	Description string                      `json:"description,omitempty"`
	Visibility  db_model.ResourceVisibility `json:"visibility"`
}

func (val *ContentEntryMetaBase) Valid() error {

	if val.Name = strings.TrimSpace(val.Name); val.Name == "" {
		return &Error{Message: "Summary invalid: 'name' field is empty"}
	} else if len(val.Name) > math.MaxUint8 {
		return &Error{Message: "Summary invalid: 'name' field too long"}
	} else if len(val.Description) > math.MaxUint8 {
		return &Error{Message: "Summary invalid: 'description' field too long"}
	}

	return nil
}

type CollectionMetadata struct {
	ContentEntryMetaBase
	ID      uuid.UUID `json:"id"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Size    int       `json:"size"`
}

func (meta *CollectionMetadata) Valid() error {

	if meta.ID == uuid.Nil {
		return &Error{Message: "Invalid collection ID"}
	}

	return meta.ContentEntryMetaBase.Valid()
}

func (meta *CollectionMetadata) FromRow(row db_gen.Collection) {
	meta.ContentEntryMetaBase = ContentEntryMetaBase{
		Name:        row.Name,
		Description: row.Description.String,
		Visibility:  row.Visibility,
	}
	meta.ID = row.ID
	meta.Created = row.CreatedAt.Time
	meta.Updated = row.UpdatedAt.Time
}

func (meta *CollectionMetadata) FromBatchRow(row db_gen.GetCollectionBatchRow) {
	meta.ContentEntryMetaBase = ContentEntryMetaBase{
		Name:        row.Name,
		Description: row.Description.String,
		Visibility:  row.Visibility,
	}
	meta.ID = row.ID
	meta.Created = row.CreatedAt.Time
	meta.Updated = row.UpdatedAt.Time
	meta.Size = int(row.Size)
}

type CollectionSearchResult struct {
	CollectionMetadata
	Rank int `json:"rank"`
}

type Collection struct {
	CollectionMetadata
	Decks []CardDeckMetadata `json:"decks"`
}

type CollectionBundle struct {
	CollectionMetadata
	Decks []CardDeckBundle `json:"decks"`
}

func (bundle *CollectionBundle) Valid() error {

	if err := bundle.CollectionMetadata.Valid(); err != nil {
		return err
	}

	for idx, deck := range bundle.Decks {

		if err := deck.CardDeckMetadata.Valid(); err != nil {
			return &Error{Message: fmt.Sprintf("Deck at index %d: %v", idx, err)}
		}
	}

	return nil
}

type CardDeckMetadata struct {
	ContentEntryMetaBase
	ID           uuid.UUID `json:"id"`
	CollectionID uuid.UUID `json:"collection_id"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
	Size         int       `json:"size"`
}

func (meta *CardDeckMetadata) FromRow(row db_gen.Deck) {
	meta.ContentEntryMetaBase = ContentEntryMetaBase{
		Name:        row.Name,
		Description: row.Description.String,
		Visibility:  row.Visibility,
	}
	meta.ID = row.ID
	meta.CollectionID = row.CollectionID
	meta.Created = row.CreatedAt.Time
	meta.Updated = row.UpdatedAt.Time
}

func (meta *CardDeckMetadata) FromBatchRow(row db_gen.GetDecksBatchRow) {
	meta.ContentEntryMetaBase = ContentEntryMetaBase{
		Name:        row.Name,
		Description: row.Description.String,
		Visibility:  row.Visibility,
	}
	meta.ID = row.ID
	meta.CollectionID = row.CollectionID
	meta.Created = row.CreatedAt.Time
	meta.Updated = row.UpdatedAt.Time
	meta.Size = int(row.Size)
}

type CardDeck struct {
	CardDeckMetadata
	Labels []string `json:"labels"`
	Cards  []Card   `json:"cards"`
}

type CardDeckBundle struct {
	CardDeckMetadata
	Cards []Card `json:"cards"`
}

type Card struct {
	db_model.CardNodeContent
	ID      uuid.UUID `json:"id"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

func (card *Card) FromRow(row db_gen.Card) {
	card.CardNodeContent = row.Content
	card.ID = row.ID
	card.Created = row.CreatedAt.Time
	card.Updated = row.UpdatedAt.Time
}

type CollectionPatch struct {
	ContentEntryMetaBase
}

type CardDeckPatch struct {
	CollectionID uuid.NullUUID         `json:"collection_id"`
	Meta         *ContentEntryMetaBase `json:"meta,omitempty"`
	Content      *CardDeckContentPatch `json:"content,omitempty"`
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

type ImageMetadata struct {
	ID               string    `json:"id"`
	Created          time.Time `json:"created"`
	Mimetype         string    `json:"mimetype"`
	SourceName       string    `json:"source_name"`
	SourceSha512Hash string    `json:"source_sha512_hash"`
	DataSha512Hash   string    `json:"data_sha512_hash"`
	DataSize         int       `json:"data_size"`
}

func (meta *ImageMetadata) FromRow(row db_gen.Image) {
	meta.ID = row.ID
	meta.Created = row.CreatedAt.Time
	meta.Mimetype = row.Mimetype
	meta.SourceName = row.SourceName
	meta.SourceSha512Hash = hex.EncodeToString(row.SourceSha512Hash)
	meta.DataSize = int(row.DataSize)
	meta.DataSha512Hash = hex.EncodeToString(row.DataSha512Hash)
}
