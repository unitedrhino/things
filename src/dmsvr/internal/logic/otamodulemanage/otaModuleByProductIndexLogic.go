package otamodulemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/logic"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaModuleByProductIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OmDB *relationDB.OtaModuleInfoRepo
}

func NewOtaModuleByProductIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaModuleByProductIndexLogic {
	return &OtaModuleByProductIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OmDB:   relationDB.NewOtaModuleInfoRepo(ctx),
	}
}

// 查询产品下的OTA模块列表
func (l *OtaModuleByProductIndexLogic) OtaModuleByProductIndex(in *dm.OTAModuleIndexReq) (*dm.OTAModuleIndexResp, error) {
	filter := relationDB.OtaModuleInfoFilter{
		ProductId:  in.ProductId,
		DeviceName: in.DeviceName,
	}
	total, err := l.OmDB.CountByFilter(l.ctx, filter)
	var otaModuleInfo []*dm.OtaModuleInfo
	otaModule, err := l.OmDB.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.PageInfo))
	if err != nil {
		l.Errorf("%s.ModuleInfo.OtaModuleInfo Delete failure err=%+v", utils.FuncName(), err)
		return nil, err
	}
	_ = copier.Copy(&otaModuleInfo, &otaModule)
	return &dm.OTAModuleIndexResp{OtaModuleInfoList: otaModuleInfo, Total: total}, nil
}
