package scene

import (
	"context"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	devicemsg "gitee.com/unitedrhino/things/service/dmsvr/client/devicemsg"
	"gitee.com/unitedrhino/things/service/dmsvr/dmExport"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

// TermProperty 物模型类型 属性
type TermProperty struct {
	AreaID           int64  `json:"areaID,string,omitempty"` //仅做记录
	ProductName      string `json:"productName,omitempty"`   //产品名称,只读
	ProductID        string `json:"productID,omitempty"`     //产品id
	DeviceName       string `json:"deviceName,omitempty"`
	DeviceAlias      string `json:"deviceAlias,omitempty"`
	SchemaAffordance string `json:"schemaAffordance,omitempty"` //只读,返回物模型定义
	Compare
}

func (c *TermProperty) Validate(repo CheckRepo) error {
	if c == nil {
		return nil
	}
	if repo.Info.DeviceMode == DeviceModeSingle {
		c.ProductID = repo.Info.ProductID
		c.DeviceName = repo.Info.DeviceName
	}
	if c.ProductID == "" {
		return errors.Parameter.AddMsg("执行条件设备类型中的产品id需要填写")
	}
	if c.DeviceName == "" {
		return errors.Parameter.AddMsg("执行条件设备类型中的设备名需要填写")
	}
	if len(c.DataID) == 0 {
		return errors.Parameter.AddMsg("执行条件设备类型中的标识符需要填写")
	}
	if repo.Info.DeviceMode != DeviceModeSingle {
		c.DeviceAlias = GetDeviceAlias(repo.Ctx, repo.DeviceCache, c.ProductID, c.DeviceName)
	}
	_, err := dmExport.SchemaAccess(repo.Ctx, repo.DeviceCache, repo.UserShareCache, def.AuthRead, devices.Core{
		ProductID:  c.ProductID,
		DeviceName: c.DeviceName,
	}, nil)
	if err != nil {
		return err
	}
	v, err := repo.SchemaCache.GetData(repo.Ctx, devices.Core{ProductID: c.ProductID, DeviceName: c.DeviceName})
	if err != nil {
		return err
	}
	dataID := strings.Split(c.DataID, ".")
	p := v.Property[dataID[0]]
	if p == nil {
		return errors.Parameter.AddMsg("dataID不存在")
	}
	if err := c.PropertyValidate(p); err != nil {
		return err
	}
	c.SchemaAffordance = schema.DoToAffordanceStr(p)
	if c.DataName == "" {
		c.DataName = p.Name
	}
	pi, err := repo.ProductCache.GetData(repo.Ctx, c.ProductID)
	if err != nil {
		return err
	}
	c.ProductName = pi.ProductName
	return nil
}
func (c *TermProperty) IsHit(ctx context.Context, columnType TermColumnType, repo CheckRepo, t *Term) bool {
	if c == nil {
		return false
	}
	dev := devices.Core{ProductID: c.ProductID, DeviceName: c.DeviceName}
	di, err := repo.DeviceCache.GetData(ctx, dev)
	if err != nil {
		if !errors.Cmp(err, errors.NotFind) {
			return false
		}
		t.IsAbnormal = true
		t.Reason = ReasonDeviceDelete
		repo.Info.IsAbnormal = true
		repo.Info.Reason = ReasonDeviceDelete
		return false
	}
	if di.ProjectID != repo.Info.ProjectID {
		t.IsAbnormal = true
		t.Reason = ReasonDeviceDelete
		repo.Info.IsAbnormal = true
		repo.Info.Reason = ReasonDeviceDelete
		return false
	}
	sm, err := repo.SchemaCache.GetData(ctx, dev)
	if err != nil {
		logx.WithContext(ctx).Errorf("%s.GetSchemaModel err:%v", utils.FuncName(), err)
		return false
	}
	switch columnType {
	case TermColumnTypeProperty:
		info, err := repo.DeviceMsg.PropertyLogLatestIndex(ctx, &devicemsg.PropertyLogLatestIndexReq{
			ProductID: c.ProductID, DeviceName: c.DeviceName, DataIDs: []string{c.DataID}})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.PropertyLatestIndex err:%v", utils.FuncName(), err)
			return false
		}
		if len(info.List) == 0 {
			logx.WithContext(ctx).Errorf("%s.PropertyLatestIndex err:dataID is not right:%v", utils.FuncName(), c.DataID[0])
			return false
		}
		if info.List[0].Timestamp != 0 { //如果有值
			dataIDs := strings.Split(c.DataID, ".")
			p := sm.Property[dataIDs[0]]
			if p == nil {
				return false
			}
			return func() bool {
			RUN:
				switch p.Define.Type {
				case schema.DataTypeStruct:
					var dataMap = map[string]any{}
					utils.Unmarshal([]byte(info.List[0].Value), &dataMap)
					return c.PropertyIsHit(p, c.DataID, dataMap)
				case schema.DataTypeArray:
					p.Define = *p.Define.ArrayInfo
					goto RUN
				default:
					return c.PropertyIsHit(p, c.DataID, info.List[0].Value)
				}
			}()
		}
		return false
	case TermColumnTypeEvent:
		logx.WithContext(ctx).Errorf("scene not support event yet")
		return false
	}
	return true
}
