package schema

import (
	"fmt"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	sq "gitee.com/unitedrhino/squirrel"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	"strings"
)

func GenArray(identifier string, num int) string {
	return fmt.Sprintf("%s.%d", identifier, num)
}

func WhereArray(db *gorm.DB, identifier string, pos string) *gorm.DB {
	arrs := strings.Split(identifier, ".")
	if len(arrs) <= 1 {
		return db
	}
	num, err := cast.ToIntE(arrs[1])
	if err != nil {
		if arrs[1] == "*" {
			return db
		}
		nums, err := utils.ParseNumberString(arrs[1])
		if err == nil {
			return db.Where(fmt.Sprintf("%s IN ?", pos), nums)
		}
		return db
	}
	db.Where(fmt.Sprintf("%s = ?", pos), num)
	return db
}

func WhereArray2(db sq.SelectBuilder, identifier string, pos string) sq.SelectBuilder {
	arrs := strings.Split(identifier, ".")
	if len(arrs) <= 1 {
		return db
	}
	num, err := cast.ToIntE(arrs[1])
	if err != nil {
		if arrs[1] == "*" {
			return db
		}
		nums, err := utils.ParseNumberString(arrs[1])
		if err == nil {
			return db.Where(fmt.Sprintf("%s IN (%s)", pos, stores.ArrayToSql(nums)))
		}
		return db
	}
	db = db.Where(fmt.Sprintf("%s = ?", pos), num)
	return db
}
func GetArray(identifier string) (ident string, num int, ok bool) {
	arrs := strings.Split(identifier, ".")
	if len(arrs) <= 1 {
		return identifier, 0, false
	}
	num, err := cast.ToIntE(arrs[1])
	if err != nil {
		return arrs[0], 0, false
	}
	return arrs[0], cast.ToInt(arrs[1]), true
}

type DataIDInfo struct {
	ID     string
	Nums   []int64
	Column string
}

func ParseDataID(identifier string) (*DataIDInfo, error) {
	arrs := strings.Split(identifier, ".")
	switch len(arrs) {
	case 0, 1:
		return &DataIDInfo{
			ID: identifier,
		}, nil
	case 2:
		num, err := cast.ToInt64E(arrs[1])
		if err != nil {
			if arrs[1] == "*" {
				return &DataIDInfo{
					ID: arrs[0],
				}, nil
			}
			nums, err := utils.ParseNumberString(arrs[1])
			if err == nil {
				return &DataIDInfo{
					ID:   arrs[0],
					Nums: nums,
				}, nil
			}
			return &DataIDInfo{
				ID:     arrs[0],
				Column: arrs[1],
			}, nil
		}
		return &DataIDInfo{
			ID:   arrs[0],
			Nums: []int64{num},
		}, nil
	case 3:
		num, err := cast.ToInt64E(arrs[1])
		if err != nil {
			if arrs[1] == "*" {
				return &DataIDInfo{
					ID:     arrs[0],
					Column: arrs[2],
				}, nil
			}
			nums, err := utils.ParseNumberString(arrs[1])
			if err != nil {
				return nil, err
			}
			return &DataIDInfo{
				ID:     arrs[0],
				Nums:   nums,
				Column: arrs[2],
			}, nil
		}
		return &DataIDInfo{
			ID:     arrs[0],
			Nums:   []int64{num},
			Column: arrs[2],
		}, nil
	}
	return &DataIDInfo{ID: identifier}, errors.Parameter.AddMsgf("get:%v,need:dataID or dataID.num or dataID.num.column", identifier)
}
