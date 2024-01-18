package usermanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAreaApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserAreaApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAreaApplyLogic {
	return &UserAreaApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserAreaApplyLogic) UserAreaApply(in *sys.UserAreaApplyReq) (*sys.Response, error) {
	err := relationDB.NewUserAreaApplyRepo(l.ctx).Insert(l.ctx, &relationDB.SysUserAreaApply{
		AreaID:   stores.AreaID(in.AreaID),
		AuthType: in.AuthType,
	})
	if err != nil {
		if errors.Cmp(err, errors.Duplicate) {
			return &sys.Response{}, nil
		}
		return nil, err
	}
	return &sys.Response{}, nil
}
