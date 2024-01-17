package areamanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AreaInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AiDB *relationDB.AreaInfoRepo
}

func NewAreaInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreaInfoIndexLogic {
	return &AreaInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewAreaInfoRepo(ctx),
	}
}

// 获取区域信息列表
func (l *AreaInfoIndexLogic) AreaInfoIndex(in *sys.AreaInfoIndexReq) (*sys.AreaInfoIndexResp, error) {
	var (
		poArr []*relationDB.SysAreaInfo
		f     = relationDB.AreaInfoFilter{ProjectID: in.ProjectID, AreaIDs: in.AreaIDs, ParentAreaID: in.ParentAreaID}
	)

	poArr, err := l.AiDB.FindByFilter(l.ctx,
		f, logic.ToPageInfo(in.Page))
	if err != nil {
		l.Errorf("AreaInfoIndex find menu_info err,menuIds:%d,err:%v", in.AreaIDs, err)
		return nil, err
	}
	total, err := l.AiDB.CountByFilter(l.ctx, f)
	if err != nil {
		l.Errorf("AreaInfoIndex find menu_info err,menuIds:%d,err:%v", in.AreaIDs, err)
		return nil, err
	}

	return &sys.AreaInfoIndexResp{List: AreaInfosToPb(poArr), Total: total}, nil
}
