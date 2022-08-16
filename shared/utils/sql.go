package utils

import "database/sql"

func GetNullTime(time sql.NullTime) int64 {
	if time.Valid == false {
		return 0
	}
	return time.Time.Unix()
}
