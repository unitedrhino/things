package logic

import (
	"context"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/shared/utils"

	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"

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
	l.Infof("Captcha|req=%+v",req)
	id,url,err:= l.svcCtx.Captcha.Get()
	if err !=nil {
		l.Errorf("%s|get Captcha err=%+v",utils.FuncName(),err)
		return  nil,errors.System.AddDetail(err.Error())
	}
	return &types.GetCaptchaResp{
		CodeID: id,
		Expire: l.svcCtx.Config.KeepTime,
		Url: url,
	}, nil
}
