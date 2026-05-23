package utils

import "strings"

type RecordMapper struct {
	columns map[string]int
}

func (mapper *RecordMapper) Get(record []string, column string) (string, bool) {

	idx, has := mapper.columns[column]
	if !has {
		return "", false
	}

	if idx >= len(record) {
		return "", false
	}

	return record[idx], true
}

func NewRecordMapper(headerRow []string) *RecordMapper {

	mapper := RecordMapper{columns: map[string]int{}}

	for idx, item := range headerRow {
		mapper.columns[strings.ToLower(strings.TrimSpace(item))] = idx
	}

	return &mapper
}
