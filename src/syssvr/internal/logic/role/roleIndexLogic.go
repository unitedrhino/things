package rolelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleIndexLogic {
	return &RoleIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleIndexLogic) RoleIndex(in *sys.RoleIndexReq) (*sys.RoleIndexResp, error) {
	ros, total, err := l.svcCtx.RoleModel.Index(in)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	info := make([]*sys.RoleIndexData, 0, len(ros))
	for _, ro := range ros {
		info = append(info, &sys.RoleIndexData{
			Id:          ro.Id,
			Name:        ro.Name,
			Remark:      ro.Remark,
			CreatedTime: ro.CreatedTime.Unix(),
			Status:      ro.Status,
		})
	}

	for i, v := range info {
		MmuIDs, err := l.svcCtx.RoleModel.IndexRoleIDMenuID(v.Id)
		if err != nil {
			info[i].RoleMenuID = nil
			continue
		}
		info[i].RoleMenuID = MmuIDs
	}

	return &sys.RoleIndexResp{
		List:  info,
		Total: total,
	}, nil

	return &sys.RoleIndexResp{}, nil
}
