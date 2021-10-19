package logic

import (
	"context"

	"gitee.com/godLei6/things/src/dcsvr/dc"
	"gitee.com/godLei6/things/src/dcsvr/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

type SendPropertyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendPropertyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendPropertyLogic {
	return &SendPropertyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 同步调用设备属性
func (l *SendPropertyLogic) SendProperty(in *dc.SendPropertyReq) (*dc.SendPropertyResp, error) {
	// todo: add your logic here and delete this line

	return &dc.SendPropertyResp{}, nil
}
