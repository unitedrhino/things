package role

import (
	"context"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleIndexLogic {
	return &RoleIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleIndexLogic) RoleIndex(req *types.RoleIndexReq) (resp *types.RoleIndexResp, err error) {
	var page sys.PageInfo
	copier.Copy(&page, req.Page)
	info, err := l.svcCtx.RoleRpc.RoleIndex(l.ctx, &sys.RoleIndexReq{
		Page:   &page,
		Name:   req.Name,
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}

	var roleInfo []*types.RoleIndexData
	var total int64
	total = info.Total

	roleInfo = make([]*types.RoleIndexData, 0, len(roleInfo))
	for _, i := range info.List {
		roleInfo = append(roleInfo, &types.RoleIndexData{
			ID:          i.Id,
			Name:        i.Name,
			Remark:      i.Remark,
			CreatedTime: i.CreatedTime,
			Status:      i.Status,
			RoleMenuID:  i.RoleMenuID,
		})
	}

	return &types.RoleIndexResp{roleInfo, total}, nil
}
