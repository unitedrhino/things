package otamanagelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/domain/deviceMsg/msgOta"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaModuleInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaModuleInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaModuleInfoCreateLogic {
	return &OtaModuleInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OtaModuleInfoCreateLogic) OtaModuleInfoCreate(in *dm.OtaModuleInfo) (*dm.WithID, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	l.ctx = ctxs.WithRoot(l.ctx)
	if in.Code == msgOta.ModuleCodeDefault || in.Code == "" {
		return nil, errors.Parameter.AddMsg("编码不能为空及不能为default")
	}
	_, err := l.svcCtx.ProductCache.GetData(l.ctx, in.ProductID)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddMsg("product not find")
		}
		return nil, err
	}
	po := utils.Copy[relationDB.DmOtaModuleInfo](in)
	in.Id = 0
	err = relationDB.NewOtaModuleInfoRepo(l.ctx).Insert(l.ctx, po)
	if err != nil {
		return nil, err
	}
	return &dm.WithID{Id: po.ID}, nil
}
