package types

import "github.com/google/uuid"

func NewNullUUID(val uuid.UUID) uuid.NullUUID {
	return uuid.NullUUID{UUID: val, Valid: true}
}
