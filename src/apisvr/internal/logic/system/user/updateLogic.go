package user

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.UserUpdateReq) error {
	_, err := l.svcCtx.UserRpc.Update(l.ctx, &sys.UserUpdateReq{
		Uid:        req.Uid,
		UserName:   req.UserName,
		Email:      req.Email,
		Phone:      req.Phone,
		Wechat:     req.Wechat,
		NickName:   req.NickName,
		City:       req.City,
		Country:    req.Country,
		Province:   req.Province,
		Language:   req.Language,
		HeadImgUrl: req.HeadImgUrl,
		Sex:        req.Sex,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.user.upadte failure err=%+v", utils.FuncName(), er)
		return er
	}
	return nil

	return nil
}
