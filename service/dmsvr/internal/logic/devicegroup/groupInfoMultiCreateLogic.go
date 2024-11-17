package devicegrouplogic

import (
	"context"
	"gitee.com/unitedrhino/share/errors"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInfoMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupInfoMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInfoMultiCreateLogic {
	return &GroupInfoMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupInfoMultiCreateLogic) GroupInfoMultiCreate(in *dm.GroupInfoMultiCreateReq) (*dm.Empty, error) {
	lg := NewGroupInfoCreateLogic(l.ctx, l.svcCtx)
	for _, g := range in.Groups {
		_, err := lg.GroupInfoCreate(g)
		if err != nil && !errors.Cmp(errors.Duplicate, err) {
			return nil, err
		}
	}
	return &dm.Empty{}, nil
}
