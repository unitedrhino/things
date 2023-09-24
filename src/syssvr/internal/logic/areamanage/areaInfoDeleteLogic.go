package areamanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AreaInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AiDB *relationDB.AreaInfoRepo
}

func NewAreaInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreaInfoDeleteLogic {
	return &AreaInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewAreaInfoRepo(ctx),
	}
}

// 删除区域
func (l *AreaInfoDeleteLogic) AreaInfoDelete(in *sys.AreaInfoDeleteReq) (*sys.Response, error) {
	if in.AreaID == 0 {
		return nil, errors.Parameter
	}

	areaPo, err := checkArea(l.ctx, in.AreaID)
	if err != nil {
		return nil, errors.Database.AddDetail(err).WithMsg("检查区域出错")
	} else if areaPo == nil {
		return nil, errors.Parameter.AddDetail(in.AreaID).WithMsg("检查区域不存在")
	}

	areaAndChildIDs, err := l.AiDB.FindIDsWithChildren(l.ctx, in.AreaID)
	if err != nil {
		return nil, errors.Database.AddDetail(err).WithMsg("查询区域及子区域出错")
	}

	err = l.AiDB.DeleteByFilter(l.ctx, relationDB.AreaInfoFilter{AreaIDs: areaAndChildIDs})
	if err != nil {
		return nil, errors.Database.AddDetail(err).WithMsg("删除区域及子区域出错")
	}

	return &sys.Response{}, nil
}
