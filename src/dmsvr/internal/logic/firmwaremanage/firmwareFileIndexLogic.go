package firmwaremanagelogic

import (
	"context"

	"gitee.com/i-Things/core/shared/def"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type FirmwareFileIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFirmwareFileIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FirmwareFileIndexLogic {
	return &FirmwareFileIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 附件列表搜索
func (l *FirmwareFileIndexLogic) FirmwareFileIndex(in *dm.OtaFirmwareFileIndexReq) (*dm.OtaFirmwareFileIndexResp, error) {
	l.Infof("GetFirmwareFileIndex|req=%+v", in)
	var (
		info     []*dm.OtaFirmwareFileInfo
		size     int64
		page     int64 = 1
		pageSize int64 = 20
		err      error
		piDB     = relationDB.NewOtaFirmwareFileRepo(l.ctx)
	)
	if !(in.Page == nil || in.Page.Page == 0 || in.Page.Size == 0) {
		page = in.Page.Page
		pageSize = in.Page.Size
	}
	var inSize *int64
	if in.Size != nil {
		tmpSize := in.Size.GetValue()
		inSize = &tmpSize
	}
	filter := relationDB.OtaFirmwareFileFilter{
		FirmwareID: in.FirmwareID,
		Size:       inSize,
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
	info = make([]*dm.OtaFirmwareFileInfo, 0, len(di))
	for _, v := range di {
		info = append(info, ToFirmwareFileResp(v))
	}
	return &dm.OtaFirmwareFileIndexResp{List: info, Total: size}, nil
}
