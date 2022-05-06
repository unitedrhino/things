package dc

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic/assemble"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ManageGroupInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewManageGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) ManageGroupInfoLogic {
	return ManageGroupInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

//todo 这里需要添加权限管理,只有组的管理员才可以编写
func (l *ManageGroupInfoLogic) ManageGroupInfo(req types.ManageGroupInfoReq) (*types.GroupInfo, error) {
	l.Infof("ManageGroupInfo|req=%+v", req)
	dcReq, err := assemble.ManageGroupInfoReqToRpc(&req)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.DcRpc.ManageGroupInfo(l.ctx, dcReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.ManageGroupInfo|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return assemble.GrouInfoToApi(resp), nil
}
