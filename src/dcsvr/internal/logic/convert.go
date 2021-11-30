package logic

import (
	"database/sql"
	"fmt"
	"gitee.com/godLei6/things/src/dcsvr/dc"
	"gitee.com/godLei6/things/src/dcsvr/model"
)

func GetNullTime(time sql.NullTime) int64 {
	if time.Valid == false {
		return 0
	}
	return time.Time.Unix()
}

func DBToRPCFmt(db interface{}) interface{} {
	switch db.(type) {
	case *model.GroupInfo:
		di := db.(*model.GroupInfo)
		return &dc.GroupInfo{
			GroupID:     di.GroupID,
			Name:        di.Name,
			Uid:         di.Uid,
			CreatedTime: di.CreatedTime.Unix(),
		}
	case *model.GroupMember:
		pi := db.(*model.GroupMember)
		dpi := &dc.GroupMember{
			GroupID:     pi.GroupID,            //产品名
			MemberID:    pi.MemberID,           //认证方式:0:账密认证,1:秘钥认证
			MemberType:  pi.MemberType,         //设备类型:0:设备,1:网关,2:子设备
			CreatedTime: pi.CreatedTime.Unix(), //创建时间
		}
		return dpi
	default:
		panic(fmt.Sprintf("ToRPCFmt not suppot:%#v", db))
	}
}
