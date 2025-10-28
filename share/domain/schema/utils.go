package schema

import (
	"fmt"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	sq "gitee.com/unitedrhino/squirrel"
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
	num, err := utils.ToIntE(arrs[1])
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
	num, err := utils.ToIntE(arrs[1])
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
	num, err := utils.ToIntE(arrs[1])
	if err != nil {
		return arrs[0], 0, false
	}
	return arrs[0], utils.ToInt(arrs[1]), true
}

type DataIDInfo struct {
	ID     string
	Nums   []int64
	NumStr string
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
		num, err := utils.ToInt64E(arrs[1])
		if err != nil {
			if arrs[1] == "*" {
				return &DataIDInfo{
					ID:     arrs[0],
					NumStr: arrs[1],
				}, nil
			}
			nums, err := utils.ParseNumberString(arrs[1])
			if err == nil {
				return &DataIDInfo{
					ID:     arrs[0],
					Nums:   nums,
					NumStr: arrs[1],
				}, nil
			}
			return &DataIDInfo{
				ID:     arrs[0],
				Column: arrs[1],
			}, nil
		}
		return &DataIDInfo{
			ID:     arrs[0],
			Nums:   []int64{num},
			NumStr: arrs[1],
		}, nil
	case 3:
		num, err := utils.ToInt64E(arrs[1])
		if err != nil {
			if arrs[1] == "*" {
				return &DataIDInfo{
					ID:     arrs[0],
					Column: arrs[2],
					NumStr: arrs[1],
				}, nil
			}
			nums, err := utils.ParseNumberString(arrs[1])
			if err != nil {
				return nil, err
			}
			return &DataIDInfo{
				ID:     arrs[0],
				Nums:   nums,
				NumStr: arrs[1],
				Column: arrs[2],
			}, nil
		}
		return &DataIDInfo{
			ID:     arrs[0],
			Nums:   []int64{num},
			NumStr: arrs[1],
			Column: arrs[2],
		}, nil
	}
	return &DataIDInfo{ID: identifier}, errors.Parameter.AddMsgf("get:%v,need:dataID or dataID.num or dataID.num.column", identifier)
}

func GetDataName(m *Model, dataID string) string {
	d, err := ParseDataID(dataID)
	if err != nil {
		return dataID
	}
	p := m.Property[d.ID]
	if p == nil {
		return dataID
	}
	ret := p.Name
	num := d.NumStr
	if num != "" {
		num = "." + num
	}
	def := &p.Define
	if len(d.Nums) > 0 && p.Define.Type == DataTypeArray {
		def = def.ArrayInfo
		ret = ret + num
		num = ""
	}
	if d.Column != "" {
		if def.Type == DataTypeStruct {
			pp := def.Spec[d.Column]
			if pp == nil {
				return ret
			}
			return ret + "." + pp.Name + num
		}
	}
	return ret + num
}
