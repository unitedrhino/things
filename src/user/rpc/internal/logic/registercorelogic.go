package logic

import (
	"context"
	"database/sql"
	"time"
	"yl/shared/define"
	"yl/shared/errors"
	"yl/shared/utils"
	"yl/src/user/common"
	"yl/src/user/model"

	"yl/src/user/rpc/internal/svc"
	"yl/src/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type RegisterCoreLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterCoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterCoreLogic {
	return &RegisterCoreLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}


func (l *RegisterCoreLogic)getRet(uc *model.UserCore)(*user.RegisterCoreResp, error){
	return &user.RegisterCoreResp{
		Uid: uc.Uid,
	},nil
}

func (l *RegisterCoreLogic) handlePhone(req *user.RegisterCoreReq) (*user.RegisterCoreResp, error){
	if !utils.IsMobile(req.Note){
		return nil, errors.Parameter.ToRpc()
	}
	if req.CodeID != "6666"{
		return nil, errors.Captcha.ToRpc()
	}
	//ip,err:=utils.GetIP(l.r)
	//fmt.Printf("ip=%s|err=%#v\n",ip)
	uc,err := l.svcCtx.UserCoreModel.FindOneByPhone(req.Note)
	switch err{
	case nil://如果已经有该账号,如果是注册了第一步,第二步没有注册,那么直接放行
		if uc.Status == define.NotRegistStatus{
			return l.getRet(uc)
		}
		return nil, errors.DuplicateMobile.AddDetail(req.Note).ToRpc()
	case model.ErrNotFound: //如果没有注册过,那么注册账号并进入下一步
		uc := model.UserCore{
			Uid: common.UserID.GetSnowflakeId(),
			Phone: req.Note,
			CreatedTime: sql.NullTime{Valid: true,Time: time.Now()},
		}
		_,err := l.svcCtx.UserCoreModel.Insert(uc)
		if err != nil {
			break
		}
		return l.getRet(&uc)
	default:
		break
	}
	return nil, errors.System.ToRpc()
}



func (l *RegisterCoreLogic) RegisterCore(in *user.RegisterCoreReq) (*user.RegisterCoreResp, error) {
	switch in.ReqType {
	case "wechat":
	case "phone":
		return l.handlePhone(in)
	default:
	}
	l.Errorf("%s|ReqType=%s| not suppot yet",utils.FuncName(),in.ReqType)
	return nil, errors.Parameter.AddDetail("reqType:"+in.ReqType).ToRpc()
}
