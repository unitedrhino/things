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
		Content  string `json:"content"`
		LogLevel int64  `json:"log_level,optional"`
	}
)

func (d *SdkLogReq) GetTimeStamp(defaultTime time.Time) time.Time {
	if d.Timestamp == 0 {
		return defaultTime
	}
	return time.UnixMilli(d.Timestamp)
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
