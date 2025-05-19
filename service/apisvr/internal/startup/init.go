package startup

import (
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/things/service/apisvr/internal/logic/things/device"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
)

func Init(svcCtx *svc.ServiceContext) {
	device.Init(svcCtx)
}

func init() {
	var (
		TagsTypes []*types.Tag
		TagMap    map[string]string

		GM  map[string]types.IDsInfo
		GM2 map[string]*dm.IDsInfo
	)
	utils.AddConverter(
		utils.TypeConverter{SrcType: TagsTypes, DstType: TagMap, Fn: func(src interface{}) (dst interface{}, err error) {
			return logic.ToTagsMap(src.([]*types.Tag)), nil
		}},
		utils.TypeConverter{SrcType: TagMap, DstType: TagsTypes, Fn: func(src interface{}) (dst interface{}, err error) {
			return logic.ToTagsType(src.(map[string]string)), nil
		}},
		utils.TypeConverter{SrcType: GM, DstType: GM2, Fn: func(src interface{}) (dst interface{}, err error) {
			return utils.CopyMap2[dm.IDsInfo](src.(map[string]types.IDsInfo)), nil
		}},
		utils.TypeConverter{SrcType: GM2, DstType: GM, Fn: func(src interface{}) (dst interface{}, err error) {
			return utils.CopyMap3[types.IDsInfo](src.(map[string]*dm.IDsInfo)), nil
		}},
	)

}
