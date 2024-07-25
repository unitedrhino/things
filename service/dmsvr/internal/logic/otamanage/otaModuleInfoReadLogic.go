package otamanagelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaModuleInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaModuleInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaModuleInfoReadLogic {
	return &OtaModuleInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OtaModuleInfoReadLogic) OtaModuleInfoRead(in *dm.WithIDCode) (*dm.OtaModuleInfo, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	l.ctx = ctxs.WithRoot(l.ctx)
	po, err := relationDB.NewOtaModuleInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.OtaModuleInfoFilter{ID: in.Id, Code: in.Code})
	if err != nil {
		return nil, err
	}
	return utils.Copy[dm.OtaModuleInfo](po), nil
}
