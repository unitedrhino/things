package otamanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

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
