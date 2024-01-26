package areamanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"gorm.io/gorm"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AreaInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAreaInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreaInfoDeleteLogic {
	return &AreaInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除区域
func (l *AreaInfoDeleteLogic) AreaInfoDelete(in *sys.AreaWithID) (*sys.Response, error) {
	if in.AreaID == 0 {
		return nil, errors.Parameter
	}
	conn := stores.GetTenantConn(l.ctx)
	err := conn.Transaction(func(tx *gorm.DB) error {

		areaPo, err := checkArea(l.ctx, tx, in.AreaID)
		if err != nil {
			return errors.Fmt(err).WithMsg("检查区域出错")
		} else if areaPo == nil {
			return errors.Parameter.AddDetail(in.AreaID).WithMsg("检查区域不存在")
		}
		AiDB := relationDB.NewAreaInfoRepo(tx)

		if areaPo.ParentAreaID != def.RootNode { //如果有父节点
			parent, err := AiDB.FindOne(l.ctx, areaPo.ParentAreaID, nil)
			if err != nil {
				return err
			}
			parent.LowerLevelCount--
			subSubAreaIDs(l.ctx, tx, parent, in.AreaID)
			err = AiDB.Update(l.ctx, parent)
			if err != nil {
				return err
			}
		}
		areas, err := AiDB.FindByFilter(l.ctx, relationDB.AreaInfoFilter{AreaIDPath: areaPo.AreaIDPath}, nil)
		if err != nil {
			return errors.Fmt(err).WithMsg("查询区域及子区域出错")
		}
		var areaIDs []int64
		for _, area := range areas {
			areaIDs = append(areaIDs, int64(area.AreaID))
		}
		err = relationDB.NewAreaInfoRepo(tx).DeleteByFilter(l.ctx, relationDB.AreaInfoFilter{AreaIDs: areaIDs})
		if err != nil {
			return errors.Fmt(err).WithMsg("删除区域及子区域出错")
		}
		err = relationDB.NewUserAreaRepo(tx).DeleteByFilter(l.ctx, relationDB.UserAreaFilter{AreaIDs: areaIDs})
		if err != nil {
			return err
		}
		err = relationDB.NewUserAreaApplyRepo(tx).DeleteByFilter(l.ctx, relationDB.UserAreaApplyFilter{AreaIDs: areaIDs})
		if err != nil {
			return err
		}
		return nil
	})

	return &sys.Response{}, err
}
