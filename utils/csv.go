package utils

import "strings"

type RecordMapper struct {
	header []string
}

type RecordRow struct {
	mapped map[string]string
}

func (row *RecordRow) Get(column string) (string, bool) {
	if row.mapped == nil {
		return "", false
	}
	val, ok := row.mapped[column]
	return val, ok
}

func (mapper *RecordMapper) WithRow(row []string) *RecordRow {

	kv := map[string]string{}

	for idx := 0; idx < len(row) && idx < len(mapper.header); idx++ {
		kv[mapper.header[idx]] = row[idx]
	}

	return &RecordRow{mapped: kv}
}

func NewRecordMapper(headerRow []string) *RecordMapper {

	mapper := RecordMapper{header: make([]string, len(headerRow))}

	for idx, item := range headerRow {
		mapper.header[idx] = strings.ToLower(strings.TrimSpace(item))
	}

	return &mapper
}
