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

type BgManageGroupInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBgManageGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) BgManageGroupInfoLogic {
	return BgManageGroupInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BgManageGroupInfoLogic) BgManageGroupInfo(req types.ManageGroupInfoReq) (*types.GroupInfo, error) {
	l.Infof("ManageGroupInfo|req=%+v", req)
	dcReq,err := dto.ManageGroupInfoReqToRpc(&req)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.DcRpc.ManageGroupInfo(l.ctx, dcReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.ManageGroupInfo|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return dto.GrouInfoToApi(resp), nil

	return &types.GroupInfo{}, nil
}
