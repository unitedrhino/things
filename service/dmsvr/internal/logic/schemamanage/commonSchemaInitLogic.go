package schemamanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommonSchemaInitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommonSchemaInitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommonSchemaInitLogic {
	return &CommonSchemaInitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CommonSchemaInitLogic) CommonSchemaInit(in *dm.Empty) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	pos, err := relationDB.NewCommonSchemaRepo(l.ctx).FindByFilter(l.ctx, relationDB.CommonSchemaFilter{Type: 1}, nil)
	if err != nil {
		return nil, err
	}
	for _, po := range pos {
		if err := l.svcCtx.SchemaManaRepo.CreateProperty(l.ctx, relationDB.ToPropertyDo(&po.DmSchemaCore), ""); err != nil {
			l.Errorf("%s.CreateProperty failure,err:%v", utils.FuncName(), err)
			continue
		}
	}
	return &dm.Empty{}, nil
}
