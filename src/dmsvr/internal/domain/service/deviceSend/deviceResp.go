package deviceSend

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/domain/thing"
	"time"
)

type (
	DeviceResp struct {
		Method      string                 `json:"method"`      //操作方法
		ClientToken string                 `json:"clientToken"` //方便排查随机数
		Timestamp   int64                  `json:"timestamp,omitempty"`
		Version     string                 `json:"version,omitempty"`  //协议版本，默认为1.0。
		Code        int64                  `json:"code"`               //状态码
		Status      string                 `json:"status,omitempty"`   //返回信息
		Type        string                 `json:"type,omitempty"`     //	表示什么类型的信息。report:表示设备上报的信息
		Data        map[string]interface{} `json:"data,omitempty"`     //返回具体设备上报的最新数据内容
		Response    map[string]interface{} `json:"response,omitempty"` //设备行为中定义的返回参数，设备行为执行成功后，向云端返回执行结果
	}
)

func (d DeviceResp) AddStatus(err error) DeviceResp {
	e := errors.Fmt(err)
	d.Code = e.Code
	d.Status = e.GetDetailMsg()
	return d
}

func (d *DeviceResp) GetTimeStamp(defaultTime time.Time) time.Time {
	if d.Timestamp == 0 {
		return defaultTime
	}
	return time.UnixMilli(d.Timestamp)
}

func (d *DeviceResp) VerifyRespParam(t *thing.Template, id string,
	tt thing.TEMP_TYPE) (map[string]TempParam, error) {
	getParam := make(map[string]TempParam, len(d.Response))
	switch tt {
	case thing.ACTION_OUTPUT:
		p, ok := t.Action[id]
		if ok == false {
			return nil, errors.Parameter.AddDetail("need right ActionID")
		}
		for k, v := range p.Out {
			tp := TempParam{
				ID:   v.ID,
				Name: v.Name,
			}
			param, ok := d.Response[v.ID]
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
