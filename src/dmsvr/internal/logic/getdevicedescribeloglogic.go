package logic

import (
	"context"

	"github.com/go-things/things/src/dmsvr/dm"
	"github.com/go-things/things/src/dmsvr/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeviceDescribeLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeviceDescribeLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeviceDescribeLogLogic {
	return &GetDeviceDescribeLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备调试信息记录登入登出,操作
func (l *GetDeviceDescribeLogLogic) GetDeviceDescribeLog(in *dm.GetDeviceDescribeLogReq) (*dm.GetDeviceDescribeLogResp, error) {
	// todo: add your logic here and delete this line

	return &dm.GetDeviceDescribeLogResp{}, nil
}
