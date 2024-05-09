package scene

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
)

/*
	设备告警通知模版如下,需要在params中填写body,如: 电量过低
		你好,{{.deviceAlias}}设备告警:{{.body}}
*/

type ActionNotify struct {
	Type       def.NotifyType    `json:"type"`
	NotifyCode def.NotifyCode    `json:"notifyCode"` //支持: ruleScene: 场景联动通知(多设备模式使用)   ruleDeviceAlarm: 设备告警通知(单设备模式使用)
	Accounts   []string          `json:"accounts"`   //通知的账号列表
	Params     map[string]string `json:"params"`     //需要填写告警的
	Str1       string            `json:"str1"`
	Str2       string            `json:"str2"`
	Str3       string            `json:"str3"`
}

func (a *ActionNotify) Validate(repo ValidateRepo) error {
	if a == nil {
		return nil
	}
	if !utils.SliceIn(a.Type, def.NotifyTypeSms, def.NotifyTypeEmail,
		def.NotifyTypeDingTalk, def.NotifyTypeWx, def.NotifyTypeMessage, def.NotifyTypePhoneCall) {
		return errors.Parameter.AddMsg("消息通知不支持的类型:" + string(a.Type))
	}
	if repo.Info.DeviceMode == DeviceModeSingle {
		if a.Params == nil {
			a.Params = map[string]string{}
		}
		a.Params["productID"] = repo.Info.ProductID
		a.Params["deviceName"] = repo.Info.DeviceName
		a.Params["deviceAlias"] = repo.Info.DeviceAlias
	}
	return nil
}
func (a *ActionNotify) Execute(ctx context.Context, repo ActionRepo) error {
	_, err := repo.NotifyM.NotifyInfoSend(ctx, &sys.NotifyInfoSendReq{
		UserIDs:    []int64{repo.UserID},
		Accounts:   a.Accounts,
		NotifyCode: a.NotifyCode,
		Type:       a.Type,
		Params:     a.Params,
		Str1:       a.Str1,
		Str2:       a.Str2,
		Str3:       a.Str3,
	})
	return err
}
