package oss

import (
	"context"
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
	//{
	//	objName := req.Sign
	//	info, err := l.svcCtx.OSS.PutObject(l.ctx, "mymusic", objName,
	//		r.Body, r.ContentLength, minio.PutObjectOptions{ContentType: r.Header.Get("Content-Type")})
	//	fmt.Println(info, err)
	//	err = l.svcCtx.OSS.FGetObject(l.ctx, "mymusic", objName, "./"+objName, minio.GetObjectOptions{})
	//	fmt.Println(err)
	//}
	return nil
}
