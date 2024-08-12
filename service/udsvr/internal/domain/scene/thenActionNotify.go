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

type NotifyUserType = string

const (
	NotifyUserAccount        NotifyUserType = "account"
	NotifyUserID                            = "userID"
	NotifyDeviceOwner                       = "deviceOwner"        //设备拥有人
	NotifyDeviceProjectAdmin                = "deviceProjectAdmin" //拥有设备项目权限的管理员
	NotifyDeviceAreaAdmin                   = "deviceAreaAdmin"    //拥有设备区域权限的管理员
	NotifyDeviceProjectAll                  = "deviceProjectAll"   //拥有设备项目权限的管理员
	NotifyDeviceAreaAll                     = "deviceAreaAll"      //拥有设备区域权限的管理员
)

type ActionNotify struct {
	Type       def.NotifyType    `json:"type"`
	NotifyCode def.NotifyCode    `json:"notifyCode"` //支持: ruleScene: 场景联动通知(多设备模式使用)   ruleDeviceAlarm: 设备告警通知(单设备模式使用)
	UserType   NotifyUserType    `json:"userType"`
	Accounts   []string          `json:"accounts,omitempty"` //通知的账号列表
	UserIDs    []string          `json:"userIDs,omitempty"`
	Params     map[string]string `json:"params,omitempty"` //需要填写告警的
	Str1       string            `json:"str1,omitempty"`
	Str2       string            `json:"str2,omitempty"`
	Str3       string            `json:"str3,omitempty"`
}

func (a *ActionNotify) Validate(repo CheckRepo) error {
	if a == nil {
		return nil
	}
	if !utils.SliceIn(a.Type, def.NotifyTypeSms, def.NotifyTypeEmail,
		def.NotifyTypeDingTalk, def.NotifyTypeWxMini, def.NotifyTypeMessage, def.NotifyTypePhoneCall, def.NotifyTypeWxEWebhook) {
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
	_, err := repo.NotifyM.NotifyConfigSend(ctx, &sys.NotifyConfigSendReq{
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
