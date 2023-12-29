package usermanagelogic

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCaptchaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCaptchaLogic {
	return &UserCaptchaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserCaptchaLogic) UserCaptcha(in *sys.UserCaptchaReq) (*sys.UserCaptchaResp, error) {
	var (
		codeID = utils.Random(20, 1)
		code   = utils.Random(6, 0)
	)
	switch in.Type {
	case def.CaptchaTypeImage:

	case def.CaptchaTypeEmail:
		account := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeEmail, in.CodeID, in.Code)
		if account == "" {
			return nil, errors.Captcha
		}
		count, err := relationDB.NewUserInfoRepo(l.ctx).CountByFilter(l.ctx, relationDB.UserInfoFilter{Emails: []string{in.Account}})
		if err != nil {
			return nil, err
		}
		if count == 0 && in.Use == def.CaptchaUseLogin {
			return nil, errors.UnRegister
		}
		c, err := relationDB.NewTenantConfigRepo(l.ctx).FindOne(l.ctx)
		if err != nil {
			return nil, err
		}
		err = utils.SenEmail(conf.Email{
			From:     c.Email.From,
			Host:     c.Email.Host,
			Secret:   c.Email.Secret,
			Nickname: c.Email.Nickname,
			Port:     c.Email.Port,
			IsSSL:    c.Email.IsSSL == def.True,
		}, []string{in.Account}, "验证码校验",
			fmt.Sprintf("您的验证码为：%s，有效期为%d分钟", code, def.CaptchaExpire/60))
		if err != nil {
			return nil, err
		}
	}
	err := l.svcCtx.Captcha.Store(l.ctx, in.Type, codeID, code, in.Account, def.CaptchaExpire)
	if err != nil {
		return nil, err
	}
	l.Infof("code:%v codeID:%v", code, codeID)
	return &sys.UserCaptchaResp{Code: code, CodeID: codeID, Expire: def.CaptchaExpire}, nil
}
