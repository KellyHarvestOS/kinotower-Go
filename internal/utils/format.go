package utils

import (
	"database/sql"
	"time"
)

func TimeISO(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05.000Z")
}

func DateISO(t time.Time) string {
	return t.Format("2006-01-02")
}

func NullString(value sql.NullString) any {
	if value.Valid {
		return value.String
	}
	return nil
}

func BoolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
