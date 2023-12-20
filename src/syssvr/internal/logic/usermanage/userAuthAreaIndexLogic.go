package usermanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAuthAreaIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UaaDB *relationDB.UserAuthAreaRepo
}

func NewUserAuthAreaIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAuthAreaIndexLogic {
	return &UserAuthAreaIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UaaDB:  relationDB.NewUserAuthAreaRepo(ctx),
	}
}

func (l *UserAuthAreaIndexLogic) UserAuthAreaIndex(in *sys.UserAreaIndexReq) (*sys.UserAreaIndexResp, error) {
	var (
		list  []*sys.UserArea
		total int64
		err   error
	)
	if in.UserID == 0 || in.ProjectID == 0 {
		return nil, errors.Parameter.AddDetail(in.UserID).WithMsg("用户id及项目id必填")
	}

	filter := relationDB.UserAuthAreaFilter{
		UserID:    in.UserID,
		ProjectID: in.ProjectID,
	}

	total, err = l.UaaDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}

	poArr, err := l.UaaDB.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}

	list = make([]*sys.UserArea, 0, len(poArr))
	for _, po := range poArr {
		list = append(list, transAreaPoToPb(po))
	}
	return &sys.UserAreaIndexResp{List: list, Total: total}, nil
}
