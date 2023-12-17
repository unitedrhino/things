package user

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptchaLogic {
	return &CaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CaptchaLogic) Captcha(req *types.UserCaptchaReq) (resp *types.UserCaptchaResp, err error) {
	l.Infof("%s req=%v", utils.FuncName(), req)
	switch req.Type {
	case def.CaptchaTypeImage, "":
		id, url, answer, err := l.svcCtx.Captcha.Get()
		if err != nil {
			l.Errorf("%s get Captcha err=%+v", utils.FuncName(), err)
			return nil, errors.System.AddDetail(err)
		}
		l.Infof("%s id=%v,answer:%v", utils.FuncName(), id, answer)
		return &types.UserCaptchaResp{
			CodeID: id,
			Expire: l.svcCtx.Config.Captcha.KeepTime,
			Url:    url,
		}, nil
	case def.CaptchaTypeEmail:
		err := utils.SenEmail(l.svcCtx.Config.Email, []string{req.Account}, "邮箱验证码测试", "验证码为123123")
		return nil, err
	}
	return nil, errors.Parameter.AddDetail("type参数错误,暂不支持")
}
