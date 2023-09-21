package userlogic

import (
	"context"
	"database/sql"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/users"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegister2Logic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UiDB *relationDB.UserInfoRepo
}

func NewUserRegister2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegister2Logic {
	return &UserRegister2Logic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UiDB:   relationDB.NewUserInfoRepo(ctx),
	}
}

func (l *UserRegister2Logic) UserRegister2(in *sys.UserRegister2Req) (*sys.Response, error) {
	var tokenInfo users.RegisterClaims
	err := users.ParseToken(&tokenInfo, in.Token, l.svcCtx.Config.Register.SecondSecret)
	if err != nil {
		return nil, err
	}
	userID := l.svcCtx.UserID.GetSnowflakeId()
	var userInfo = relationDB.SysUserInfo{
		UserID:     userID,
		UserName:   sql.NullString{Valid: true, String: in.Info.UserName},
		Password:   utils.MakePwd(in.Info.Password, userID, false),
		RegIP:      in.RegIP,
		NickName:   in.Info.NickName,
		Sex:        in.Info.Sex,
		City:       in.Info.City,
		Country:    in.Info.Country,
		Province:   in.Info.Province,
		Language:   in.Info.Language,
		HeadImgUrl: in.Info.HeadImgUrl,
		Role:       l.svcCtx.Config.Register.DefaultRole,
		IsAllData:  def.False,
	}
	switch tokenInfo.RejType {
	case users.RegWxOpen:
		_, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Wechat: tokenInfo.Note})
		if err == nil { //已经注册过
			return nil, errors.DuplicateRegister
		}
		if !errors.Cmp(err, errors.NotFind) {
			return nil, err
		}
		userInfo.Wechat = sql.NullString{String: tokenInfo.Note, Valid: true}
	}
	err = l.UiDB.Insert(l.ctx, &userInfo)
	if err != nil {
		if errors.Cmp(err, errors.Duplicate) {
			return nil, errors.Parameter.AddMsgf("用户名已注册:%v", userInfo.UserName)
		}
		return nil, err
	}
	return &sys.Response{}, nil
}
