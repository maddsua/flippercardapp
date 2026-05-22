package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

func NewNullUUID(val uuid.UUID) uuid.NullUUID {
	return uuid.NullUUID{UUID: val, Valid: true}
}

func NewNullUUIDs(val uuid.UUIDs) NullUUIDs {
	return NullUUIDs{UUIDs: UUIDs(val), Valid: len(val) > 0}
}

type UUIDs []uuid.UUID

func (val *UUIDs) Scan(src any) (err error) {

	if src == nil {
		return fmt.Errorf("unable to scan a nil value into UUIDs")
	}

	switch src := src.(type) {
	case []byte:
		return json.Unmarshal(src, val)
	case string:
		return json.Unmarshal([]byte(src), val)
	default:
		return fmt.Errorf("unable to scan %T into UUIDs", src)
	}
}

func (val UUIDs) Value() (driver.Value, error) {
	return json.Marshal(val)
}

type NullUUIDs struct {
	UUIDs UUIDs
	Valid bool
}

func (val *NullUUIDs) Scan(src any) (err error) {

	if src == nil {
		val.UUIDs = nil
		val.Valid = false
		return
	}

	switch src := src.(type) {
	case []byte:
		err = json.Unmarshal(src, &val.UUIDs)
	case string:
		err = json.Unmarshal([]byte(src), &val.UUIDs)
	default:
		err = fmt.Errorf("unable to scan %T into NullUUIDs", src)
	}

	val.Valid = err == nil

	return
}

func (val NullUUIDs) Value() (driver.Value, error) {
	if !val.Valid {
		return nil, nil
	}
	return json.Marshal(val.UUIDs)
}
