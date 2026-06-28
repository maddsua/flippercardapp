package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

type ResourceVisibility byte

const (
	ResourceVisibilityPrivate = iota
	ResourceVisibilityHidden
	ResourceVisibilityPublic
)

func (val ResourceVisibility) String() string {
	switch val {
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

func (val *ResourceVisibility) Valid() bool {
	return val != nil && val.String() != ""
}

func (val ResourceVisibility) MarshalText() ([]byte, error) {
	text := val.String()
	if text == "" {
		return nil, nil
	}
	return []byte(text), nil
}

func (val *ResourceVisibility) UnmarshalText(data []byte) error {

	if enum, err := strconv.ParseInt(string(data), 10, 64); err == nil {
		return val.UnmarshalEnum(byte(enum))
	}

	switch string(data) {
	case "PRIVATE":
		*val = ResourceVisibilityPrivate
	case "HIDDEN":
		*val = ResourceVisibilityHidden
	case "PUBLIC":
		*val = ResourceVisibilityPublic
	default:
		return fmt.Errorf("invalid ResourceVisibility string value (%v)", string(data))
	}
	return nil
}

func (val ResourceVisibility) Value() (driver.Value, error) {
	return int64(val), nil
}

func (val *ResourceVisibility) Scan(src any) error {
	switch src := src.(type) {
	case int:
		return val.UnmarshalEnum(byte(src))
	case int64:
		return val.UnmarshalEnum(byte(src))
	case string:
		return val.UnmarshalText([]byte(src))
	case []byte:
		return val.UnmarshalText(src)
	default:
		return fmt.Errorf("unable to scan %T into ResourceVisibility", src)
	}
}

func (val *ResourceVisibility) UnmarshalEnum(enum byte) error {
	switch enum {
	case ResourceVisibilityPrivate,
		ResourceVisibilityHidden,
		ResourceVisibilityPublic:
		*val = ResourceVisibility(enum)
	default:
		return fmt.Errorf("invalid ResourceVisibility enum value (%v)", enum)
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

func (rvset *ResourceVisibilities) ScanJSON(data []byte) error {

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
		return rvset.ScanJSON(src)
	case string:
		return rvset.ScanJSON([]byte(src))
	default:
		return fmt.Errorf("unable to scan %T into ResourceVisibilitySet", src)
	}
}

func ResourceVisibilityFromPtr(ptr *ResourceVisibility) ResourceVisibility {
	if ptr.Valid() {
		return *ptr
	}
	return ResourceVisibilityPrivate
}
