package otamodulemanagelogic

import (
	"context"
	"github.com/i-Things/things/src/dmsvr/internal/logic"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type OtaModuleVersionsByDeviceIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OmDB *relationDB.OtaModuleInfoRepo
}

func NewOtaModuleVersionsByDeviceIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaModuleVersionsByDeviceIndexLogic {
	return &OtaModuleVersionsByDeviceIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OmDB:   relationDB.NewOtaModuleInfoRepo(ctx),
	}
}

// 查询设备上报过的OTA模块版本列表
func (l *OtaModuleVersionsByDeviceIndexLogic) OtaModuleVersionsByDeviceIndex(in *dm.OTAModuleIndexReq) (*dm.OTAModuleVersionsIndexResp, error) {
	filter := relationDB.OtaModuleInfoFilter{
		ProductId:  in.ProductId,
		DeviceName: in.DeviceName,
	}
	total, err := l.OmDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}
	var otaModuleDetailList []*dm.OTAModuleDetail
	otaModuleList, err := l.OmDB.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.PageInfo))
	_ = copier.Copy(&otaModuleDetailList, &otaModuleList)
	return &dm.OTAModuleVersionsIndexResp{OtaModuleDetailList: otaModuleDetailList, Total: total}, nil
}
