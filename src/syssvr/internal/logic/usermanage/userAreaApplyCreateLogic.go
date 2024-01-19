package usermanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAreaApplyCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserAreaApplyCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAreaApplyCreateLogic {
	return &UserAreaApplyCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserAreaApplyCreateLogic) UserAreaApplyCreate(in *sys.UserAreaApplyCreateReq) (*sys.Response, error) {
	_, err := relationDB.NewAreaInfoRepo(l.ctx).FindOne(l.ctx, in.AreaID, nil)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddMsgf("区域不存在")
		}
		return nil, err
	}
	err = relationDB.NewUserAreaApplyRepo(l.ctx).Insert(l.ctx, &relationDB.SysUserAreaApply{
		UserID:   ctxs.GetUserCtx(l.ctx).UserID,
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
