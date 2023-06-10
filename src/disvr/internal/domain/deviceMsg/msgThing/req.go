package msgThing

import (
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"time"
)

type (
	Req struct {
		deviceMsg.CommonMsg
		Params      map[string]any `json:"params,omitempty"`      //参数列表
		Identifiers []string       `json:"identifiers,omitempty"` //内为希望设备上报的属性列表,不填为获取全部
		Version     string         `json:"version,omitempty"`     //协议版本，默认为1.0。
		EventID     string         `json:"eventId,omitempty"`     //事件的 Id，在数据模板事件中定义。
		ActionID    string         `json:"actionId,omitempty"`    //数据模板中的行为标识符，由开发者自行根据设备的应用场景定义
		Type        string         `json:"type,omitempty"`        //表示获取什么类型的信息（report:表示设备上报的信息 info:信息 alert:告警 fault:故障）
	}
	//设备基础信息
	DeviceBasicInfo struct {
		devices.Core
		Imei     string            `json:"imei,omitempty"`     //设备的 IMEI 号信息，非必填项
		Mac      string            `json:"mac,omitempty"`      //设备的 MAC 信息，非必填项
		Version  string            `json:"version,omitempty"`  //固件版本
		HardInfo string            `json:"hardInfo,omitempty"` //模组具体硬件型号
		SoftInfo string            `json:"softInfo,omitempty"` //模组软件版本
		Position *def.Point        `json:"position,omitempty"` //设备基础信息-坐标信息
		Tags     map[string]string `json:"tags,omitempty"`     //设备标签信息
	}
)

func (d Req) AddStatus(err error) Req {
	e := errors.Fmt(err)
	d.Code = e.Code
	d.Status = e.GetDetailMsg()
	return d
}

func (d *Req) GetTimeStamp(defaultTime int64) time.Time {
	if d.Timestamp == 0 {
		return time.UnixMilli(defaultTime)
	}
	return time.UnixMilli(d.Timestamp)
}

// 校验设备上报的参数合法性
func (d *Req) VerifyReqParam(t *schema.Model, tt schema.ParamType) (map[string]Param, error) {
	if len(d.Params) == 0 {
		return nil, errors.Parameter.AddDetail("need add params")
	}
	getParam := make(map[string]Param, len(d.Params))

	switch tt {
	case schema.ParamProperty:
		for k, v := range d.Params {
			p, ok := t.Property[k]
			if ok == false {
				continue
			}

			tp := Param{
				Identifier: p.Identifier,
				Name:       p.Name,
				Desc:       p.Desc,
				Mode:       p.Mode,
				Required:   p.Required,
			}

			err := tp.SetByDefine(&p.Define, v)
			if err == nil {
				getParam[k] = tp
			} else if !errors.Cmp(err, errors.NotFind) {
				return nil, errors.Fmt(err).AddDetail(p.Identifier)
			}
		}
	case schema.ParamEvent:
		p, ok := t.Event[d.EventID]
		if ok == false {
			return nil, errors.Parameter.AddDetail("need right eventId")
		}
		if d.Type != string(p.Type) {
			return nil, errors.Parameter.AddDetailf("err type:%v", d.Type)
		}

		for k, v := range p.Param {
			tp := Param{
				Identifier: v.Identifier,
				Name:       v.Name,
			}

			param, ok := d.Params[k]
			if ok == false {
				return nil, errors.Parameter.AddDetail("need param:" + k)
			}

			err := tp.SetByDefine(&v.Define, param)
			if err == nil {
				getParam[k] = tp
			} else if !errors.Cmp(err, errors.NotFind) {
				return nil, errors.Fmt(err).AddDetail(p.Identifier)
			}
		}
	case schema.ParamActionInput:
		p, ok := t.Action[d.ActionID]
		if ok == false {
			return nil, errors.Parameter.AddDetail("need right ActionID")
		}
		for k, v := range p.In {
			tp := Param{
				Identifier: v.Identifier,
				Name:       v.Name,
			}

			param, ok := d.Params[v.Identifier]
			if ok == false {
				return nil, errors.Parameter.AddDetail("need param:" + k)
			}

			err := tp.SetByDefine(&v.Define, param)
			if err == nil {
				getParam[k] = tp
			} else if !errors.Cmp(err, errors.NotFind) {
				return nil, errors.Fmt(err).AddDetail(p.Identifier)
			}
		}
	case schema.ParamActionOutput:
		p, ok := t.Action[d.ActionID]
		if ok == false {
			return nil, errors.Parameter.AddDetail("need right ActionID")
		}

		for k, v := range p.In {
			tp := Param{
				Identifier: v.Identifier,
				Name:       v.Name,
			}
			param, ok := d.Params[v.Identifier]
			if ok == false {
				return nil, errors.Parameter.AddDetail("need param:" + k)
			}

			err := tp.SetByDefine(&v.Define, param)
			if err == nil {
				getParam[k] = tp
			} else if !errors.Cmp(err, errors.NotFind) {
				return nil, errors.Fmt(err).AddDetail(p.Identifier)
			}
		}
	}
	return getParam, nil
}
