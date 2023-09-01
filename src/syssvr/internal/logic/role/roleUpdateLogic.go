package rolelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	RiDB *relationDB.RoleInfoRepo
}

func NewRoleUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleUpdateLogic {
	return &RoleUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		RiDB:   relationDB.NewRoleInfoRepo(ctx),
	}
}

func (l *RoleUpdateLogic) RoleUpdate(in *sys.RoleUpdateReq) (*sys.Response, error) {
	ro, err := l.RiDB.FindOne(l.ctx, in.Id, nil)
	if err != nil {
		l.Logger.Error("RoleInfoModel.FindOne err , sql:%s", l.svcCtx)
		return nil, err
	}
	if in.Name == "" {
		in.Name = ro.Name
	}

	if in.Remark == "" {
		in.Remark = ro.Remark
	}

	if in.Status == 0 {
		in.Status = ro.Status
	}

	err = l.RiDB.Update(l.ctx, &relationDB.SysRoleInfo{
		ID:     in.Id,
		Name:   in.Name,
		Remark: in.Remark,
		Status: in.Status,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &sys.Response{}, nil
}
