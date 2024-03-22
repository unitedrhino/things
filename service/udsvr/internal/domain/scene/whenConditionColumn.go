package scene

import (
	"context"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	devicemsg "github.com/i-Things/things/service/dmsvr/client/devicemsg"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

type TermColumnType string

const (
	TermColumnTypeProperty TermColumnType = "property"
	TermColumnTypeEvent    TermColumnType = "event"
	//TermColumnTypeReportTime TermColumnType = "reportTime"
	TermColumnTypeSysTime TermColumnType = "sysTime"
)

// TermProperty 物模型类型 属性
type TermProperty struct {
	ProductID   string   `json:"productID"` //产品id
	DeviceName  string   `json:"deviceName"`
	DeviceAlias string   `json:"deviceAlias"`
	DataID      string   `json:"dataID"` //属性的id   aa.bb.cc
	DataName    string   `json:"dataName"`
	TermType    CmpType  `json:"termType"` //动态条件类型  eq: 相等  not:不相等  btw:在xx之间  gt: 大于  gte:大于等于 lt:小于  lte:小于等于   in:在xx值之间
	Values      []string `json:"values"`   //条件值 参数根据动态条件类型会有多个参数
}

func (t TermColumnType) Validate() error {
	if !utils.SliceIn(t, TermColumnTypeProperty, TermColumnTypeEvent, TermColumnTypeSysTime) {
		return errors.Parameter.AddMsg("条件类型不支持:" + string(t))
	}
	return nil
}

func (c *TermProperty) Validate(repo ValidateRepo) error {
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
	c.DeviceAlias = GetDeviceAlias(repo.Ctx, repo.DeviceCache, c.ProductID, c.DeviceName)
	v, err := repo.ProductSchemaCache.GetData(repo.Ctx, c.ProductID)
	if err != nil {
		return err
	}
	p := v.Property[c.DataID]
	if p == nil {
		return errors.Parameter.AddMsg("dataID不存在")
	}
	c.DataName = p.Name
	return nil
}
func (c *TermProperty) IsHit(ctx context.Context, columnType TermColumnType, repo TermRepo) bool {
	sm, err := repo.SchemaRepo.GetSchemaModel(ctx, c.ProductID)
	if err != nil {
		logx.WithContext(ctx).Errorf("%s.GetSchemaModel err:%v", utils.FuncName(), err)
		return false
	}
	var val string
	var dataType schema.DataType
	switch columnType {
	case TermColumnTypeProperty:
		dataID := strings.Split(c.DataID, ".")
		info, err := repo.DeviceMsg.PropertyLogLatestIndex(ctx, &devicemsg.PropertyLogLatestIndexReq{ProductID: c.ProductID, DeviceName: c.DeviceName, DataIDs: dataID[:1]})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.PropertyLatestIndex err:%v", utils.FuncName(), err)
			return false
		}
		if len(info.List) == 0 {
			logx.WithContext(ctx).Errorf("%s.PropertyLatestIndex err:dataID is not right:%v", utils.FuncName(), c.DataID[0])
			return false
		}
		if info.List[0].Timestamp != 0 { //如果有值
			dataType = sm.Property[dataID[0]].Define.Type
			def := sm.Property[dataID[0]].Define
			switch def.Type {
			case schema.DataTypeStruct:
				if len(dataID) < 2 { //必须指定到结构体的成员
					return false
				}
				var dataMap = map[string]any{}
				utils.Unmarshal([]byte(info.List[0].Value), &dataMap)
				v, ok := dataMap[dataID[1]]
				if ok {
					val = cast.ToString(v)
					dataType = def.Spec[dataID[1]].DataType.Type
				}
			case schema.DataTypeArray:
				logx.WithContext(ctx).Errorf("%s scene not support array yet", utils.FuncName())
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
