package logic

import (
	"context"

	"github.com/i-Things/things/src/filesvr/internal/svc"
	"github.com/i-Things/things/src/filesvr/pb/file"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UploadLogic) Upload(stream file.File_UploadServer) error {
	// todo: add your logic here and delete this line
	stream.Recv()
	return nil
}
