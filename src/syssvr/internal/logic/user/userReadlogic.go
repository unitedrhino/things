package userlogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReadLogic) UserRead(in *sys.UserReadReq) (*sys.UserInfo, error) {
	ui, err := l.svcCtx.UserInfoModel.FindOne(l.ctx, cast.ToInt64(in.UserID))
	if err != nil {
		l.Logger.Error("UserInfoModel.FindOne err , sql:%s", l.svcCtx)
		return nil, err
	}

	return &sys.UserInfo{
		UserID:      ui.UserID,
		UserName:    ui.UserName.String,
		Email:       ui.Email.String,
		Phone:       ui.Phone.String,
		Wechat:      ui.Wechat.String,
		LastIP:      ui.LastIP,
		RegIP:       ui.RegIP,
		NickName:    ui.NickName,
		City:        ui.City,
		Country:     ui.Country,
		Province:    ui.Province,
		Language:    ui.Language,
		HeadImgUrl:  ui.HeadImgUrl,
		CreatedTime: ui.CreatedTime.Unix(),
		Role:        ui.Role,
		Sex:         ui.Sex,
		IsAllData:   ui.IsAllData,
	}, nil
}
