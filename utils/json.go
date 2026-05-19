package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"reflect"
	"strings"
)

func DecodeGenericJSON[T any](data []byte) (val T, err error) {
	err = json.Unmarshal(data, &val)
	return
}

// Extracts the first field value that appears on a serialized objecy.
// Used for dynamic object type sniffing
func ExtractJSONField[T any](data []byte, field string) (any, error) {

	decoder := json.NewDecoder(bytes.NewReader(data))

	var depth int

	for {

		token, err := decoder.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("decoder.Token: %v", err)
		}

		switch token := token.(type) {

		case json.Delim:

			//	don't process arrays
			if depth == 0 && token == '[' {
				return nil, nil
			}

			switch token {
			case '{', '[':
				depth++
			case '}', ']':
				depth--
			}

		case string:
			if depth == 1 && token == field {
				var val T
				if err := decoder.Decode(&val); err != nil {
					return nil, fmt.Errorf("decoder.Decode: %v", err)
				}
				return val, nil
			}
		}
	}

	return nil, nil
}

// Pulls first-level json-marked fields out of a struct into a map.
// Used for dynamic object type sniffing
func ExtractStructJSONFields(src any) (map[string]any, error) {

	if src == nil {
		return map[string]any{}, nil
	}

	elem := reflect.ValueOf(src)

	switch elem.Kind() {

	case reflect.Struct:
		break

	case reflect.Pointer, reflect.Interface:

		if elem.IsNil() {
			return map[string]any{}, nil
		}

		return ExtractStructJSONFields(elem.Elem().Interface())

	default:
		return nil, fmt.Errorf("unsupported source value type: %T", src)
	}

	elemType := elem.Type()

	fields := map[string]any{}

	for idx := 0; idx < elem.NumField(); idx++ {

		fieldType := elemType.Field(idx)
		fieldType.Tag.Get("json")

		if !fieldType.IsExported() {
			continue
		}

		tag := fieldType.Tag.Get("json")
		if tag == "-" {
			continue
		}

		field := elem.Field(idx)

		name, opts := parseJsonTag(tag)
		if opts["omitempty"] && isEmptyJsonValue(field) {
			continue
		} else if opts["omitzero"] && isZeroJsonValue(field) {
			continue
		}

		if name == "" {
			name = fieldType.Name
		}

		if fieldType.Anonymous {

			subfields, err := ExtractStructJSONFields(field.Interface())
			if err != nil {
				return nil, err
			}

			maps.Copy(fields, subfields)
			continue
		}

		fields[name] = field.Interface()
	}

	return fields, nil
}

func parseJsonTag(tag string) (name string, optionSet map[string]bool) {

	optionSet = map[string]bool{}

	name, optList, _ := strings.Cut(tag, ",")

	for opt := range strings.SplitSeq(optList, ",") {
		optionSet[strings.TrimSpace(opt)] = true
	}

	return strings.TrimSpace(name), optionSet
}

func isEmptyJsonValue(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return val.Len() == 0
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Interface, reflect.Pointer:
		return val.IsZero()
	}
	return false
}

func isZeroJsonValue(val reflect.Value) bool {
	if zeroer, ok := val.Interface().(interface{ IsZero() bool }); ok {
		return zeroer == nil || zeroer.IsZero()
	}
	return val.IsZero()
}
