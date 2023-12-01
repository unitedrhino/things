package otafirmwaremanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/dmsvr/internal/logic"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaFirmwareIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProductInfoRepo
	OfDB *relationDB.OtaFirmwareRepo
}

func NewOtaFirmwareIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaFirmwareIndexLogic {
	return &OtaFirmwareIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
		OfDB:   relationDB.NewOtaFirmwareRepo(ctx),
	}
}

// 升级包列表
func (l *OtaFirmwareIndexLogic) OtaFirmwareIndex(in *dm.OtaFirmwareIndexReq) (*dm.OtaFirmwareIndexResp, error) {
	var (
		info []*dm.OtaFirmwareInfo
		size int64
		err  error
	)
	filter := relationDB.OtaFirmwareFilter{
		ProductID:   in.ProductID,
		Module:      in.ModuleName,
		WithProduct: true,
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
	info = make([]*dm.OtaFirmwareInfo, len(list))
	err = copier.Copy(&info, &list)
	if err != nil {
		logx.Errorf("DmOtaFirmware copy to OtaFirmwareInfo failed")
		return nil, err
	}
	return &dm.OtaFirmwareIndexResp{Total: size, List: info}, nil
}
