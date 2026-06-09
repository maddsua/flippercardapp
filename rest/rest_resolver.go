package rest

import (
	"bytes"
	"context"
	"crypto/sha512"
	"database/sql"
	"fmt"
	"image"
	"io"
	"log/slog"
	"math"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/maddsua/flippercardapp/auth"
	db_pkg "github.com/maddsua/flippercardapp/db"
	db_gen "github.com/maddsua/flippercardapp/db/generated"
	db_model "github.com/maddsua/flippercardapp/db/model"
	"github.com/maddsua/flippercardapp/db/types"
	"github.com/maddsua/flippercardapp/rest/model"

	"github.com/maddsua/flippercardapp/utils"
	libwebp "github.com/pixiv/go-libwebp/webp"
	"golang.org/x/image/draw"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"
)

type resolver struct {
	db *db_pkg.Wrapper
}

func (rslv *resolver) LoadCardDeck(ctx context.Context, deckID uuid.UUID) (*model.CardDeck, error) {

	deck, err := rslv.db.GetDeckById(ctx, deckID)
	if db_pkg.IsNull(err) {
		return nil, &model.Error{Message: "deck not found", Code: http.StatusNotFound}
	} else if err != nil {
		return nil, InternalError("sqlc.GetDeckById", err)
	} else if err := EnforceResourceVisibility(ctx, deck.Visibility); err != nil {
		return nil, err
	}

	result := model.CardDeck{
		CardDeckMetadata: db_pkg.TransformRow[model.CardDeckMetadata](deck),
		Labels:           []string{deck.Name},
	}

	if version, err := rslv.getDeckLatestVersion(ctx, &rslv.db.Queries, deck.ID, deck.LatestVersionID); err != nil {
		return nil, err
	} else if version.Valid {
		result.VersionID = types.NewNullUUID(version.V.ID)
		result.Cards = version.V.Content.Cards
		result.Size = int(version.V.CardCount)
	}

	if collection, err := rslv.db.GetCollectionById(ctx, deck.CollectionID); err != nil {
		return nil, InternalError("sqlc.GetCollectionById", err)
	} else {
		result.Labels = append(result.Labels, collection.Name)
	}

	return &result, nil
}

func (rslv *resolver) getDeckLatestVersion(ctx context.Context, tx *db_gen.Queries, deckID uuid.UUID, latestVersionID uuid.NullUUID) (sql.Null[db_gen.DeckVersion], error) {

	if latestVersionID.Valid {

		if version, err := tx.GetDeckVersion(ctx, latestVersionID.UUID); err != nil && !db_pkg.IsNull(err) {
			return sql.Null[db_gen.DeckVersion]{}, InternalError("sqlc.GetDeckVersion", err)
		} else if err == nil {
			return sql.Null[db_gen.DeckVersion]{V: version, Valid: true}, nil
		}

		slog.Warn("REST DATA: Deck latest version wasn't found by it's referenced ID",
			slog.String("deck_id", deckID.String()))

		return sql.Null[db_gen.DeckVersion]{}, nil
	}

	version, err := tx.GetDeckLatestVersion(ctx, deckID)
	if err != nil && !db_pkg.IsNull(err) {
		return sql.Null[db_gen.DeckVersion]{}, InternalError("sqlc.GetDeckLatestVersion", err)
	}

	return sql.Null[db_gen.DeckVersion]{V: version, Valid: true}, nil
}

