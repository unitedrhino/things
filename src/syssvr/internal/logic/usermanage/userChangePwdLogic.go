package usermanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserChangePwdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserChangePwdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserChangePwdLogic {
	return &UserChangePwdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserChangePwdLogic) UserChangePwd(in *sys.UserChangePwdReq) (*sys.Response, error) {
	var account string
	uc := ctxs.GetUserCtx(l.ctx)
	var oldUi *relationDB.SysTenantUserInfo
	switch in.Type {
	case def.CaptchaTypeEmail:
		account = l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeEmail, def.CaptchaUseChangePwd, in.CodeID, in.Code)
		if account == "" {
			return nil, errors.Captcha
		}
		ui, err := relationDB.NewUserInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Emails: []string{account}})
		if err != nil {
			return nil, err
		}
		oldUi = ui
	case def.CaptchaTypePhone:
		account = l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypePhone, def.CaptchaUseChangePwd, in.CodeID, in.Code)
		if account == "" {
			return nil, errors.Captcha
		}
		ui, err := relationDB.NewUserInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Phones: []string{account}})
		if err != nil {
			return nil, err
		}
		oldUi = ui
	}
	if oldUi.UserID != uc.UserID {
		return nil, errors.Permissions.AddMsgf("只能修改自己的密码")
	}
	err := CheckPwd(l.svcCtx, in.Password)
	if err != nil {
		return nil, err
	}
	oldUi.Password = utils.MakePwd(oldUi.Password, oldUi.UserID, false)
	err = relationDB.NewUserInfoRepo(l.ctx).Update(l.ctx, oldUi)
	if err != nil {
		return nil, err
	}
	return &sys.Response{}, nil
}
