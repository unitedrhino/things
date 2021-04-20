package logic

import (
	"context"
	"yl/src/webapi/internal/svc"
	"yl/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type CaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) CaptchaLogic {
	return CaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CaptchaLogic) Captcha(req types.GetCaptchaReq) (*types.GetCaptchaResp, error) {
	// todo: add your logic here and delete this line
	l.Infof("Captcha|req=%#v",req)
	return &types.GetCaptchaResp{
		CodeID: "6666",
		Expire: 10000000000000,
	}, nil
}
