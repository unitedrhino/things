package deviceSend

import (
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"time"
)

type (
	DeviceReq struct {
		Method      string         `json:"method"`             //操作方法
		ClientToken string         `json:"clientToken"`        //方便排查随机数
		Params      map[string]any `json:"params,omitempty"`   //参数列表
		Version     string         `json:"version,omitempty"`  //协议版本，默认为1.0。
		EventID     string         `json:"eventId,omitempty"`  //事件的 Id，在数据模板事件中定义。
		ActionID    string         `json:"actionId,omitempty"` //数据模板中的行为标识符，由开发者自行根据设备的应用场景定义
		Timestamp   int64          `json:"timestamp,omitempty"`
		Showmeta    int64          `json:"showmeta,omitempty"` //标识回复消息是否带 metadata，缺省为0表示不返回 metadata
		Type        string         `json:"type,omitempty"`     //	表示获取什么类型的信息。report:表示设备上报的信息 info:信息 alert:告警 fault:故障
		Code        int64          `json:"code,omitempty"`     //状态码
		Status      string         `json:"status,omitempty"`   //返回信息
	}
)

func (d DeviceReq) AddStatus(err error) DeviceReq {
	e := errors.Fmt(err)
	d.Code = e.Code
	d.Status = e.GetDetailMsg()
	return d
}

func (d *DeviceReq) GetTimeStamp(defaultTime time.Time) time.Time {
	if d.Timestamp == 0 {
		return defaultTime
	}
	return time.UnixMilli(d.Timestamp)
}

func (d *DeviceReq) VerifyReqParam(t *schema.Model, tt schema.ParamType) (map[string]TempParam, error) {
	if len(d.Params) == 0 {
		return nil, errors.Parameter.AddDetail("need add params")
	}
	getParam := make(map[string]TempParam, len(d.Params))
	switch tt {
	case schema.ParamProperty:
		for k, v := range d.Params {
			p, ok := t.Property[k]
			if ok == false {
				continue
			}
			tp := TempParam{
				Identifier: p.Identifier,
				Name:       p.Name,
				Desc:       p.Desc,
				Mode:       p.Mode,
				Required:   p.Required,
			}
			err := tp.AddDefine(&p.Define, v)
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
			tp := TempParam{
				Identifier: v.Identifier,
				Name:       v.Name,
			}
			param, ok := d.Params[k]
			if ok == false {
				return nil, errors.Parameter.AddDetail("need param:" + k)
			}
			err := tp.AddDefine(&v.Define, param)
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
			tp := TempParam{
				Identifier: v.Identifier,
				Name:       v.Name,
			}
			param, ok := d.Params[v.Identifier]
			if ok == false {
				return nil, errors.Parameter.AddDetail("need param:" + k)
			}
			err := tp.AddDefine(&v.Define, param)
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
			tp := TempParam{
				Identifier: v.Identifier,
				Name:       v.Name,
			}
			param, ok := d.Params[v.Identifier]
			if ok == false {
				return nil, errors.Parameter.AddDetail("need param:" + k)
			}
			err := tp.AddDefine(&v.Define, param)
			if err == nil {
				getParam[k] = tp
			} else if !errors.Cmp(err, errors.NotFind) {
				return nil, errors.Fmt(err).AddDetail(p.Identifier)
			}
		}
	}
	return getParam, nil
}
