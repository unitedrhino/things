package msgSdkLog

import (
	"gitee.com/unitedrhino/things/share/domain/deviceMsg"
	"time"

	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
)

type (
	Req struct {
		deviceMsg.CommonMsg
		Params []SdkLog `json:"params,optional"` //参数列表
	}
	SdkLog struct {
		Content   string `json:"content"`
		Timestamp int64  `json:"timestamp,optional"`
		LogLevel  int64  `json:"logLevel,optional"`
	}
)

const (
	TypeOperation = "operation" //获取日志级别
	TypeReport    = "report"    //日志上报
	TypeUpdate    = "update"    //日志级别改变推送
)

func (d *Req) GetTimeStamp(logTime int64) time.Time {
	if logTime == 0 {
		return d.CommonMsg.GetTimeStamp()
	}
	return time.UnixMilli(logTime)
}

func (d *Req) VerifyReqParam() error {
	if len(d.Params) == 0 {
		return errors.Parameter.AddDetail("need add params")
	}
	for k, logObj := range d.Params {
		if logObj.Content == "" {
			return errors.Parameter.AddDetail("need param: content")
		}
		if logObj.LogLevel == 0 {
			d.Params[k].LogLevel = def.LogDebug
		}
	}

	return nil
}
