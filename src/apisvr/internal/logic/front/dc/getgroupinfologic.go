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

type GetGroupInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetGroupInfoLogic {
	return GetGroupInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

//todo 这里需要添加权限管理,只有组的成员才可以获取
func (l *GetGroupInfoLogic) GetGroupInfo(req types.GetGroupInfoReq) (*types.GetGroupInfoResp, error) {
	l.Infof("GetGroupInfo|req=%+v", req)
	dcReq, err := assemble.GetGroupInfoReqToRpc(&req)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.DcRpc.GetGroupInfo(l.ctx, dcReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetGroupInfo|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return assemble.GetGroupInfoRespToApi(resp)
}
