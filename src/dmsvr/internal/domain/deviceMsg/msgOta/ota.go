package msgOta

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg"
)

const (
	TypeReport   = "report"
	TypeUpdate   = "update" //固件升级消息下行  返回升级信息，版本、固件地址
	TypeProgress = "progress"
)

type (
	Req struct {
		deviceMsg.CommonMsg
		Params params `json:"params,optional"`
	}
	Process struct {
		deviceMsg.CommonMsg
		Params processParams `json:"params,optional"`
	}
	params struct {
		ID      int64  `json:"id"`
		Version string `json:"version"`
		Module  string `json:"module"`
	}
	processParams struct {
		ID     int64  `json:"id"`
		Step   int64  `json:"step"`
		Desc   string `json:"desc"`
		Module string `json:"module"`
	}
)

func (d *Req) VerifyReqParam() error {
	if d.Params.Module == "" {
		return errors.Parameter.AddDetail("need add module")
	}

	return nil
}
func (d *Req) GetVersion() string {
	return d.Params.Version
}
func (d *Process) VerifyReqParam() error {
	if d.Params.Module == "" {
		return errors.Parameter.AddDetail("need add module")
	}
	if d.Params.Step == 0 {
		return errors.Parameter.AddDetail("need add Step")
	}
	return nil
}
