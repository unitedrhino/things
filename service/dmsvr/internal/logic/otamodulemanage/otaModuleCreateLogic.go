package otamodulemanagelogic

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaModuleCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OmDB *relationDB.OtaModuleInfoRepo
}

func NewOtaModuleCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaModuleCreateLogic {
	return &OtaModuleCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OmDB:   relationDB.NewOtaModuleInfoRepo(ctx),
	}
}

// 创建产品的OTA模块
func (l *OtaModuleCreateLogic) OtaModuleCreate(in *dm.OTAModuleReq) (*dm.Response, error) {
	var otaModule relationDB.DmOtaModule
	_ = copier.Copy(&otaModule, &in)
	err := l.OmDB.Insert(l.ctx, &otaModule)
	if err != nil {
		l.Errorf("%s.ModuleInfo.OtaModuleInfo failure err=%+v", utils.FuncName(), err)
		return nil, err
	}
	return &dm.Response{}, nil
}
