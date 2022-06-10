package oss

import (
	"context"
	"github.com/minio/minio-go/v7"
	"net/http"

	"github.com/i-Things/things/src/filesvr/internal/svc"
	"github.com/i-Things/things/src/filesvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadLogic) Upload(req *types.UploadReq, r *http.Request) error {
	objName := req.Sign
	_, err := l.svcCtx.OSS.PutObject(l.ctx, "mymusic", objName,
		r.Body, r.ContentLength, minio.PutObjectOptions{ContentType: r.Header.Get("Content-Type")})
	if err != nil {
		return err
	}
	return nil
}
