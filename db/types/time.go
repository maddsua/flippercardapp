package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

func NewTime(val time.Time) Time {
	return Time{Time: val}
}

type Time struct {
	time.Time
}

func (val *Time) Scan(src any) (err error) {

	if src == nil {
		return fmt.Errorf("unable to scan a nil value into Time")
	}

	val.Time, err = scanTimeValue(src)
	return
}

func (val Time) Value() (driver.Value, error) {
	return marshalNanoTime(val.Time)
}

type NullTime struct {
	Time  time.Time
	Valid bool
}

func (val *NullTime) Scan(src any) (err error) {

	if src != nil {
		val.Time, err = scanTimeValue(src)
	}

	val.Valid = src != nil && err == nil

	return
}

func (val NullTime) Value() (driver.Value, error) {

	if !val.Valid {
		return nil, nil
	}

	return marshalNanoTime(val.Time)
}

func NewNullTime(val time.Time) NullTime {
	return NullTime{
		Time:  val,
		Valid: !val.IsZero(),
	}
}

func NewNullTimePtr(val *time.Time) NullTime {
	if val == nil {
		return NullTime{}
	}
	return NewNullTime(*val)
}

func scanTimeValue(src any) (time.Time, error) {

	switch src := src.(type) {
	case string:
		return unmarshalTextTime([]byte(src))
	case []byte:
		return unmarshalTextTime(src)
	case int:
		return unmarshalNanoTime(int64(src))
	case int64:
		return unmarshalNanoTime(src)
	default:
		return time.Time{}, fmt.Errorf("unable to scan %T into Time", src)
	}
}

func unmarshalTextTime(val []byte) (time time.Time, err error) {
	err = time.UnmarshalText(val)
	return
}

func unmarshalNanoTime(val int64) (time.Time, error) {
	if val < 0 {
		return time.Time{}, fmt.Errorf("invalid nano timestamp value: %v", val)
	}
	return time.Unix(0, val), nil
}

func marshalTextTime(val time.Time) ([]byte, error) {
	return val.MarshalText()
}

func marshalNanoTime(val time.Time) (int64, error) {
	return val.UnixNano(), nil
}
