package tsDB

import (
	"fmt"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"strings"
)

func GroupFilter(db *stores.DB, belongGroup map[string]def.IDsInfo) *stores.DB {
	if len(belongGroup) == 0 {
		return db
	}
	subQuery := stores.GetTsConn(ctxs.WithRoot(nil)).Model(&relationDB.DmGroupDevice{}).Select("product_id, device_name")

	for purpose, idInfo := range belongGroup {
		if len(idInfo.IDs) > 0 {
			subQuery = subQuery.Where("purpose=? and group_id in ?", purpose, idInfo.IDs)
		}
		if len(idInfo.IDPaths) > 0 {
			var sqls []string
			for _, idPath := range idInfo.IDPaths {
				sqls = append(sqls, fmt.Sprintf(" group_id_path like '%s' ", idPath+"%"))
			}
			subQuery = subQuery.Where(fmt.Sprintf("purpose=? and (%s)", strings.Join(sqls, " or ")), purpose)
		}
	}
	db = db.Where("(product_id, device_name) in (?)",
		subQuery)
	return db
}
func GroupFilter2(db *stores.DB, belongGroup map[string]def.IDsInfo) *stores.DB {
	if len(belongGroup) == 0 {
		return db
	}
	subQuery := stores.GetTsConn(ctxs.WithRoot(nil)).Model(&relationDB.DmGroupDevice{}).Select("product_id, device_name")

	for purpose, idInfo := range belongGroup {
		if len(idInfo.IDs) > 0 {
			subQuery = subQuery.Where("purpose=? and group_id in ?", purpose, idInfo.IDs)
		}
		if len(idInfo.IDPaths) > 0 {
			var sqls []string
			for _, idPath := range idInfo.IDPaths {
				sqls = append(sqls, fmt.Sprintf(" group_id_path like '%s' ", idPath+"%"))
			}
			subQuery = subQuery.Where(fmt.Sprintf("purpose=? and (%s)", strings.Join(sqls, " or ")), purpose)
		}
	}
	db = db.Where("(tb.product_id, tb.device_name) in (?)",
		subQuery)
	return db
}
