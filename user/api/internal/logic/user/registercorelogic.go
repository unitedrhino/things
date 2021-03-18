package logic

import (
	"context"
	"database/sql"
	"time"
	"yl/shared/define"
	"yl/shared/utils"
	"yl/user/common"
	"yl/user/model"

	"yl/user/api/internal/svc"
	"yl/user/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type RegisterCoreLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterCoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) RegisterCoreLogic {
	return RegisterCoreLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterCoreLogic)getRet(uc *model.UserCore)(*types.RegisterCoreResp, error){
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Rej.AccessExpire
	jwtToken, err := utils.GetJwtToken(l.svcCtx.Config.Rej.AccessSecret, now, accessExpire, uc.Uid)
	if err != nil {
		return nil, err
	}
	return &types.RegisterCoreResp{
		JwtToken :types.JwtToken{
			AccessToken:  jwtToken,
			AccessExpire: now + accessExpire,
			RefreshAfter: now + accessExpire/2,
		},
		Uid: uc.Uid,
	},nil
}

func (l *RegisterCoreLogic) handlePhone(req types.RegisterCoreReq) (*types.RegisterCoreResp, error){
	if !utils.IsMobile(req.Note){
		return nil, common.ErrorParameter
	}
	if req.CodeID != "6666"{
		return nil, common.ErrorCaptcha
	}
	//ip,err:=utils.GetIP(l.r)
	//fmt.Printf("ip=%s|err=%#v\n",ip)
	uc,err := l.svcCtx.UserCoreModel.FindOneByPhone(sql.NullString{String: req.Note,Valid: true})
	switch err{
	case nil://如果已经有该账号,如果是注册了第一步,第二步没有注册,那么直接放行
		if uc.Status == define.NotRegistStatus{
			return l.getRet(uc)
		}
		return nil, common.ErrorDuplicateMobile
	case model.ErrNotFound://如果没有注册过,那么注册账号并进入下一步
		uc := model.UserCore{
			Uid:common.UserID.GetSnowflakeId(),
			Phone: sql.NullString{
				String: req.Note,
				Valid: true,
			},
		}
		_,err := l.svcCtx.UserCoreModel.Insert(uc)
		if err != nil {
			break
		}
		return l.getRet(&uc)
	default:
		break
	}
	return nil, common.ErrorSystem
}
func (l *RegisterCoreLogic) RegisterCore(req types.RegisterCoreReq) (*types.RegisterCoreResp, error) {
	switch req.RegType {
	case "wechat":
		logx.Error("wechat not suppot yet")
	case "phone":
		return l.handlePhone(req)
	default:
		return nil, common.ErrorParameter
	}
	return &types.RegisterCoreResp{}, common.ErrorParameter
}
