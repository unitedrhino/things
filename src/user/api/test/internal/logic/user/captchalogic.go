package logic

import (
	"context"

	"yl/src/user/api/test/internal/svc"
	"yl/src/user/api/test/internal/types"

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

	return &types.GetCaptchaResp{}, nil
}
