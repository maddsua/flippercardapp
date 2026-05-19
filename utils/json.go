package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
