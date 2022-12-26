package utils

import (
	"database/sql"
	"strconv"
	"strings"
)

func GetNullTime(time sql.NullTime) int64 {
	if time.Valid == false {
		return 0
	}
	return time.Time.Unix()
}

func GetPositionValue(position string) (float64, float64) {
	if position == "" || len(position) <= 7 {
		return 0, 0
	}
	sTemp := position[6 : len(position)-2]
	sCli := strings.Split(sTemp, " ")
	if len(sCli) < 2 {
		return 0, 0
	}
	Longitude, _ := strconv.ParseFloat(sCli[0], 64)
	Latitude, _ := strconv.ParseFloat(sCli[1], 64)

	return Longitude, Latitude
}
