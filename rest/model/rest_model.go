package model

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
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

type CollectionPatch struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func (patch *CollectionPatch) Valid() error {

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
