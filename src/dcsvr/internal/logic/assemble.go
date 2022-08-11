package logic

import (
	"database/sql"
	"fmt"
	"github.com/i-Things/things/src/dcsvr/dc"
	"github.com/i-Things/things/src/dcsvr/internal/repo/mysql"
)

func GetNullTime(time sql.NullTime) int64 {
	if time.Valid == false {
		return 0
	}
	return time.Time.Unix()
}

func DBToRPCFmt(db any) any {
	switch db.(type) {
	case *mysql.GroupInfo:
		di := db.(*mysql.GroupInfo)
		return &dc.GroupInfo{
			GroupID:     di.GroupID,
			Name:        di.Name,
			Uid:         di.Uid,
			CreatedTime: di.CreatedTime.Unix(),
		}
	case *mysql.GroupMember:
		pi := db.(*mysql.GroupMember)
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
