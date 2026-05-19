package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	db_gen "github.com/maddsua/flippercardapp/db/generated"
	"github.com/maddsua/flippercardapp/db/model"
	"github.com/maddsua/flippercardapp/db/types"
	"github.com/maddsua/flippercardapp/utils"
	"golang.org/x/crypto/bcrypt"
)

type StateInitParams struct {
	RootUserName     string
	RootUserPassword string
	RestRootPassword bool
}

func InitDatabase(db *sql.DB, params StateInitParams) error {

	tx, err := NewWrapper(db).BeginTx(context.Background())
	if err != nil {
		return fmt.Errorf("sqlc.BeginTx: %v", err)
	}

	defer tx.Rollback()

	if err := initRootUser(context.Background(), tx, params); err != nil {
		return fmt.Errorf("root user init: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("sqlc.Commit: %v", err)
	}

	return nil
}

func initRootUser(ctx context.Context, tx *TxWrapper, params StateInitParams) error {

	// skip user init if it's not set, obviously
	if params.RootUserName == "" {
		return nil
	}

	rootUser, err := tx.GetUserByName(ctx, params.RootUserName)
	if IsNull(err) {
		return createRootUser(ctx, tx, params.RootUserName, params.RootUserPassword)
	} else if err != nil {
		return fmt.Errorf("sqlc.GetUserByName: %v", err)
	}

	if params.RestRootPassword && params.RootUserPassword != "" {
		slog.Warn("Resetting root user password",
			slog.String("id", rootUser.ID.String()),
			slog.String("username", rootUser.Name))
		return resetRootUserPassword(ctx, tx, rootUser.ID, params.RootUserPassword)
	}

	return nil
}

func createRootUser(ctx context.Context, tx *TxWrapper, username, password string) error {

	displayPassword := "**"
	if password == "" {
		password = utils.NewRandomBcryptPassword(64)
		displayPassword = password
	}

	pwHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypt.GenerateFromPassword: %v", err)
	}

	user, err := tx.InsertUser(ctx, db_gen.InsertUserParams{
		ID:           uuid.New(),
		CreatedAt:    types.NewTime(time.Now()),
		Name:         username,
		PasswordHash: pwHash,
		Permissions: model.NullUserPermissions{
			Permissions: model.UserPermissions{
				Administrative: true,
				ContentEdit:    true,
			},
			Valid: true,
		},
	})

	if err != nil {
		return fmt.Errorf("sqlc.InsertUser: %v", err)
	}

	slog.Warn("Add a root user",
		slog.String("id", user.ID.String()),
		slog.String("username", username),
		slog.String("password", displayPassword))

	return nil
}

func resetRootUserPassword(ctx context.Context, tx *TxWrapper, userID uuid.UUID, password string) error {

	pwHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypt.GenerateFromPassword: %v", err)
	}

	if _, err = tx.SetUserPassword(ctx, db_gen.SetUserPasswordParams{
		ID:           userID,
		PasswordHash: pwHash,
	}); err != nil {
		return fmt.Errorf("sqlc.SetUserPassword: %v", err)
	}

	return nil
}
