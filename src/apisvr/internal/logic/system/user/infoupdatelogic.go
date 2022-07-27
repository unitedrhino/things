package user

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/users"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/usersvr/user"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
)

type InfoUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoUpdateLogic {
	return &InfoUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoUpdateLogic) InfoUpdate(req *types.UserInfoUpdateReq) error {
	//l.Infof("ModifyUserInfo|uid=%d|req=%+v", uid, req)
	_, err := l.svcCtx.UserRpc.InfoUpdate(l.ctx, &user.UserInfoUpdateReq{
		Uid:        cast.ToString(users.GetClaimsFromToken(l.ctx, types.USER_UID, "").Uid),
		NickName:   req.NickName,
		InviterUid: req.InviterUid,
		InviterId:  req.InviterId,
		Sex:        req.Sex,
		City:       req.City,
		Country:    req.Country,
		Province:   req.Province,
		Language:   req.Language,
		HeadImgUrl: req.HeadImgUrl,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("ModifyUserInfo failure|err=%+v", er)
		return er
	}
	return nil

	return nil
}
