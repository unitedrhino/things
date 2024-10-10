package startup

import (
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/things/service/apisvr/internal/logic/things/device"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
)

func Init(svcCtx *svc.ServiceContext) {
	device.Init(svcCtx)
}

func init() {
	var (
		TagsTypes []*types.Tag
		TagMap    map[string]string
	)
	utils.AddConverter(
		utils.TypeConverter{SrcType: TagsTypes, DstType: TagMap, Fn: func(src interface{}) (dst interface{}, err error) {
			return logic.ToTagsMap(src.([]*types.Tag)), nil
		}},
		utils.TypeConverter{SrcType: TagMap, DstType: TagsTypes, Fn: func(src interface{}) (dst interface{}, err error) {
			return logic.ToTagsType(src.(map[string]string)), nil
		}},
	)

}
