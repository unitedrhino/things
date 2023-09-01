package rolelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/shared/utils/cast"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleApiMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleApiMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleApiMultiUpdateLogic {
	return &RoleApiMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleApiMultiUpdateLogic) RoleApiMultiUpdate(in *sys.RoleApiMultiUpdateReq) (*sys.Response, error) {
	// clear old policies
	var oldPolicies [][]string
	oldPolicies = l.svcCtx.Casbin.GetFilteredPolicy(0, cast.ToString(in.RoleID))
	if len(oldPolicies) != 0 {
		removeResult, err := l.svcCtx.Casbin.RemoveFilteredPolicy(0, cast.ToString(in.RoleID))
		if err != nil {
			l.Errorf("%s.Casbin.RemoveFilteredPolicy req=%v err=%+v", utils.FuncName(), in, err)
			return nil, errors.Permissions.AddDetail(err)
		}
		if !removeResult {
			l.Errorf("%s.Casbin.RemoveFilteredPolicy req=%v", utils.FuncName(), in)
			return nil, errors.System.AddDetail("RemoveFilteredPolicy Failed")
		}
	}

	// add new policies
	var policies [][]string
	for _, v := range in.List {
		policies = append(policies, []string{cast.ToString(in.RoleID), v.Route, cast.ToString(v.Method)})
	}
	addResult, err := l.svcCtx.Casbin.AddPolicies(policies)
	if err != nil {
		err := errors.Fmt(err)
		l.Errorf("%s.Casbin.AddPolicies req=%v err=%+v", utils.FuncName(), in, err)
		return nil, errors.Permissions.AddDetail(err)
	}
	if !addResult {
		l.Errorf("%s Casbin.AddPolicies return nil req=%+v", utils.FuncName(), in)
		return nil, errors.System.AddDetail(err)
	}

	return &sys.Response{}, nil
}
