package logic

import (
	"context"
	"github.com/spf13/cast"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/usersvr/internal/svc"
	"gitee.com/godLei6/things/src/usersvr/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
	defer func() {
		if p := recover(); p != nil {
			utils.HandleThrow(p)
		}
	}()
	l.Infof("GetUserInfo|req=%+v",in)
	if len(in.Uid) == 0 {
		return nil, errors.Parameter.AddDetail("uid num = 0")
	}
	uis := make ([]*user.UserInfo,0,len(in.Uid))
	failUids := make([]int64,0,cap(uis)/2 + 1)
	for _,uid := range in.Uid {
		ui,err := l.svcCtx.UserInfoModel.FindOne(uid)
		if err != nil {
			failUids = append(failUids,uid)
		}else{
			uis = append(uis,&(user.UserInfo{
				Uid        : ui.Uid,
				UserName   : ui.UserName,
				NickName   : ui.NickName,
				InviterUid : ui.InviterUid,
				InviterId  : ui.InviterId,
				Sex        : ui.Sex,
				City       : ui.City,
				Country    : ui.Country,
				Province   : ui.Province,
				Language   : ui.Language,
				HeadImgUrl : ui.HeadImgUrl,
				CreateTime : ui.CreatedTime.Time.Unix(),
			}))
		}
	}
	l.Infof("GetUserInfo|allNum=%d|getNum=%d|failNum=%d|failUin=%+v|userInfo=%+v",
		len(in.Uid),len(uis),len(failUids),failUids,uis)
	var err  *errors.CodeError = nil
	if len(failUids) > 0 {
		err = errors.GetInfoPartFailure
		for _, Uid := range failUids  {
			err.AddDetail(cast.ToString(Uid))
		}
		return &user.GetUserInfoResp{Info: uis}, err
	}
	return &user.GetUserInfoResp{Info: uis}, nil
}
