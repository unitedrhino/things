package user

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.UserCreateReq) (resp *types.UserCreateResp, err error) {
	l.Infof("%s req=%+v", utils.FuncName(), req)
	resp1, err1 := l.svcCtx.UserRpc.Create(l.ctx, &sys.UserCreateReq{
		ReqType:    req.ReqType,
		UserName:   req.UserName,
		Password:   req.Password,
		Wechat:     req.Wechat,
		LastIP:     req.LastIP,
		RegIP:      req.RegIP,
		NickName:   req.NickName,
		City:       req.City,
		Country:    req.Country,
		Province:   req.Province,
		Language:   req.Language,
		HeadImgUrl: req.HeadImgUrl,
		Role:       req.Role,
		Sex:        req.Sex,
	})
	if err1 != nil {
		er := errors.Fmt(err1)
		l.Errorf("%s.rpc.Register req=%v err=%v rpc_err=%v", utils.FuncName(), req, er, err)
		return &types.UserCreateResp{}, er
	}
	if resp1 == nil {
		l.Errorf("%s.rpc.Register return nil req=%v", utils.FuncName(), req)
		return &types.UserCreateResp{}, errors.System.AddDetail("register core rpc return nil")
	}

	return &types.UserCreateResp{Uid: resp1.Uid}, nil

	return
}
