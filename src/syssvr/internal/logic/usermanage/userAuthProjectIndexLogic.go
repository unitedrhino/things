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

type UserAuthProjectIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UapDB *relationDB.UserAuthProjectRepo
}

func NewUserAuthProjectIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAuthProjectIndexLogic {
	return &UserAuthProjectIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UapDB:  relationDB.NewUserAuthProjectRepo(ctx),
	}
}

func (l *UserAuthProjectIndexLogic) UserAuthProjectIndex(in *sys.UserProjectIndexReq) (*sys.UserProjectIndexResp, error) {
	var (
		list  []*sys.UserProject
		total int64
		err   error
	)
	if in.UserID == 0 {
		return nil, errors.Parameter.AddDetail(in.UserID).WithMsg("用户id必填")
	}

	filter := relationDB.UserAuthProjectFilter{
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

	list = make([]*sys.UserProject, 0, len(poArr))
	for _, po := range poArr {
		list = append(list, transProjectPoToPb(po))
	}
	return &sys.UserProjectIndexResp{List: list, Total: total}, nil
}
