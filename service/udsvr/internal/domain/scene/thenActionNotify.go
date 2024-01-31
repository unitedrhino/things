package scene

import (
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
)

type NotifyType string

const (
	NotifyTypeMsgCenter NotifyType = "msgCenter" //app上直接弹窗或者消息通知
)

type ActionNotify struct {
	Type NotifyType
}

func (a *ActionNotify) Validate() error {
	if a == nil {
		return nil
	}
	if !utils.SliceIn(a.Type, NotifyTypeMsgCenter) {
		return errors.Parameter.AddMsg("消息通知不支持的类型:" + string(a.Type))
	}
	return nil
}
