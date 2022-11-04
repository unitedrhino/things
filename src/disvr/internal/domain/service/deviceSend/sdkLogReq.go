package deviceSend

import (
	"time"

	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
)

type (
	SdkLogReq struct {
		Method      string   `json:"method"`          //操作方法
		ClientToken string   `json:"clientToken"`     //方便排查随机数
		Params      []sdklog `json:"params,optional"` //参数列表
		Timestamp   int64    `json:"timestamp,omitempty"`
	}
	sdklog struct {
		Content   string `json:"content"`
		Timestamp int64  `json:"timestamp,optional"`
		LogLevel  int64  `json:"logLevel,optional"`
	}
)

func (d *SdkLogReq) GetTimeStamp(logTime int64) time.Time {
	if logTime == 0 {
		if d.Timestamp != 0 {
			return time.UnixMilli(d.Timestamp)
		}
		return time.Now()
	}
	return time.UnixMilli(logTime)
}

func (d *SdkLogReq) VerifyReqParam() error {
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
