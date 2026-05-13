package types

import "database/sql"

func NewNullString(val string) sql.NullString {
	return sql.NullString{String: val, Valid: val != ""}
}
