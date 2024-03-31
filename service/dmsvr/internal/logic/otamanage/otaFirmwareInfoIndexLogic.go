package otamanagelogic

import (
	"context"
	"gitee.com/i-Things/share/def"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaFirmwareInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProductInfoRepo
	OfDB *relationDB.OtaFirmwareInfoRepo
}

func NewOtaFirmwareInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaFirmwareInfoIndexLogic {
	return &OtaFirmwareInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
		OfDB:   relationDB.NewOtaFirmwareInfoRepo(ctx),
	}
}

// 升级包列表
func (l *OtaFirmwareInfoIndexLogic) OtaFirmwareInfoIndex(in *dm.OtaFirmwareInfoIndexReq) (*dm.OtaFirmwareInfoIndexResp, error) {
	var (
		info []*dm.OtaFirmwareInfo
		size int64
		err  error
	)
	filter := relationDB.OtaFirmwareInfoFilter{
		ProductID: in.ProductID,
		Name:      in.Name,
		WithFiles: true,
	}
	size, err = l.OfDB.CountByFilter(l.ctx, filter)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	list, err := l.OfDB.FindByFilter(l.ctx, filter, logic.ToPageInfoWithDefault(in.Page, logic.ToPageInfo(in.Page,
		def.OrderBy{Filed: "created_time", Sort: def.OrderDesc},
		def.OrderBy{Filed: "product_id", Sort: def.OrderDesc})))
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		info = append(info, ToFirmwareInfoPb(l.ctx, l.svcCtx, v))
	}
	return &dm.OtaFirmwareInfoIndexResp{Total: size, List: info}, nil
}