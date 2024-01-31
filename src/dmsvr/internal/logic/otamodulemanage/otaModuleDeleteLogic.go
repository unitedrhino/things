package otamodulemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaModuleDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OmDB *relationDB.OtaModuleInfoRepo
}

func NewOtaModuleDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaModuleDeleteLogic {
	return &OtaModuleDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OmDB:   relationDB.NewOtaModuleInfoRepo(ctx),
	}
}

// 删除自定义OTA模块
func (l *OtaModuleDeleteLogic) OtaModuleDelete(in *dm.OTAModuleDeleteReq) (*dm.Response, error) {
	filter := relationDB.OtaModuleInfoFilter{
		ModuleName: in.ModuleName,
		ProductId:  in.ProductId,
	}
	err := l.OmDB.DeleteByFilter(l.ctx, filter)
	if err != nil {
		l.Errorf("%s.ModuleInfo.OtaModuleInfo Delete failure err=%+v", utils.FuncName(), err)
		return nil, err
	}
	return &dm.Response{}, nil
}
