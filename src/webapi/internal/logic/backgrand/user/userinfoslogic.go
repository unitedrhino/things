package user

import (
	"context"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/usersvr/user"
	"github.com/spf13/cast"

	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type UserInfosLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfosLogic(ctx context.Context, svcCtx *svc.ServiceContext) UserInfosLogic {
	return UserInfosLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfosLogic) UserInfos(req types.GetUserInfosReq) (*types.GetUserInfosResp, error) {
	l.Infof("UserInfos|req=%d", req)
	uids := make([]int64, 0, len(req.Uid))
	for _, uid := range req.Uid {
		uids = append(uids, cast.ToInt64(uid))
	}
	uis, err := l.svcCtx.UserRpc.GetUserInfo(l.ctx, &user.GetUserInfoReq{Uid: uids})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("[%s]|rpc.Login|uid=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	resp := types.GetUserInfosResp{}
	resp.Info = make([]*types.UserInfo, 0, len(uis.Info))
	for _, ui := range uis.Info {
		resp.Info = append(resp.Info, types.UserInfoToApi(ui))
	}
	return &resp, nil
}
