package util

import "database/sql"

func GetSqlNullString(s string) sql.NullString {
	var res sql.NullString
	if len(s) <= 0 {
		res = sql.NullString{String: "", Valid: false}
	} else {
		res = sql.NullString{String: s, Valid: true}
	}
	return res
}

func GetSqlNullInt64(s int64) sql.NullInt64 {
	var res sql.NullInt64
	if s == 0 {
		res = sql.NullInt64{Int64: 0, Valid: false}
	} else {
		res = sql.NullInt64{Int64: s, Valid: false}
	}
	return res
}
