package areamanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AreaInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AiDB *relationDB.AreaInfoRepo
}

func NewAreaInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreaInfoReadLogic {
	return &AreaInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewAreaInfoRepo(ctx),
	}
}

// 获取区域信息详情
func (l *AreaInfoReadLogic) AreaInfoRead(in *sys.AreaWithID) (*sys.AreaInfo, error) {
	po, err := l.AiDB.FindOne(l.ctx, in.AreaID, nil)
	if err == nil {
		return transPoToPb(po), nil
	}
	if errors.Cmp(err, errors.NotFind) {
		return nil, errors.Parameter.AddDetail(in.AreaID).WithMsg("区域ID错误")
	}
	return nil, err
}