func (rslv *resolver) ListCardDeckVersions(ctx context.Context, deckID uuid.UUID, page PagePointers) (*Page[model.CardDeckVersionMetadata], error) {

	if perms, err := auth.For(ctx).Permissions(); err != nil {
		return nil, err
	} else if err := perms.AsTeamMember(); err != nil {
		return nil, err
	}

	deck, err := rslv.db.GetDeckById(ctx, deckID)
	if db_pkg.IsNull(err) {
		return nil, &model.Error{Message: "deck not found", Code: http.StatusNotFound}
	} else if err != nil {
		return nil, InternalError("sqlc.GetDeckById", err)
	}

	entries, err := rslv.db.GetDeckVersionsBatch(ctx, db_gen.GetDeckVersionsBatchParams{
		DeckID: deck.ID,
		Limit:  page.QueryLimit(),
		Offset: page.QueryOffset(),
	})
	if err != nil {
		return nil, InternalError("sqlc.GetDeckVersionsBatch", err)
	}

	result := TransformPage(page, entries, db_pkg.TransformBatchRow[model.CardDeckVersionMetadata, db_gen.GetDeckVersionsBatchRow])

	if deck.LatestVersionID.Valid {
		for idx, val := range result.Entries {
			if val.ID == deck.LatestVersionID.UUID {
				result.Entries[idx].IsLatest = true
				break
			}
		}
	}

	return result, nil
}

