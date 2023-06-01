package rolelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/mysql"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleUpdateLogic {
	return &RoleUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleUpdateLogic) RoleUpdate(in *sys.RoleUpdateReq) (*sys.Response, error) {
	ro, err := l.svcCtx.RoleInfoModle.FindOne(l.ctx, in.Id)
	if err != nil {
		l.Logger.Error("RoleInfoModle.FindOne err , sql:%s", l.svcCtx)
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

	err = l.svcCtx.RoleInfoModle.Update(l.ctx, &mysql.SysRoleInfo{
		Id:     in.Id,
		Name:   in.Name,
		Remark: in.Remark,
		Status: in.Status,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &sys.Response{}, nil
}
