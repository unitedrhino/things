package oss

import (
	"context"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownLoadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDownLoadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownLoadLogic {
	return &DownLoadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownLoadLogic) DownLoad(req *types.DownloadReq) error {
	// todo: add your logic here and delete this line

	return nil
}
