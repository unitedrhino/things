package dc

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/webapi/internal/dto"

	"github.com/i-Things/things/src/webapi/internal/svc"
	"github.com/i-Things/things/src/webapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetGroupMemberLogic {
	return GetGroupMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

//todo 这里需要添加权限管理,只有组的成员才可以获取
func (l *GetGroupMemberLogic) GetGroupMember(req types.GetGroupMemberReq) (*types.GetGroupMemberResp, error) {
	l.Infof("GetGroupMember|req=%+v", req)
	dcReq, err := dto.GetGroupMemberReqToRpc(&req)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.DcRpc.GetGroupMember(l.ctx, dcReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetGroupMember|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return dto.GetGroupMemberRespToApi(resp)
}
