package model

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
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

type Hash []byte

func (hash Hash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(hash)), nil
}

func (hash *Hash) UnmarshalText(data []byte) (err error) {
	*hash, err = hex.DecodeString(string(data))
	return
}

type SignInParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ContentEntryMeta struct {
	db_model.ContentSummary
	ID         uuid.UUID                   `json:"id"`
	Created    time.Time                   `json:"created"`
	Updated    time.Time                   `json:"updated"`
	Visibility db_model.ResourceVisibility `json:"visibility"`
}

type CollectionMeta struct {
	ContentEntryMeta
	ContentUpdated *time.Time `json:"content_updated,omitempty"`
	Size           int        `json:"size"`
	ThemeColor     string     `json:"theme_color,omitempty"`
}

func (model *CollectionMeta) FromRow(row db_gen.Collection) {

	model.ContentEntryMeta = ContentEntryMeta{

		ContentSummary: db_model.ContentSummary{
			Name:        row.Name,
			Description: row.Description.String,
		},

		ID:      row.ID,
		Created: row.CreatedAt.Time,
		Updated: row.UpdatedAt.Time,

		Visibility: row.Visibility,
	}

	if row.ContentUpdatedAt.Valid {
		model.ContentUpdated = &row.ContentUpdatedAt.Time
	}

	model.ThemeColor = row.ThemeColor.String
}

func (model *CollectionMeta) FromBatchRow(row db_gen.GetCollectionBatchRow) {

	model.ContentEntryMeta = ContentEntryMeta{

		ContentSummary: db_model.ContentSummary{
			Name:        row.Name,
			Description: row.Description.String,
		},

		ID:         row.ID,
		Created:    row.CreatedAt.Time,
		Updated:    row.UpdatedAt.Time,
		Visibility: row.Visibility,
	}

	if row.ContentUpdatedAt.Valid {
		model.ContentUpdated = &row.ContentUpdatedAt.Time
	}

	model.Size = int(row.Size)
	model.ThemeColor = row.ThemeColor.String
}

type CollectionSearchResult struct {
	CollectionMeta
	Rank int `json:"rank"`
}

type Collection struct {
	CollectionMeta
	Decks []CardDeckMeta `json:"decks"`
}

type CardDeckMeta struct {
	ContentEntryMeta
	CollectionID         uuid.UUID     `json:"collection_id"`
	VersionID            uuid.NullUUID `json:"version_id"`
	Size                 int           `json:"size"`
	CollectionThemeColor string        `json:"collection_theme_color"`
}

func (model *CardDeckMeta) FromRow(row db_gen.Deck) {

	model.ContentEntryMeta = ContentEntryMeta{

		ContentSummary: db_model.ContentSummary{
			Name:        row.Name,
			Description: row.Description.String,
		},

		ID:         row.ID,
		Created:    row.CreatedAt.Time,
		Updated:    row.UpdatedAt.Time,
		Visibility: row.Visibility,
	}

	model.CollectionID = row.CollectionID
	model.VersionID = row.LatestVersionID
}

func (model *CardDeckMeta) FromBatchRow(row db_gen.GetDecksBatchRow) {

	model.ContentEntryMeta = ContentEntryMeta{

		ContentSummary: db_model.ContentSummary{
			Name:        row.Name,
			Description: row.Description.String,
		},

		ID:         row.ID,
		Created:    row.CreatedAt.Time,
		Updated:    row.UpdatedAt.Time,
		Visibility: row.Visibility,
	}

	model.CollectionID = row.CollectionID
	model.VersionID = row.LatestVersionID

	model.Size = int(row.Size.Int64)
	model.CollectionThemeColor = row.CollectionThemeColor.String
}

type CardDeck struct {
	CardDeckMeta
	Labels []string            `json:"labels"`
	Cards  []db_model.CardNode `json:"cards"`
}

type CardDeckVersionMetaBase struct {
	ID        uuid.UUID `json:"id"`
	Created   time.Time `json:"created"`
	DeckID    uuid.UUID `json:"deck_id"`
	CardCount int       `json:"card_count"`
	IsLatest  bool      `json:"is_latest"`
	Label     string    `json:"label,omitempty"`
}

func (model *CardDeckVersionMetaBase) FromRow(row db_gen.DeckVersion) {
	*model = CardDeckVersionMetaBase{
		ID:        row.ID,
		Created:   row.CreatedAt.Time,
		DeckID:    row.DeckID,
		CardCount: int(row.CardCount),
		Label:     row.Label.String,
	}
}

type CardDeckVersionMeta struct {
	CardDeckVersionMetaBase
	Summary *db_model.ContentSummary `json:"summary,omitempty"`
}

func (model *CardDeckVersionMeta) FromRow(row db_gen.DeckVersion) {

	model.CardDeckVersionMetaBase.FromRow(row)

	if !row.Content.Summary.Empty() {
		model.Summary = &row.Content.Summary
	}
}

type CardDeckVersion struct {
	CardDeckVersionMetaBase
	Content db_model.CardDeckVersionContent `json:"content"`
}

func (model *CardDeckVersion) FromRow(row db_gen.DeckVersion) {

	model.CardDeckVersionMetaBase.FromRow(row)
	model.Content = row.Content

	// ensure that an empty array is always returned instead of a null
	if model.Content.Cards == nil {
		model.Content.Cards = make([]db_model.CardNode, 0)
	}
}

type CollectionPatch struct {
	db_model.ContentSummary
	Visibility db_model.ResourceVisibility `json:"visibility"`
	ThemeColor string                      `json:"theme_color"`
}

type CardDeckPatch struct {
	CollectionID uuid.NullUUID                `json:"collection_id"`
	Label        string                       `json:"label,omitempty"`
	Summary      *db_model.ContentSummary     `json:"summary,omitempty"`
	Visibility   *db_model.ResourceVisibility `json:"visibility,omitempty"`
	Content      *CardDeckContentPatch        `json:"content,omitempty"`
}

type CardDeckContentPatch struct {
	Cards []db_model.CardNode `json:"cards"`
}

type ImageMeta struct {
	ID               string    `json:"id"`
	Created          time.Time `json:"created"`
	Mimetype         string    `json:"mimetype"`
	SourceName       string    `json:"source_name"`
	SourceSha512Hash Hash      `json:"source_sha512_hash"`
	DataSha512Hash   Hash      `json:"data_sha512_hash"`
	DataSize         int       `json:"data_size"`
}

func (meta *ImageMeta) FromRow(row db_gen.Image) {
	meta.ID = row.ID
	meta.Created = row.CreatedAt.Time
	meta.Mimetype = row.Mimetype
	meta.SourceName = row.SourceName
	meta.SourceSha512Hash = row.SourceSha512Hash
	meta.DataSize = int(row.DataSize)
	meta.DataSha512Hash = row.DataSha512Hash
}
