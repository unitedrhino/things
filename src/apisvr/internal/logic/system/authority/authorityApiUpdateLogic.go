package authority

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/domain/userHeader"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthorityApiUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthorityApiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthorityApiUpdateLogic {
	return &AuthorityApiUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthorityApiUpdateLogic) AuthorityApiUpdate(req *types.AuthorityApiUpdateReq) error {
	// clear old policies
	var oldPolicies [][]string
	RoleId := cast.ToString(userHeader.GetUserCtx(l.ctx).Role)
	oldPolicies = l.svcCtx.Casbin.GetFilteredPolicy(0, RoleId)
	if len(oldPolicies) != 0 {
		removeResult, err := l.svcCtx.Casbin.RemoveFilteredPolicy(0, RoleId)
		if err != nil {
			l.Errorf("%s.Casbin.GetFilteredPolicy req=%v err=%+v", utils.FuncName(), req, err)
			return err
		}
		if !removeResult {
			l.Errorf("%s.Casbin.GetFilteredPolicy req=%v", utils.FuncName(), req)
			return errors.System.AddDetail("RemoveFilteredPolicy Failed")
		}
	}

	// add new policies
	var policies [][]string
	for _, v := range req.List {
		policies = append(policies, []string{RoleId, v.Path, v.Method})
	}
	addResult, err := l.svcCtx.Casbin.AddPolicies(policies)
	if err != nil {
		err := errors.Fmt(err)
		l.Errorf("%s.Casbin.AddPolicies req=%v err=%+v", utils.FuncName(), req, err)
		return err
	}
	if !addResult {
		l.Errorf("%s Casbin.AddPolicies return nil req=%+v", utils.FuncName(), req)
		return errors.System.AddDetail("AddPolicies Failed")
	}
	return nil
}
