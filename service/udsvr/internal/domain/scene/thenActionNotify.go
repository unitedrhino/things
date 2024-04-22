package scene

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
)

type ActionNotify struct {
	Type       def.NotifyType    `json:"type"`
	NotifyCode string            `json:"notifyCode"`
	Accounts   []string          `json:"accounts"` //通知的账号列表
	Params     map[string]string `json:"params"`
	Str1       string            `json:"str1"`
	Str2       string            `json:"str2"`
	Str3       string            `json:"str3"`
}

func (a *ActionNotify) Validate(repo ValidateRepo) error {
	if a == nil {
		return nil
	}
	if !utils.SliceIn(a.Type, def.NotifyTypeSms, def.NotifyTypeEmail,
		def.NotifyTypeDingTalk, def.NotifyTypeWx, def.NotifyTypeMessage) {
		return errors.Parameter.AddMsg("消息通知不支持的类型:" + string(a.Type))
	}
	return nil
}
func (a *ActionNotify) Execute(ctx context.Context, repo ActionRepo) error {
	repo.NotifyM.NotifyInfoSend(ctx, &sys.NotifyInfoSendReq{
		UserIDs:    []int64{repo.UserID},
		Accounts:   a.Accounts,
		NotifyCode: a.NotifyCode,
		Type:       a.Type,
		Params:     a.Params,
		Str1:       a.Str1,
		Str2:       a.Str2,
		Str3:       a.Str3,
	})
	return nil
}
