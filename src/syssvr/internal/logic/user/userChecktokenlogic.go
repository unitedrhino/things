package userlogic

import (
	"context"
	"github.com/i-Things/things/shared/users"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCheckTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckTokenLogic {
	return &CheckTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckTokenLogic) UserCheckToken(in *sys.UserCheckTokenReq) (*sys.UserCheckTokenResp, error) {
	jwt, err := users.ParseToken(in.Token, l.svcCtx.Config.UserToken.AccessSecret)
	if err != nil {
		l.Errorf("%s parse token fail err=%s", utils.FuncName(), err.Error())
		return nil, err
	}
	var token string

	if (jwt.ExpiresAt-time.Now().Unix())*2 < l.svcCtx.Config.UserToken.AccessExpire {
		token, _ = users.RefreshToken(in.Token, l.svcCtx.Config.UserToken.AccessSecret, time.Now().Unix()+l.svcCtx.Config.UserToken.AccessExpire)
	}
	return &sys.UserCheckTokenResp{
		Token:     token,
		UserID:    jwt.UserID,
		Role:      jwt.Role,
		IsAllData: jwt.IsAllData,
	}, nil
}
