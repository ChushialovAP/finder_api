package utils

import (
	"database/sql"
	"strconv"
)

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	if input_num == 0 {
		return ""
	}
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}
