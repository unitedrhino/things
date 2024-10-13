package shadow

import (
	"context"
	"gitee.com/unitedrhino/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/share/domain/schema"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"time"
)

const (
	ControlAuto             = 0 //自动,当设备不在线的时候设置设备影子,设备在线时直接下发给设备
	ControlNo               = 1 //只实时下发,不在线报错
	ControlOnly             = 2 //如果有设备影子只修改影子,没有的也不下发
	ControlOnlyCloud        = 3 //只修改云端的值
	ControlOnlyCloudWithLog = 4 //只修改云端的值并记录操作日志
	ControlNoLog            = 5 //只实时下发,不在线报错,且不记录日志
)
const (
	UpdatedDevice   = 1 //已经更新到过设备
	NotUpdateDevice = 2 //尚未更新到过设备
)

type (
	Info struct {
		ID                int64
		ProductID         string     // 产品id
		DeviceName        string     // 设备名称
		DataID            string     // 属性id
		Value             string     // 属性值
		UpdatedDeviceTime *time.Time //更新到设备的时间
		CreatedTime       time.Time
		UpdatedTime       time.Time
	}
	Filter struct {
		ProductID           string
		DeviceName          string
		DataIDs             []string
		UpdatedDeviceStatus int64 //1 已经更新到过设备 2 尚未更新到过设备
	}
	Repo interface {
		FindByFilter(ctx context.Context, f Filter) ([]*Info, error)
		// MultiUpdate 批量更新 LightStrategyDevice 记录
		MultiUpdate(ctx context.Context, data []*Info) error
		// MultiDelete 批量删除 LightStrategyDevice 记录
		MultiDelete(ctx context.Context, f Filter) error
	}
)

func CheckEnableShadow(params map[string]any, model *schema.Model) error {
	for k := range params {
		kk, _, _ := schema.GetArray(k)
		if prop, ok := model.Property[kk]; !ok {
			return errors.Parameter.AddMsgf("属性: %v 未定义该物模型属性", k)
		} else if prop.IsUseShadow != true {
			return errors.Parameter.AddMsgf("属性: %v 未开启设备影子模式", k)
		}
	}
	return nil
}
func NewInfo(productID, deviceName string, params map[string]any) []*Info {
	var ret []*Info
	for k, v := range params {
		ret = append(ret, &Info{
			ID:         0,
			ProductID:  productID,
			DeviceName: deviceName,
			DataID:     k,
			Value:      utils.MarshalNoErr(v),
		})
	}
	return ret
}

func ToValues(in []*Info, property schema.PropertyMap) map[string]msgThing.Param {
	var ret = map[string]msgThing.Param{}
	for _, v := range in {
		p := property[v.DataID]
		if p == nil {
			continue
		}
		val, err := p.Define.FmtValue(v.Value)
		if err != nil {
			continue
		}
		ret[v.DataID] = msgThing.Param{
			Identifier: p.Identifier,
			Name:       p.Name,
			Desc:       p.Desc,
			Mode:       p.Mode,
			Required:   p.Required,
			Define:     &p.Define,
			Value:      val,
		}
	}
	return ret
}

func ToMap(in []*Info) map[string]*Info {
	var ret = map[string]*Info{}
	for _, v := range in {
		ret[v.DataID] = v
	}
	return ret
}
