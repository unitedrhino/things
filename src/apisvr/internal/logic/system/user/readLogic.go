package user

import (
	"context"
	"github.com/i-Things/things/src/apisvr/internal/domain/userHeader"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.UserIndexReq) (resp *types.UserIndexResp, err error) {
	//l.Infof("UserCoreList|req=%+v", req)
	info, err := l.svcCtx.UserRpc.Read(l.ctx,
		&sys.UserReadReq{Uid: userHeader.GetUserCtx(l.ctx).Uid})
	if err != nil {
		return nil, err
	}

	return &types.UserIndexResp{
		Uid:        info.Uid,
		UserName:   info.UserName,
		InviterUid: info.InviterUid,
		InviterId:  info.InviterId,
		Sex:        info.Sex,
		City:       info.City,
		Country:    info.Country,
		Province:   info.Province,
		Language:   info.Language,
		HeadImgUrl: info.HeadImgUrl,
		CreateTime: info.CreateTime,
		Password:   info.Password,
		Email:      info.Email,
		Phone:      info.Phone,
		Wechat:     info.Wechat,
		LastIP:     info.LastIP,
		RegIP:      info.RegIP,
		Status:     info.Status,
	}, nil
}
