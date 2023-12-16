package usermanagelogic

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
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegister1Logic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UiDB *relationDB.UserInfoRepo
}

func NewUserRegister1Logic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegister1Logic {
	return &UserRegister1Logic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UiDB:   relationDB.NewUserInfoRepo(ctx),
	}
}

func (l *UserRegister1Logic) UserRegister1(in *sys.UserRegister1Req) (*sys.UserRegister1Resp, error) {
	l.Infof("%v.req:%v", utils.FuncName(), utils.Fmt(in))
	if in.TenantCode == "" {
		in.TenantCode = def.TenantCodeDefault
	}
	switch in.RegType {
	case users.RegWxMiniP:
		return l.handleWxminip(in)
	default:
		return nil, errors.NotRealize.AddMsgf(in.RegType)
	}
}
func (l *UserRegister1Logic) handleWxminip(in *sys.UserRegister1Req) (*sys.UserRegister1Resp, error) {
	auth := l.svcCtx.WxMiniProgram.GetAuth()
	ret, err := auth.Code2SessionContext(l.ctx, in.Code)
	if err != nil {
		l.Errorf("%v.Code2SessionContext err:%v", err)
		if ret.ErrCode != 0 {
			return nil, errors.System.AddDetail(ret.ErrMsg)
		}
		return nil, errors.System.AddDetail(err)
	} else if ret.ErrCode != 0 {
		return nil, errors.Parameter.AddDetail(ret.ErrMsg)
	}
	_, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{WechatUnionID: ret.UnionID})
	if err == nil { //已经注册过
		return nil, errors.DuplicateRegister
	}
	if !errors.Cmp(err, errors.NotFind) {
		return nil, err
	}
	//如果没有注册过,则进入注册逻辑
	if l.svcCtx.Config.Register.NeedDetail == true { //如果需要填写用户信息
		token, err := users.GetRegisterToken(l.svcCtx.Config.Register.SecondSecret, time.Now().Unix(),
			l.svcCtx.Config.Register.SecondExpire, users.RegWxOpen, ret.UnionID, 0)
		if err != nil {
			return nil, errors.System.AddDetail(err)
		}
		return &sys.UserRegister1Resp{Token: token}, nil
	}
	err = l.UiDB.Insert(l.ctx, &relationDB.SysUserInfo{
		UserID:        l.svcCtx.UserID.GetSnowflakeId(),
		WechatUnionID: sql.NullString{Valid: true, String: ret.UnionID},
		RegIP:         in.RegIP,
		Role:          l.svcCtx.Config.Register.DefaultRole,
	})
	if err != nil {
		return nil, err
	}
	return &sys.UserRegister1Resp{}, nil
}
