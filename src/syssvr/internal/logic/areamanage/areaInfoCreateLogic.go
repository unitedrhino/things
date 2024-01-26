package areamanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type AreaInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAreaInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreaInfoCreateLogic {
	return &AreaInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 新增区域
func (l *AreaInfoCreateLogic) AreaInfoCreate(in *sys.AreaInfo) (*sys.AreaWithID, error) {
	if in.AreaName == "" || in.ParentAreaID == 0 || ////root节点不为0
		in.ParentAreaID == def.NotClassified { //未分类不能有下属的区域
		return nil, errors.Parameter
	}
	if in.ProjectID == 0 {
		in.ProjectID = ctxs.GetUserCtx(l.ctx).ProjectID
	}
	var areaID = l.svcCtx.AreaID.GetSnowflakeId()
	var areaIDPath string = cast.ToString(areaID) + "-"
	var areaNamePath = in.AreaName + "-"
	areaPo := &relationDB.SysAreaInfo{
		AreaID:       stores.AreaID(areaID),
		ParentAreaID: in.ParentAreaID,                //创建时必填
		ProjectID:    stores.ProjectID(in.ProjectID), //创建时必填
		AreaIDPath:   areaIDPath,
		AreaNamePath: areaNamePath,
		AreaName:     in.AreaName,
		Position:     logic.ToStorePoint(in.Position),
		Desc:         utils.ToEmptyString(in.Desc),
	}
	conn := stores.GetTenantConn(l.ctx)
	err := conn.Transaction(func(tx *gorm.DB) error {
		projPo, err := checkProject(l.ctx, tx, in.ProjectID)
		if err != nil {
			return errors.Fmt(err).WithMsg("检查项目出错")
		} else if projPo == nil {
			return errors.Parameter.AddDetail(in.ProjectID).WithMsg("检查项目不存在")
		}
		aiRepo := relationDB.NewAreaInfoRepo(tx)
		if in.ParentAreaID != def.RootNode { //有选了父级项目区域
			pa, err := checkParentArea(l.ctx, tx, in.ParentAreaID)
			if err != nil {
				return err
			}
			areaPo.AreaIDPath = pa.AreaIDPath + cast.ToString(areaID) + "-"
			areaPo.AreaNamePath = pa.AreaNamePath + in.AreaName + "-"
			pa.LowerLevelCount++
			err = addSubAreaIDs(l.ctx, tx, pa, int64(areaPo.AreaID))
			if err != nil {
				return err
			}
			err = aiRepo.Update(l.ctx, pa)
			if err != nil {
				return err
			}
		}
		err = aiRepo.Insert(l.ctx, areaPo)
		if err != nil {
			l.Errorf("%s.Insert err=%+v", utils.FuncName(), err)
			return err
		}
		return nil
	})

	return &sys.AreaWithID{AreaID: int64(areaPo.AreaID)}, err
}
