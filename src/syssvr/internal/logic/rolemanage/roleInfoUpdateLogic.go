package rolemanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	RiDB *relationDB.RoleInfoRepo
}

func NewRoleInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleInfoUpdateLogic {
	return &RoleInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		RiDB:   relationDB.NewRoleInfoRepo(ctx),
	}
}

func (l *RoleInfoUpdateLogic) RoleInfoUpdate(in *sys.RoleInfo) (*sys.Response, error) {
	ro, err := l.RiDB.FindOne(l.ctx, in.Id)
	if err != nil {
		l.Logger.Error("RoleInfoModel.FindOne err , sql:%s", l.svcCtx)
		return nil, err
	}
	if in.Name == "" {
		in.Name = ro.Name
	}

	if in.Desc == "" {
		in.Desc = ro.Desc
	}

	if in.Status == 0 {
		in.Status = ro.Status
	}

	err = l.RiDB.Update(l.ctx, &relationDB.SysRoleInfo{
		ID:     in.Id,
		Name:   in.Name,
		Desc:   in.Desc,
		Status: in.Status,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &sys.Response{}, nil
}
