package msgThing

import (
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"time"
)

type (
	Resp struct {
		deviceMsg.CommonMsg
		Version  string         `json:"version,omitempty"`  //协议版本，默认为1.0。
		Type     string         `json:"type,omitempty"`     //	表示什么类型的信息。report:表示设备上报的信息
		Response map[string]any `json:"response,omitempty"` //设备行为中定义的返回参数，设备行为执行成功后，向云端返回执行结果
		ActionID string         `json:"actionID,omitempty"` //数据模板中的行为标识符，由开发者自行根据设备的应用场景定义
	}
)

func (d *Resp) GetTimeStamp(defaultTime time.Time) time.Time {
	if d.Timestamp == 0 {
		return defaultTime
	}
	return time.UnixMilli(d.Timestamp)
}

func (d *Resp) FmtRespParam(t *schema.Model, id string, tt schema.ParamType) error {
	param, err := d.VerifyRespParam(t, id, tt)
	if err != nil {
		return err
	}
	d.Response, err = ToVal(param)
	if err != nil {
		return err
	}
	return nil
}

func (d *Resp) VerifyRespParam(t *schema.Model, id string,
	tt schema.ParamType) (map[string]Param, error) {
	getParam := make(map[string]Param, len(d.Response))
	switch tt {
	case schema.ParamActionOutput:
		p, ok := t.Action[id]
		if ok == false {
			return nil, errors.Parameter.AddDetail("need right ActionID")
		}
		for k, v := range p.Out {
			tp := Param{
				Identifier: v.Identifier,
				Name:       v.Name,
			}
			param, ok := d.Response[v.Identifier]
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
