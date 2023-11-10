package firmwaremanagelogic

import (
	"context"

	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type FirmwareFileUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFirmwareFileUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FirmwareFileUpdateLogic {
	return &FirmwareFileUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 附件信息更新
func (l *FirmwareFileUpdateLogic) FirmwareFileUpdate(in *dm.OtaFirmwareFileReq) (*dm.OtaFirmwareFileResp, error) {
	db := relationDB.NewOtaFirmwareFileRepo(l.ctx)
	df, err := db.FindOneByFilter(l.ctx, relationDB.OtaFirmwareFileFilter{
		ID: in.FileID,
	})
	if err != nil {
		return nil, err
	}

	df.Size = in.Size
	df.Signature = in.Signature
	err = db.Update(l.ctx, df)
	if err != nil {
		return nil, err
	}
	return &dm.OtaFirmwareFileResp{}, nil
}
