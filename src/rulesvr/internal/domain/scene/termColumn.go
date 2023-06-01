package scene

import (
	"context"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	devicemsg "github.com/i-Things/things/src/disvr/client/devicemsg"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
)

type TermColumnType string

const (
	TermColumnTypeProperty TermColumnType = "property"
	TermColumnTypeEvent    TermColumnType = "event"
	//TermColumnTypeReportTime TermColumnType = "reportTime"
	TermColumnTypeSysTime TermColumnType = "sysTime"
)

// ColumnSchema 物模型类型 属性,事件
type ColumnSchema struct {
	ProductID  string   `json:"productID"` //产品id
	DeviceName string   `json:"deviceName"`
	DataID     []string `json:"dataID"`   //属性的id及事件的id aa.bb.cc
	TermType   CmpType  `json:"termType"` //动态条件类型  eq: 相等  not:不相等  btw:在xx之间  gt: 大于  gte:大于等于 lt:小于  lte:小于等于   in:在xx值之间
	Values     []string `json:"values"`   //条件值 参数根据动态条件类型会有多个参数
}

func (t TermColumnType) Validate() error {
	if !utils.SliceIn(t, TermColumnTypeProperty, TermColumnTypeEvent, TermColumnTypeSysTime) {
		return errors.Parameter.AddMsg("条件类型不支持:" + string(t))
	}
	return nil
}

func (c *ColumnSchema) Validate() error {
	if c == nil {
		return nil
	}
	if err := c.TermType.Validate(c.Values); err != nil {
		return err
	}
	if c.ProductID == "" {
		return errors.Parameter.AddMsg("触发设备类型中的产品id需要填写")
	}
	if c.DeviceName == "" {
		return errors.Parameter.AddMsg("触发设备类型中的设备名需要填写")
	}
	if len(c.DataID) == 0 {
		return errors.Parameter.AddMsg("触发设备类型中的标识符需要填写")
	}

	return nil
}
func (c *ColumnSchema) IsHit(ctx context.Context, columnType TermColumnType, repo TermRepo) bool {
	sm, err := repo.SchemaRepo.GetSchemaModel(ctx, c.ProductID)
	if err != nil {
		logx.WithContext(ctx).Errorf("%s.GetSchemaModel err:%v", utils.FuncName(), err)
		return false
	}
	var val string
	var dataType schema.DataType
	switch columnType {
	case TermColumnTypeProperty:
		info, err := repo.DeviceMsg.PropertyLatestIndex(ctx, &devicemsg.PropertyLatestIndexReq{ProductID: c.ProductID, DeviceName: c.DeviceName, DataIDs: c.DataID[:1]})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.PropertyLatestIndex err:%v", err)
			return false
		}
		if len(info.List) == 0 {
			logx.WithContext(ctx).Errorf("%s.PropertyLatestIndex err:dataID is not right:%v", c.DataID[0])
			return false
		}
		if info.List[0].Timestamp != 0 { //如果有值
			dataType = sm.Property[c.DataID[0]].Define.Type
			def := sm.Property[c.DataID[0]].Define
			switch def.Type {
			case schema.DataTypeStruct:
				if len(c.DataID) < 2 { //必须指定到结构体的成员
					return false
				}
				var dataMap = map[string]any{}
				utils.Unmarshal([]byte(info.List[0].Value), &dataMap)
				v, ok := dataMap[c.DataID[1]]
				if ok {
					val = cast.ToString(v)
					dataType = def.Spec[c.DataID[1]].DataType.Type
				}
			case schema.DataTypeArray:
				logx.WithContext(ctx).Errorf("%s scene not support array yet")
				return false
			default:
				val = info.List[0].Value
			}
		}
		return c.TermType.IsHit(dataType, val, c.Values)
	case TermColumnTypeEvent:
		logx.WithContext(ctx).Errorf("scene not support event yet")
		return false
	}
	return true
}