func (rslv *resolver) RollbackCardDeckVersion(ctx context.Context, deckID uuid.UUID, versionID uuid.UUID) (*model.CardDeckMetadata, error) {

	if perms, err := auth.For(ctx).Permissions(); err != nil {
		return nil, err
	} else if err := perms.AsContentEditor(); err != nil {
		return nil, err
	}

	tx, err := rslv.db.BeginTx(ctx)
	if err != nil {
		return nil, InternalError("sqlc.BeginTx", err)
	}
	defer tx.Rollback()

	deck, err := tx.GetDeckById(ctx, deckID)
	if db_pkg.IsNull(err) {
		return nil, &model.Error{Message: "Deck not found", Code: http.StatusNotFound}
	} else if err != nil {
		return nil, InternalError("sqlc.GetDeckById", err)
	}

	version, err := tx.GetDeckVersion(ctx, versionID)
	if db_pkg.IsNull(err) {
		return nil, &model.Error{Message: "Version not found", Code: http.StatusNotFound}
	} else if err != nil {
		return nil, InternalError("sqlc.GetDeckVersion", err)
	} else if version.ID == deck.LatestVersionID.UUID {
		return nil, &model.Error{Message: "Version is already the latest", Code: http.StatusBadRequest}
	}

	if deck, err = tx.SetDeckLatestVersion(ctx, db_gen.SetDeckLatestVersionParams{
		DeckID:          deck.ID,
		LatestVersionID: types.NewNullUUID(version.ID),
	}); err != nil {
		return nil, InternalError("sqlc.SetDeckLatestVersion", err)
	}

	if deck, err = tx.SetDeckUpdateTime(ctx, db_gen.SetDeckUpdateTimeParams{
		DeckID:    deck.ID,
		UpdatedAt: types.NewTime(time.Now()),
	}); err != nil {
		return nil, InternalError("sqlc.SetDeckUpdateTime", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, InternalError("sqlc.Commit", err)
	}

	result := db_pkg.TransformRow[model.CardDeckMetadata](deck)
	return &result, nil
}

func (rslv *resolver) ListCardDeckBatch(ctx context.Context, ids uuid.UUIDs, page PagePointers) (*Page[model.CardDeckMetadata], error) {

	entries, err := rslv.db.GetDecksBatch(ctx, db_gen.GetDecksBatchParams{
		IdsSet:        types.NewNullUUIDs(ids),
		VisibilitySet: ResourceVisibilityFilter(ctx, ids),
		Limit:         page.QueryLimit(),
		Offset:        page.QueryOffset(),
	})

	if err != nil {
		return nil, InternalError("sqlc.GetDecksBatch", err)
	}

	return TransformPage(page, entries, db_pkg.TransformBatchRow[model.CardDeckMetadata, db_gen.GetDecksBatchRow]), nil
}

func (rslv *resolver) LoadCollection(ctx context.Context, collectionID uuid.UUID) (*model.Collection, error) {

	collection, err := rslv.db.GetCollectionById(ctx, collectionID)
	if db_pkg.IsNull(err) {
		return nil, &model.Error{Message: "collection not found", Code: http.StatusNotFound}
	} else if err != nil {
		return nil, InternalError("sqlc.GetCollectionById", err)
	} else if err := EnforceResourceVisibility(ctx, collection.Visibility); err != nil {
		return nil, err
	}

	decks, err := rslv.db.GetDecksBatch(ctx, db_gen.GetDecksBatchParams{
		CollectionID:  types.NewNullUUID(collection.ID),
		VisibilitySet: ResourceVisibilityFilter(ctx, nil),
		Limit:         math.MaxInt32,
	})

	if err != nil {
		return nil, InternalError("sqlc.GetDecksBatch", err)
	}

	result := model.Collection{
		CollectionMetadata: db_pkg.TransformRow[model.CollectionMetadata](collection),
		Decks:              make([]model.CardDeckMetadata, len(decks)),
	}

	result.CollectionMetadata.Size = len(decks)

	for idx, val := range decks {
		result.Decks[idx].FromBatchRow(val)
	}

	return &result, nil
}

func (rslv *resolver) ListCollectionsBatch(ctx context.Context, ids uuid.UUIDs, page PagePointers) (*Page[model.CollectionMetadata], error) {

	entries, err := rslv.db.GetCollectionBatch(ctx, db_gen.GetCollectionBatchParams{
		IdsSet:        types.NewNullUUIDs(ids),
		VisibilitySet: ResourceVisibilityFilter(ctx, ids),
		Limit:         page.QueryLimit(),
		Offset:        page.QueryOffset(),
	})
	if err != nil {
		return nil, InternalError("sqlc.GetCollectionBatch", err)
	}

	return TransformPage(page, entries, db_pkg.TransformBatchRow[model.CollectionMetadata, db_gen.GetCollectionBatchRow]), nil
}

func (rslv *resolver) SearchCollections(ctx context.Context, term string, page PagePointers) (*Page[model.CollectionSearchResult], error) {

	matchIndex, err := rslv.fuzzyIndexCollection(ctx, term)
	if err != nil {
		return nil, err
	} else if len(matchIndex) == 0 {
		return WrapPage(page, []model.CollectionSearchResult{}), nil
	}

	entries, err := rslv.db.GetCollectionBatch(ctx, db_gen.GetCollectionBatchParams{
		IdsSet: types.NewNullUUIDs(UnwrapSearchIndex(matchIndex, page)),
		Limit:  page.QueryLimit(),
	})

	if err != nil {
		return nil, InternalError("sqlc.GetCollectionBatch", err)
	}

	result := TransformPage(page, entries, func(row db_gen.GetCollectionBatchRow) model.CollectionSearchResult {
		var next model.CollectionSearchResult
		next.FromBatchRow(row)
		next.Rank = matchIndex[row.ID]
		return next
	})

	sort.SliceStable(result.Entries, func(i, j int) bool {
		return result.Entries[i].Rank < result.Entries[j].Rank
	})

	return result, nil
}

func (rslv *resolver) fuzzyIndexCollection(ctx context.Context, term string) (map[uuid.UUID]int, error) {

	if term = strings.ToLower(strings.TrimSpace(term)); len(term) < 2 {
		return nil, &model.Error{Message: "Search term too short"}
	} else if len(term) > math.MaxUint8 {
		return nil, &model.Error{Message: "Search term too long", Code: http.StatusRequestEntityTooLarge}
	}

	tx, err := rslv.db.BeginTx(ctx)
	if err != nil {
		return nil, InternalError("sqlc.BeginTx", err)
	}
	defer tx.Rollback()

	const indexBatchSize = 100

	matchIndex := map[uuid.UUID]int{}

	for offset := 0; offset < math.MaxInt; offset += indexBatchSize {

		next, err := tx.GetCollectionSearchBatch(ctx, db_gen.GetCollectionSearchBatchParams{
			VisibilitySet: ResourceVisibilityFilter(ctx, nil),
			Offset:        int64(offset),
			Limit:         indexBatchSize,
		})

		if err != nil {
			return nil, InternalError("sqlc.GetCollectionSearchBatch", err)
		} else if len(next) == 0 {
			break
		}

		for _, item := range next {
			doc := strings.ToLower(item.Name)
			if rank := fuzzy.RankMatch(term, doc); rank >= 0 {
				matchIndex[item.ID] = rank
			}
		}
	}

	return matchIndex, nil
}

func (rslv *resolver) CreateContentCollection(ctx context.Context, params model.CollectionPatch) (*model.CollectionMetadata, error) {

	if perms, err := auth.For(ctx).Permissions(); err != nil {
		return nil, err
	} else if err := perms.AsContentEditor(); err != nil {
		return nil, err
	}

	tx, err := rslv.db.BeginTx(ctx)
	if err != nil {
		return nil, InternalError("sqlc.BeginTx", err)
	}
	defer tx.Rollback()

	if err := params.Valid(); err != nil {
		return nil, &model.Error{Message: fmt.Sprintf("Invalid collection details: %v", err)}
	}

	if exists, err := tx.CollectionNameExists(ctx, params.Name); err != nil {
		return nil, InternalError("sqlc.CollectionNameExists", err)
	} else if exists {
		return nil, &model.Error{
			Message: "Collection name already exists",
			Code:    http.StatusConflict,
		}
	}

	entry, err := tx.InsertCollection(ctx, db_gen.InsertCollectionParams{
		ID:          uuid.New(),
		CreatedAt:   types.NewTime(time.Now()),
		UpdatedAt:   types.NewTime(time.Now()),
		Name:        params.Name,
		Description: types.NewNullString(params.Description),
		Visibility:  params.Visibility,
	})

	if err != nil {
		return nil, InternalError("sqlc.InsertCollection", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, InternalError("sqlc.Commit", err)
	}

	result := db_pkg.TransformRow[model.CollectionMetadata](entry)
	return &result, nil
}

func (rslv *resolver) UpdateContentCollection(ctx context.Context, collectionID uuid.UUID, params model.CollectionPatch) (*model.CollectionMetadata, error) {

	if perms, err := auth.For(ctx).Permissions(); err != nil {
		return nil, err
	} else if err := perms.AsContentEditor(); err != nil {
		return nil, err
	}

	if err := params.Valid(); err != nil {
		return nil, &model.Error{Message: fmt.Sprintf("Invalid collection details: %v", err)}
	}

	tx, err := rslv.db.BeginTx(ctx)
	if err != nil {
		return nil, InternalError("sqlc.BeginTx", err)
	}
	defer tx.Rollback()

	entry, err := tx.GetCollectionById(ctx, collectionID)
	if db_pkg.IsNull(err) {
		return nil, &model.Error{Message: "Collection not found", Code: http.StatusNotFound}
	} else if err != nil {
		return nil, InternalError("sqlc.GetCollectionById", err)
	}

	if entry.Visibility != params.Visibility {
		if _, err = tx.UpdateCollectionChildrenVisibility(ctx, db_gen.UpdateCollectionChildrenVisibilityParams{
			CollectionID:  entry.ID,
			OldVisibility: entry.Visibility,
			NewVisibility: params.Visibility,
		}); err != nil {
			return nil, InternalError("sqlc.UpdateCollectionChildrenVisibility", err)
		}
	}

	if entry, err = tx.UpdateCollection(ctx, db_gen.UpdateCollectionParams{
		ID:          collectionID,
		Name:        params.Name,
		Description: types.NewNullString(params.Description),
		Visibility:  params.Visibility,
		UpdatedAt:   types.NewTime(time.Now()),
	}); err != nil {
		return nil, InternalError("sqlc.UpdateCollection", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, InternalError("sqlc.Commit", err)
	}

	result := db_pkg.TransformRow[model.CollectionMetadata](entry)
	return &result, nil
}

func (rslv *resolver) DeleteCollection(ctx context.Context, collectionID uuid.UUID, recursive bool) error {

	if perms, err := auth.For(ctx).Permissions(); err != nil {
		return err
	} else if err := perms.AsContentEditor(); err != nil {
		return err
	}

	tx, err := rslv.db.BeginTx(ctx)
	if err != nil {
		return InternalError("sqlc.BeginTx", err)
	}
	defer tx.Rollback()

	if !recursive {

		if count, err := tx.CollectionSize(ctx, collectionID); err != nil {

			if db_pkg.IsNull(err) {
				return &model.Error{Message: "Collection not found", Code: http.StatusNotFound}
			}

			return InternalError("sqlc.CollectionSize", err)

		} else if count > 0 {
			return &model.Error{Message: "Collection is not empty", Code: http.StatusConflict}
		}
	}

	if count, err := tx.DeleteCollection(ctx, collectionID); err != nil {
		return InternalError("sqlc.DeleteCollection", err)
	} else if count == 0 {
		return &model.Error{Message: "Collectcion not found", Code: http.StatusNotFound}
	}

	if err := tx.Commit(); err != nil {
		return InternalError("sqlc.Commit", err)
	}

	return nil
}

func (rslv *resolver) CreateCardDeck(ctx context.Context, params model.CardDeckPatch) (*model.CardDeckMetadata, error) {

	if perms, err := auth.For(ctx).Permissions(); err != nil {
		return nil, err
	} else if err := perms.AsContentEditor(); err != nil {
		return nil, err
	}

	if params.Meta == nil {
		return nil, &model.Error{Message: "Deck details must be provided"}
	} else if err := params.Meta.Valid(); err != nil {
		return nil, &model.Error{Message: fmt.Sprintf("Invalid deck details: %v", err)}
	} else if params.Content == nil || len(params.Content.Cards) == 0 {
		return nil, &model.Error{Message: "Deck has no cards in it"}
	} else if !params.CollectionID.Valid {
		return nil, &model.Error{Message: "Collection ID must be provided"}
	}

	tx, err := rslv.db.BeginTx(ctx)
	if err != nil {
		return nil, InternalError("sqlc.BeginTx", err)
	}
	defer tx.Rollback()

	if exists, err := tx.CollectionIDExists(ctx, params.CollectionID.UUID); err != nil {
		return nil, InternalError("sqlc.CollectionIDExists", err)
	} else if !exists {
		return nil, &model.Error{Message: "Collection ID not found", Code: http.StatusNotFound}
	}

	deck, err := tx.InsertDeck(ctx, db_gen.InsertDeckParams{
		ID:           uuid.New(),
		CollectionID: params.CollectionID.UUID,
		CreatedAt:    types.NewTime(time.Now()),
		UpdatedAt:    types.NewTime(time.Now()),
		Name:         params.Meta.Name,
		Description:  types.NewNullString(params.Meta.Description),
		Visibility:   params.Meta.Visibility,
	})

	if err != nil {
		return nil, InternalError("sqlc.InsertDeck", err)
	}

	result := db_pkg.TransformRow[model.CardDeckMetadata](deck)

	if params.Content != nil {

		version, err := rslv.addDeckVersion(ctx, &tx.Queries, deck.ID, "Deck created", params.Content.DeckVersionContent)
		if err != nil {
			return nil, err
		}

		result.VersionID = types.NewNullUUID(version.ID)
		result.Size = len(params.Content.Cards)
	}

	if err := tx.Commit(); err != nil {
		return nil, InternalError("sqlc.Commit", err)
	}

	return &result, nil
}

func (rslv *resolver) addDeckVersion(ctx context.Context, tx *db_gen.Queries, deckID uuid.UUID, label string, content db_model.DeckVersionContent) (db_gen.DeckVersion, error) {

	entry, err := tx.InsertDeckVersion(ctx, db_gen.InsertDeckVersionParams{
		ID:        uuid.New(),
		CreatedAt: types.NewTime(time.Now()),
		DeckID:    deckID,
		CardCount: int64(len(content.Cards)),
		Content:   content,
		Label:     types.NewNullString(label),
	})
	if err != nil {
		return db_gen.DeckVersion{}, InternalError("sqlc.InsertDeckVersion", err)
	}

	if _, err = tx.SetDeckLatestVersion(ctx, db_gen.SetDeckLatestVersionParams{
		DeckID:          entry.DeckID,
		LatestVersionID: types.NewNullUUID(entry.ID),
	}); err != nil {
		return db_gen.DeckVersion{}, InternalError("sqlc.SetDeckLatestVersion", err)
	}

	return entry, nil
}

func (rslv *resolver) UpdateCardDeck(ctx context.Context, deckID uuid.UUID, params model.CardDeckPatch) (*model.CardDeckMetadata, error) {

	if perms, err := auth.For(ctx).Permissions(); err != nil {
		return nil, err
	} else if err := perms.AsContentEditor(); err != nil {
		return nil, err
	}

	tx, err := rslv.db.BeginTx(ctx)
	if err != nil {
		return nil, InternalError("sqlc.BeginTx", err)
	}
	defer tx.Rollback()

	deck, err := tx.SetDeckUpdateTime(ctx, db_gen.SetDeckUpdateTimeParams{
		DeckID:    deckID,
		UpdatedAt: types.NewTime(time.Now()),
	})

	if db_pkg.IsNull(err) {
		return nil, &model.Error{Message: "Deck ID not found", Code: http.StatusNotFound}
	} else if err != nil {
		return nil, InternalError("sqlc.SetDeckUpdateTime", err)
	}

	if params.CollectionID.Valid {
		if exists, err := tx.CollectionIDExists(ctx, params.CollectionID.UUID); err != nil {
			return nil, InternalError("sqlc.CollectionIDExists", err)
		} else if !exists {
			return nil, &model.Error{Message: "Collection ID not found", Code: http.StatusNotFound}
		}
	}

	if params.Meta != nil {

		if err := params.Meta.Valid(); err != nil {
			return nil, &model.Error{Message: fmt.Sprintf("Invalid deck details: %v", err)}
		}

		if deck, err = tx.UpdateDeckMetadata(ctx, db_gen.UpdateDeckMetadataParams{
			ID:           deckID,
			CollectionID: params.CollectionID,
			Name:         params.Meta.Name,
			Description:  types.NewNullString(params.Meta.Description),
			Visibility:   params.Meta.Visibility,
		}); err != nil {
			return nil, InternalError("sqlc.UpdateDeckMetadata", err)
		}
	}

	result := db_pkg.TransformRow[model.CardDeckMetadata](deck)

	if params.Content != nil {

		version, err := rslv.addDeckVersion(ctx, &tx.Queries, deck.ID, params.Label, params.Content.DeckVersionContent)
		if err != nil {
			return nil, err
		}

		result.VersionID = types.NewNullUUID(version.ID)
		result.Size = len(params.Content.Cards)
	}

	if err := tx.Commit(); err != nil {
		return nil, InternalError("sqlc.Commit", err)
	}

	return &result, nil
}

func (rslv *resolver) DeleteDeck(ctx context.Context, deckID uuid.UUID) error {

	if perms, err := auth.For(ctx).Permissions(); err != nil {
		return err
	} else if err := perms.AsContentEditor(); err != nil {
		return err
	}

	if count, err := rslv.db.DeleteDeck(ctx, deckID); err != nil {
		return InternalError("sqlc.DeleteDeck", err)
	} else if count == 0 {
		return &model.Error{Message: "Deck not found", Code: http.StatusNotFound}
	}

	return nil
}

func (rslv *resolver) UploadImage(ctx context.Context, name string, reader io.Reader) (*model.ImageMetadata, error) {

	if perms, err := auth.For(ctx).Permissions(); err != nil {
		return nil, err
	} else if err := perms.AsContentEditor(); err != nil {
		return nil, err
	}

	const maxCanvasSize = 1440

	sourceHash := sha512.New()
	source, _, err := image.Decode(io.TeeReader(reader, sourceHash))
	if err != nil {
		return nil, &model.Error{Message: "Unable to decode image: " + err.Error(), Code: http.StatusUnprocessableEntity}
	}

	if entry, err := rslv.db.GetImageByHash(ctx, sourceHash.Sum(nil)); err == nil {
		var result model.ImageMetadata
		result.FromRow(entry)
		return &result, nil
	} else if !db_pkg.IsNull(err) {
		return nil, InternalError("sqlc.GetImageByHash", err)
	}

	sourceSize := source.Bounds().Size()
	if min(sourceSize.X, sourceSize.Y) < 100 || max(sourceSize.X, sourceSize.Y) > 10_000 {
		return nil, &model.Error{Message: "Invalid image dimensions", Code: http.StatusUnprocessableEntity}
	}

	if sourceAspect := float64(sourceSize.X) / float64(sourceSize.Y); sourceAspect > 3 || sourceAspect < 0.65 {
		return nil, &model.Error{Message: "Invalid image aspect ratio", Code: http.StatusUnprocessableEntity}
	}

	targetBounds := image.Rect(0, 0, sourceSize.X, sourceSize.Y)

	if size := targetBounds.Size(); size.X > maxCanvasSize {
		targetBounds = image.Rect(0, 0, maxCanvasSize, int(float64(size.Y)*(float64(maxCanvasSize)/float64(size.X))))
	}

	if size := targetBounds.Size(); size.Y > maxCanvasSize {
		targetBounds = image.Rect(0, 0, int(float64(size.X)*(float64(maxCanvasSize)/float64(size.Y))), maxCanvasSize)
	}

	canvas := image.NewRGBA(targetBounds)
	draw.CatmullRom.Scale(canvas, canvas.Rect, source, source.Bounds(), draw.Over, nil)

	var buffer bytes.Buffer

	encoderConfig, err := libwebp.ConfigPreset(libwebp.PresetDefault, 90)
	if err != nil {
		return nil, InternalError("webp.ConfigPreset", err)
	}

	if err := libwebp.EncodeRGBA(&buffer, canvas, encoderConfig); err != nil {
		return nil, InternalError("webp.EncodeRGBA", err)
	}

	dataHash := sha512.New()
	data, err := io.ReadAll(io.TeeReader(&buffer, dataHash))
	if err != nil {
		return nil, InternalError("io.TeeReader -> sha512", err)
	}

	if name == "" {
		name = "Image upload"
	}

	entry, err := rslv.db.InsertImage(ctx, db_gen.InsertImageParams{
		ID:               utils.NewRandomBase64Token(64),
		CreatedAt:        types.NewTime(time.Now()),
		Mimetype:         "image/webp",
		SourceName:       name,
		SourceSha512Hash: sourceHash.Sum(nil),
		DataSha512Hash:   dataHash.Sum(nil),
		DataSize:         int64(len(data)),
		Data:             data,
	})

	if err != nil {
		return nil, InternalError("sqlc.InsertImage", err)
	}

	var result model.ImageMetadata
	result.FromRow(entry)
	return &result, nil
}

func (rslv *resolver) ImageMetadata(ctx context.Context, mediaID string) (*model.ImageMetadata, error) {

	if perms, err := auth.For(ctx).Permissions(); err != nil {
		return nil, err
	} else if err := perms.AsTeamMember(); err != nil {
		return nil, err
	}

	entry, err := rslv.db.GetImageById(ctx, mediaID)
	if db_pkg.IsNull(err) {
		return nil, &model.Error{Message: "image not found", Code: http.StatusNotFound}
	} else if err != nil {
		return nil, InternalError("sqlc.GetImageById", err)
	}

	var result model.ImageMetadata
	result.FromRow(entry)
	return &result, nil
}

func (rslv *resolver) ImageBlob(ctx context.Context, mediaID string) (io.Reader, error) {

	if perms, err := auth.For(ctx).Permissions(); err != nil {
		return nil, err
	} else if err := perms.AsTeamMember(); err != nil {
		return nil, err
	}

	entry, err := rslv.db.GetImageById(ctx, mediaID)
	if db_pkg.IsNull(err) {
		return nil, &model.Error{Message: "image not found", Code: http.StatusNotFound}
	} else if err != nil {
		return nil, InternalError("sqlc.GetImageById", err)
	}

	return bytes.NewReader(entry.Data), nil
}
