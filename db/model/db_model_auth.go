package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type UserPermissions struct {
	Administrative bool `json:"administrative,omitempty"`
	ContentEdit    bool `json:"content_edit,omitempty"`
}

func (perms UserPermissions) Value() (driver.Value, error) {
	return json.Marshal(perms)
}

func (perms *UserPermissions) Scan(src any) error {
	switch src := src.(type) {
	case []byte:
		return json.Unmarshal(src, perms)
	case string:
		return json.Unmarshal([]byte(src), perms)
	default:
		return fmt.Errorf("unable to scan %T into UserPermissions", src)
	}
}

type NullUserPermissions struct {
	Permissions UserPermissions
	Valid       bool
}

func (perms NullUserPermissions) Value() (driver.Value, error) {

	if !perms.Valid {
		return nil, nil
	}

	return perms.Permissions.Value()
}

func (perms *NullUserPermissions) Scan(src any) error {

	if src == nil {
		return nil
	}

	err := perms.Permissions.Scan(src)
	perms.Valid = err == nil
	return err
}
