package dc

import (
	"context"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/webapi/internal/dto"

	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type BgGetGroupMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBgGetGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) BgGetGroupMemberLogic {
	return BgGetGroupMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BgGetGroupMemberLogic) BgGetGroupMember(req types.GetGroupMemberReq) (*types.GetGroupMemberResp, error) {
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
