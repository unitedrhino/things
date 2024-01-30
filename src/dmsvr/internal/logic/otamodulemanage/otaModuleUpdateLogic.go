package otamodulemanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaModuleUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OmDB *relationDB.OtaModuleInfoRepo
}

func NewOtaModuleUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaModuleUpdateLogic {
	return &OtaModuleUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OmDB:   relationDB.NewOtaModuleInfoRepo(ctx),
	}
}

// 修改OTA模块别名、描述
func (l *OtaModuleUpdateLogic) OtaModuleUpdate(in *dm.OTAModuleReq) (*dm.Response, error) {
	var otaModule relationDB.DmOtaModule
	_ = copier.Copy(&otaModule, &in)
	err := l.OmDB.Update(l.ctx, &otaModule)
	if err != nil {
		l.Errorf("%s.ModuleInfo.OtaModuleInfo failure err=%+v", utils.FuncName(), err)
		return nil, err
	}
	return &dm.Response{}, nil
}
