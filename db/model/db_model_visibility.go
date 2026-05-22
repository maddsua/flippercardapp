package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type ResourceVisibility byte

const (
	ResourceVisibilityPrivate = iota
	ResourceVisibilityHidden
	ResourceVisibilityPublic
)

func (visibility ResourceVisibility) String() string {
	switch visibility {
	case ResourceVisibilityPrivate:
		return "PRIVATE"
	case ResourceVisibilityHidden:
		return "HIDDEN"
	case ResourceVisibilityPublic:
		return "PUBLIC"
	default:
		return ""
	}
}

func (visibility ResourceVisibility) MarshalText() ([]byte, error) {
	val := visibility.String()
	if val == "" {
		return nil, nil
	}
	return []byte(val), nil
}

func (visibility *ResourceVisibility) UnmarshalText(data []byte) error {
	switch string(data) {
	case "PRIVATE":
		*visibility = ResourceVisibilityPrivate
	case "HIDDEN":
		*visibility = ResourceVisibilityHidden
	case "PUBLIC":
		*visibility = ResourceVisibilityPublic
	default:
		return fmt.Errorf("invalid ResourceVisibility value")
	}
	return nil
}

func (visibility ResourceVisibility) Value() (driver.Value, error) {
	return int64(visibility), nil
}

func (visibility *ResourceVisibility) Scan(src any) error {
	switch src := src.(type) {
	case int:
		return visibility.UnmarshalEnum(byte(src))
	case int64:
		return visibility.UnmarshalEnum(byte(src))
	default:
		return fmt.Errorf("unable to scan %T into ResourceVisibility", src)
	}
}

func (visibility *ResourceVisibility) UnmarshalEnum(val byte) error {
	switch val {
	case ResourceVisibilityPrivate,
		ResourceVisibilityHidden,
		ResourceVisibilityPublic:
		*visibility = ResourceVisibility(val)
	default:
		*visibility = ResourceVisibilityPrivate
	}
	return nil
}

type ResourceVisibilities []ResourceVisibility

func (rvset ResourceVisibilities) Value() (driver.Value, error) {

	if len(rvset) == 0 {
		return nil, nil
	}

	values := make([]any, len(rvset))
	for idx, val := range rvset {
		values[idx] = int(val)
	}

	return json.Marshal(values)
}

func (rvset *ResourceVisibilities) unmarshal(data []byte) error {

	var values []int
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}

	*rvset = make(ResourceVisibilities, len(values))
	for idx, val := range values {
		if err := (*rvset)[idx].Scan(val); err != nil {
			return err
		}
	}

	return nil
}

func (rvset *ResourceVisibilities) Scan(src any) error {

	if src == nil {
		*rvset = nil
		return nil
	}

	switch src := src.(type) {
	case []byte:
		return rvset.unmarshal(src)
	case string:
		return rvset.unmarshal([]byte(src))
	default:
		return fmt.Errorf("unable to scan %T into ResourceVisibilitySet", src)
	}
}
