package areamanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type AreaInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AiDB *relationDB.AreaInfoRepo
}

func NewAreaInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreaInfoUpdateLogic {
	return &AreaInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewAreaInfoRepo(ctx),
	}
}

// 更新区域
func (l *AreaInfoUpdateLogic) AreaInfoUpdate(in *sys.AreaInfo) (*sys.Response, error) {
	if in.AreaID == 0 || utils.SliceIn(in.AreaID, def.RootNode, def.NotClassified) {
		return nil, errors.Parameter
	}

	areaPo, err := checkArea(l.ctx, in.AreaID)
	if err != nil {
		return nil, errors.Fmt(err).WithMsg("检查区域出错")
	} else if areaPo == nil {
		return nil, errors.Parameter.AddDetail(in.AreaID).WithMsg("检查区域不存在")
	}

	if in.ParentAreaID != 0 { //有选了父级项目区域
		if _, err := checkParentArea(l.ctx, in.ParentAreaID, false); err != nil {
			return nil, err
		}
	}

	projPo, err := checkProject(l.ctx, in.ProjectID)
	if err != nil {
		return nil, errors.Fmt(err).WithMsg("检查项目出错")
	} else if projPo == nil {
		return nil, errors.Parameter.AddDetail(in.ProjectID).WithMsg("检查项目不存在")
	}

	l.setPoByPb(areaPo, in)

	err = l.AiDB.Update(l.ctx, areaPo)
	if err != nil {
		return nil, errors.Fmt(err).WithMsg("检查出错")
	}

	return &sys.Response{}, nil
}
func (l *AreaInfoUpdateLogic) setPoByPb(po *relationDB.SysAreaInfo, pb *sys.AreaInfo) {
	//不支持更改 区域所属项目，因此不进行赋值

	//支持区域 改为 第一级区域（改字段前端必填）
	po.ParentAreaID = pb.ParentAreaID

	if pb.AreaName != "" {
		po.AreaName = pb.AreaName
	}
	if pb.Position != nil {
		po.Position = logic.ToStorePoint(pb.Position)
	}
	if pb.Desc != nil {
		po.Desc = pb.Desc.GetValue()
	}
}
