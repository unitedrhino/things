package firmwaremanagelogic

import (
	"context"

	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type FirmwareInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFirmwareInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FirmwareInfoIndexLogic {
	return &FirmwareInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FirmwareInfoIndexLogic) FirmwareInfoIndex(in *dm.FirmwareInfoIndexReq) (*dm.FirmwareInfoIndexResp, error) {
	l.Infof("GetFirmwareInfoIndex|req=%+v", in)
	var (
		info     []*dm.FirmwareInfo
		size     int64
		page     int64 = 1
		pageSize int64 = 20
		err      error
		piDB     = relationDB.NewOtaFirmwareRepo(l.ctx)
	)
	if !(in.Page == nil || in.Page.Page == 0 || in.Page.Size == 0) {
		page = in.Page.Page
		pageSize = in.Page.Size
	}

	filter := relationDB.OtaFirmwareFilter{
		ProductID:  in.ProductID,
		FirmwareID: in.FirmwareID,
	}
	size, err = piDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}
	di, err := piDB.FindByFilter(
		l.ctx, filter, &def.PageInfo{Size: pageSize, Page: page})
	if err != nil {
		return nil, err
	}
	info = make([]*dm.FirmwareInfo, 0, len(di))
	for _, v := range di {
		info = append(info, ToFirmwareInfo(v))
	}
	return &dm.FirmwareInfoIndexResp{List: info, Total: size}, nil
}
