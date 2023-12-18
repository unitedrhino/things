package usermanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserForgetPwdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserForgetPwdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserForgetPwdLogic {
	return &UserForgetPwdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserForgetPwdLogic) UserForgetPwd(in *sys.UserForgetPwdReq) (*sys.Response, error) {
	var account string
	var oldUi *relationDB.SysUserInfo
	switch in.Type {
	case def.CaptchaTypeEmail:
		account = l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeEmail, in.CodeID, in.Code)
		if account == "" {
			return nil, errors.Captcha
		}
		ui, err := relationDB.NewUserInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Emails: []string{account}})
		if err != nil {
			return nil, err
		}
		oldUi = ui
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
