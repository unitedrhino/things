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

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.UserLoginReq) (resp *types.UserLoginResp, err error) {
	l.Infof("%s req=%+v", utils.FuncName(), req)
	if req.LoginType == "pwd" {
		if l.svcCtx.Captcha.Verify(req.CodeID, req.Code) == false {
			return nil, errors.Captcha
		}
	}
	uResp, err := l.svcCtx.UserRpc.Login(l.ctx, &sys.LoginReq{
		UserID:    req.UserID,
		PwdType:   req.PwdType,
		Password:  req.Password,
		LoginType: req.LoginType,
		Code:      req.Code,
		CodeID:    req.CodeID,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.Login req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	if uResp == nil {
		l.Errorf("%s.rpc.Register return nil req=%v", utils.FuncName(), req)
		return nil, errors.System.AddDetail("register core rpc return nil")
	}
	return &types.UserLoginResp{
		Info: types.UserInfo{
			Uid:         uResp.Info.Uid,
			UserName:    uResp.Info.UserName,
			Password:    "",
			Email:       uResp.Info.Email,
			Phone:       uResp.Info.Phone,
			Wechat:      uResp.Info.Wechat,
			LastIP:      uResp.Info.LastIP,
			RegIP:       uResp.Info.RegIP,
			NickName:    uResp.Info.NickName,
			City:        uResp.Info.City,
			Country:     uResp.Info.Country,
			Province:    uResp.Info.Province,
			Language:    uResp.Info.Language,
			HeadImgUrl:  uResp.Info.HeadImgUrl,
			CreatedTime: uResp.Info.CreatedTime,
			Role:        uResp.Info.Role,
			Sex:         uResp.Info.Sex,
		},
		Token: types.JwtToken{
			AccessToken:  uResp.Token.AccessToken,
			AccessExpire: uResp.Token.AccessExpire,
			RefreshAfter: uResp.Token.RefreshAfter,
		},
	}, nil

	return
}
