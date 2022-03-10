package dc

import (
	"context"
	"github.com/go-things/things/shared/errors"
	"github.com/go-things/things/shared/utils"
	"github.com/go-things/things/src/webapi/internal/dto"

	"github.com/go-things/things/src/webapi/internal/svc"
	"github.com/go-things/things/src/webapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BgGetGroupInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBgGetGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) BgGetGroupInfoLogic {
	return BgGetGroupInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BgGetGroupInfoLogic) BgGetGroupInfo(req types.GetGroupInfoReq) (*types.GetGroupInfoResp, error) {
	l.Infof("GetGroupInfo|req=%+v", req)
	dcReq, err := dto.GetGroupInfoReqToRpc(&req)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.DcRpc.GetGroupInfo(l.ctx, dcReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetGroupInfo|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return dto.GetGroupInfoRespToApi(resp)
}
