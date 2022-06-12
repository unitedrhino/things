package dm

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/oss"
	"github.com/i-Things/things/shared/utils"
	"time"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenFirmwareUploadUriLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenFirmwareUploadUriLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenFirmwareUploadUriLogic {
	return &GenFirmwareUploadUriLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenFirmwareUploadUriLogic) GenFirmwareUploadUri(req *types.GenFirmwareUploadUriReq) (resp *types.GenFirmwareUploadUriResp, err error) {
	l.Infof("[%s]|req=%+v", utils.FuncName(), utils.GetJson(req))
	dir, err := GenOssDir("firmware")
	if err != nil {
		return nil, err
	}
	jwt, err := devices.GetJwtToken(l.svcCtx.Config.OSS.AccessSecret, time.Now().Unix(), l.svcCtx.Config.OSS.AccessExpire, oss.BucketFirmware, dir)
	if err != nil {
		return nil, err
	}
	return &types.GenFirmwareUploadUriResp{
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
