package deviceSend

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceTemplate"
	"time"
)

type (
	DeviceReq struct {
		Method      string                 `json:"method"`             //操作方法
		ClientToken string                 `json:"clientToken"`        //方便排查随机数
		Params      map[string]interface{} `json:"params,omitempty"`   //参数列表
		Version     string                 `json:"version,omitempty"`  //协议版本，默认为1.0。
		EventID     string                 `json:"eventId,omitempty"`  //事件的 Id，在数据模板事件中定义。
		ActionID    string                 `json:"actionId,omitempty"` //数据模板中的行为标识符，由开发者自行根据设备的应用场景定义
		Timestamp   int64                  `json:"timestamp,omitempty"`
		Showmeta    int64                  `json:"showmeta,omitempty"` //标识回复消息是否带 metadata，缺省为0表示不返回 metadata
		Type        string                 `json:"type,omitempty"`     //	表示获取什么类型的信息。report:表示设备上报的信息 info:信息 alert:告警 fault:故障
		Code        int64                  `json:"code,omitempty"`     //状态码
		Status      string                 `json:"status,omitempty"`   //返回信息
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

func (d *DeviceReq) VerifyReqParam(t *deviceTemplate.Template, tt deviceTemplate.TEMP_TYPE) (map[string]TempParam, error) {
	if len(d.Params) == 0 {
		return nil, errors.Parameter.AddDetail("need add params")
	}
	getParam := make(map[string]TempParam, len(d.Params))
	switch tt {
	case deviceTemplate.PROPERTY:
		for k, v := range d.Params {
			p, ok := t.Property[k]
			if ok == false {
				continue
			}
			tp := TempParam{
				ID:       p.ID,
				Name:     p.Name,
				Desc:     p.Desc,
				Mode:     p.Mode,
				Required: p.Required,
			}
			err := tp.AddDefine(&p.Define, v)
			if err == nil {
				getParam[k] = tp
			} else if !errors.Cmp(err, errors.NotFind) {
				return nil, errors.Fmt(err).AddDetail(p.ID)
			}
		}
	case deviceTemplate.EVENT:
		p, ok := t.Event[d.EventID]
		if ok == false {
			return nil, errors.Parameter.AddDetail("need right eventId")
		}
		if d.Type != p.Type {
			return nil, errors.Parameter.AddDetail("err type:" + d.Type)
		}

		for k, v := range p.Param {
			tp := TempParam{
				ID:   v.ID,
				Name: v.Name,
			}
			param, ok := d.Params[k]
			if ok == false {
				return nil, errors.Parameter.AddDetail("need param:" + k)
			}
			err := tp.AddDefine(&v.Define, param)
			if err == nil {
				getParam[k] = tp
			} else if !errors.Cmp(err, errors.NotFind) {
				return nil, errors.Fmt(err).AddDetail(p.ID)
			}
		}
	case deviceTemplate.ACTION_INPUT:
		p, ok := t.Action[d.ActionID]
		if ok == false {
			return nil, errors.Parameter.AddDetail("need right ActionID")
		}
		for k, v := range p.In {
			tp := TempParam{
				ID:   v.ID,
				Name: v.Name,
			}
			param, ok := d.Params[v.ID]
			if ok == false {
				return nil, errors.Parameter.AddDetail("need param:" + k)
			}
			err := tp.AddDefine(&v.Define, param)
			if err == nil {
				getParam[k] = tp
			} else if !errors.Cmp(err, errors.NotFind) {
				return nil, errors.Fmt(err).AddDetail(p.ID)
			}
		}
	case deviceTemplate.ACTION_OUTPUT:
		p, ok := t.Action[d.ActionID]
		if ok == false {
			return nil, errors.Parameter.AddDetail("need right ActionID")
		}
		for k, v := range p.In {
			tp := TempParam{
				ID:   v.ID,
				Name: v.Name,
			}
			param, ok := d.Params[v.ID]
			if ok == false {
				return nil, errors.Parameter.AddDetail("need param:" + k)
			}
			err := tp.AddDefine(&v.Define, param)
			if err == nil {
				getParam[k] = tp
			} else if !errors.Cmp(err, errors.NotFind) {
				return nil, errors.Fmt(err).AddDetail(p.ID)
			}
		}
	}
	return getParam, nil
}
