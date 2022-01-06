package logic

import (
	"context"
	"github.com/go-things/things/shared/errors"
	"github.com/go-things/things/src/dmsvr/dm"

	"github.com/go-things/things/src/dcsvr/dc"
	"github.com/go-things/things/src/dcsvr/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

type SendPropertyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendPropertyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendPropertyLogic {
	return &SendPropertyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 同步调用设备属性
func (l *SendPropertyLogic) SendProperty(in *dc.SendPropertyReq) (*dc.SendPropertyResp, error) {
	l.Infof("SendProperty|req=%+v", in)
	deviceMemberID := in.ProductID + ":" + in.DeviceName
	ok, err := l.svcCtx.DcDB.CheckMemeberWithGoupID(in.MemberID, in.MemberType, deviceMemberID, 1)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	if ok != true { //如果设备和发送者不在同一个组
		return nil, errors.Permissions
	}
	resp, err := l.svcCtx.Dmsvr.SendProperty(l.ctx, &dm.SendPropertyReq{
		ProductID:     in.ProductID,
		DeviceName:    in.DeviceName,
		Data:          in.Data,
		DataTimestamp: in.DataTimestamp,
		Method:        in.Method,
	})
	if err != nil {
		return nil, err
	}
	return &dc.SendPropertyResp{
		Data:        resp.Data,
		ClientToken: resp.ClientToken, //调用id
		Status:      resp.Status,      //返回状态
		Code:        resp.Code,        //设备返回状态码
	}, nil
}
