package utils

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
)

func GetNullTime(time sql.NullTime) int64 {
	if time.Valid == false {
		return 0
	}
	return time.Time.Unix()
}
func TimeToNullTime(in *time.Time) sql.NullTime {
	if in == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Valid: true, Time: *in}
}

func ToNullTime(in int64) sql.NullTime {
	if in == 0 {
		return sql.NullTime{}
	}
	return sql.NullTime{Valid: true, Time: time.Unix(in, 0)}
}

// position 格式: POINT(100.101 50.894)
func GetPositionValue(position string) (float64, float64) {
	if position == "" || len(position) <= 7 {
		return 0, 0
	}
	sTemp := position[6 : len(position)-1]
	sCli := strings.Split(sTemp, " ")
	if len(sCli) < 2 {
		return 0, 0
	}
	Longitude, _ := strconv.ParseFloat(sCli[0], 64)
	Latitude, _ := strconv.ParseFloat(sCli[1], 64)

	return Longitude, Latitude
}

// 生成 "?,?,..." (有num个?)
func NewFillPlace(num int) string {
	return NewFillString(num, "?", ",")
}
