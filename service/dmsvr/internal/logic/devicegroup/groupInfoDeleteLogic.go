package devicegrouplogic

import (
	"context"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/things/service/dmsvr/internal/logic"
	"gitee.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/i-Things/things/service/dmsvr/internal/svc"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	GiDB *relationDB.GroupInfoRepo
}

func NewGroupInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInfoDeleteLogic {
	return &GroupInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		GiDB:   relationDB.NewGroupInfoRepo(ctx),
	}
}

// 删除分组
func (l *GroupInfoDeleteLogic) GroupInfoDelete(in *dm.WithID) (*dm.Empty, error) {
	po, err := relationDB.NewGroupInfoRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		if po.ParentID != 0 {
			c, err := relationDB.NewGroupInfoRepo(tx).CountByFilter(l.ctx, relationDB.GroupInfoFilter{ParentID: po.ParentID})
			if err != nil {
				return err
			}
			if c == 0 { //下面没有子节点了
				err = relationDB.NewGroupInfoRepo(tx).UpdateWithField(l.ctx,
					relationDB.GroupInfoFilter{ID: po.ParentID}, map[string]any{"is_leaf": def.True})
				if err != nil {
					return err
				}
			}
		}
		err := relationDB.NewGroupInfoRepo(l.ctx).Delete(l.ctx, in.Id)
		return err
	})
	logic.FillAreaGroupCount(l.ctx, l.svcCtx, int64(po.AreaID))
	return &dm.Empty{}, nil
}
