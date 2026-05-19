package utils

import "encoding/json"

func DecodeGenericJSON[T any](data []byte) (val T, err error) {
	err = json.Unmarshal(data, &val)
	return
}
