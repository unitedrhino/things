package otataskmanagelogic

import (
	"context"

	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaTaskIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskIndexLogic {
	return &OtaTaskIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OtaTaskIndexLogic) OtaTaskIndex(in *dm.OtaTaskIndexReq) (*dm.OtaTaskIndexResp, error) {
	var (
		info     []*dm.OtaTaskInfo
		size     int64
		page     int64 = 1
		pageSize int64 = 20
		err      error
		otDB     = relationDB.NewOtaTaskRepo(l.ctx)
	)
	if !(in.Page == nil || in.Page.Page == 0 || in.Page.Size == 0) {
		page = in.Page.Page
		pageSize = in.Page.Size
	}

	size, err = otDB.CountByFilter(
		l.ctx, relationDB.OtaTaskFilter{
			FirmwareID: in.FirmwareID,
		})
	if err != nil {
		return nil, err
	}
	di, err := otDB.FindByFilter(
		l.ctx, relationDB.OtaTaskFilter{
			FirmwareID: in.FirmwareID,
		}, &def.PageInfo{Size: pageSize, Page: page})
	if err != nil {
		return nil, err
	}
	info = make([]*dm.OtaTaskInfo, 0, len(di))
	for _, v := range di {
		info = append(info, ToOtaTaskInfo(v))
	}

	return &dm.OtaTaskIndexResp{List: info, Total: size}, nil
}
