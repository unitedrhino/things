package oss

import (
	"context"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/errors"
	"github.com/minio/minio-go/v7"
	"net/http"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

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
	f, header, err := r.FormFile("file")
	if err != nil {
		return err
	}
	token, err := devices.ParseToken(req.Sign, l.svcCtx.Config.OSS.AccessSecret)
	if err != nil {
		return errors.Fmt(err)
	}
	userMetadata := map[string]string{
		"Filename": req.FileName,
	}
	_, err = l.svcCtx.OSS.PutObject(l.ctx, token.Bucket, token.Dir,
		f, header.Size, minio.PutObjectOptions{UserMetadata: userMetadata})
	if err != nil {
		return err
	}
	return nil
}
