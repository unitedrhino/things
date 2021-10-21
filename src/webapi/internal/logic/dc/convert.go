package dc

import (
	"fmt"
	"gitee.com/godLei6/things/src/dcsvr/dc"
	"gitee.com/godLei6/things/src/webapi/internal/types"
	"github.com/golang/protobuf/ptypes/wrappers"
)

func GetNullVal(val *wrappers.StringValue) *string {
	if val == nil {
		return nil
	}
	return &val.Value
}

func RPCToApiFmt(rpc interface{}) interface{} {
	switch rpc.(type) {
	case *dc.GroupInfo:
		v := rpc.(*dc.GroupInfo)
		return &types.GroupInfo{
			GroupID:v.GroupID,     //组id
			Name:v.Name,        //组名
			Uid:v.Uid,         //管理员用户id
			CreatedTime:v.CreatedTime, //创建时间 只读
		}
	case *dc.GroupMember:
		v := rpc.(*dc.GroupMember)
		return &types.GroupMember{
			GroupID     :v.GroupID,              //组id
			MemberID    :v.MemberID,             //成员id
			MemberType  :v.MemberType,           //成员类型:1:设备 2:用户
			CreatedTime :v.CreatedTime, //创建时间 只读
		}
	default:
		panic(fmt.Sprintf("RPCToApiFmt not suppot:%#v", rpc))
	}
}
