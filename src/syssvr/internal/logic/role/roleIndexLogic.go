package rolelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	RiDB *relationDB.RoleInfoRepo
}

func NewRoleIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleIndexLogic {
	return &RoleIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		RiDB:   relationDB.NewRoleInfoRepo(ctx),
	}
}

func (l *RoleIndexLogic) RoleIndex(in *sys.RoleIndexReq) (*sys.RoleIndexResp, error) {
	f := relationDB.RoleInfoFilter{
		Name:         in.Name,
		Status:       in.Status,
		RoleInfoWith: &relationDB.RoleInfoWith{WithMenus: true},
	}
	ros, err := l.RiDB.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	total, err := l.RiDB.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	info := make([]*sys.RoleIndexData, 0, len(ros))
	for _, ro := range ros {
		var menuIDs []int64
		if len(ro.Menus) != 0 {
			for _, v := range ro.Menus {
				menuIDs = append(menuIDs, v.MenuID)
			}
		}
		info = append(info, &sys.RoleIndexData{
			Id:          ro.ID,
			Name:        ro.Name,
			Remark:      ro.Remark,
			CreatedTime: ro.CreatedTime.Unix(),
			Status:      ro.Status,
			RoleMenuID:  menuIDs,
		})
	}

	return &sys.RoleIndexResp{
		List:  info,
		Total: total,
	}, nil
}
