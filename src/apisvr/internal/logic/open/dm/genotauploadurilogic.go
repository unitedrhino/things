package dm

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/oss"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenOTAUploadUriLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenOTAUploadUriLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenOTAUploadUriLogic {
	return &GenOTAUploadUriLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenOTAUploadUriLogic) GenOTAUploadUri(req *types.GenOTAUploadUriReq) (resp *types.GenOTAUploadUriResp, err error) {
	l.Infof("[%s]|req=%+v", utils.FuncName(), utils.GetJson(req))
	dir, err := GenOssDir("ota")
	if err != nil {
		return nil, err
	}
	jwt, err := devices.GetJwtToken(l.svcCtx.Config.OSS.AccessSecret, time.Now().Unix(), l.svcCtx.Config.OSS.AccessExpire, oss.BucketOta, dir)
	if err != nil {
		return nil, err
	}
	return &types.GenOTAUploadUriResp{
		Sign: jwt,
		Host: oss.GetUploadUrl(l.svcCtx.Config.OSS.Minio),
	}, nil
}

func GenOssDir(business string) (string, error) {
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		logx.Info("GenerateUUID failure")
		return "", err
	}
	now := time.Now().Format("2006/01/02")
	return fmt.Sprintf("%s/%s/%s", business, now, uuid), nil
}
