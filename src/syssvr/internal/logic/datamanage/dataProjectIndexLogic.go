package datamanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataProjectIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UapDB *relationDB.DataProjectRepo
}

func NewDataProjectIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataProjectIndexLogic {
	return &DataProjectIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UapDB:  relationDB.NewDataProjectRepo(ctx),
	}
}

func (l *DataProjectIndexLogic) DataProjectIndex(in *sys.DataProjectIndexReq) (*sys.DataProjectIndexResp, error) {
	var (
		list  []*sys.DataProject
		total int64
		err   error
	)

	filter := relationDB.UserProjectFilter{
		UserID: in.UserID,
	}

	total, err = l.UapDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}

	poArr, err := l.UapDB.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}

	list = make([]*sys.DataProject, 0, len(poArr))
	for _, po := range poArr {
		list = append(list, transProjectPoToPb(po))
	}
	return &sys.DataProjectIndexResp{List: list, Total: total}, nil
}
