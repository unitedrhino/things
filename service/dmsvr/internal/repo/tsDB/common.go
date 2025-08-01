package tsDB

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"strings"
)

func GenRedisPropertyFirstKey(productID string, deviceName string) string {
	return fmt.Sprintf("device:thing:property:first:%s:%s", productID, deviceName)
}

func GenRedisPropertyLastKey(productID string, deviceName string) string {
	return fmt.Sprintf("device:thing:property:last:%s:%s", productID, deviceName)
}

// 如果改变了才记录
func CheckIsChange(ctx context.Context, kv kv.Store, dev devices.Core, p *schema.Property, data msgThing.PropertyData) bool {
	//if p.RecordMode == schema.RecordModeAll || p.RecordMode == 0 {
	//	return true
	//}
	//if p.RecordMode == schema.RecordModeNone {
	//	return false
	//}
	data.Fmt()
	retStr, err := kv.Hget(GenRedisPropertyLastKey(dev.ProductID, dev.DeviceName), data.Identifier)
	if err != nil || retStr == "" {
		return true
	}
	var ret msgThing.PropertyData
	err = json.Unmarshal([]byte(retStr), &ret)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return true
	} else if msgThing.IsParamValEq(&p.Define, data.Param, ret.Param) { //相等不记录
		return false
	}
	return true
}

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
