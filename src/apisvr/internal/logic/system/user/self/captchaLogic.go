package self

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/pb/sys"

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
	ret, err := l.svcCtx.UserRpc.UserCaptcha(l.ctx, &sys.UserCaptchaReq{Account: req.Account, Type: req.Type, Use: req.Use})
	if err != nil {
		l.Errorf("%s UserCaptcha err=%+v", utils.FuncName(), err)
		return nil, err
	}
	switch req.Type {
	case def.CaptchaTypeImage:
		url := l.svcCtx.Captcha.GetB64(ret.Code)
		return &types.UserCaptchaResp{
			CodeID: ret.CodeID,
			Expire: ret.Expire,
			Url:    url,
		}, nil
	default:
		return &types.UserCaptchaResp{
			CodeID: ret.CodeID,
			Expire: ret.Expire,
		}, nil
	}
}
