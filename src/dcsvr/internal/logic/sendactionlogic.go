package logic

import (
	"context"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/src/dmsvr/dm"

	"gitee.com/godLei6/things/src/dcsvr/dc"
	"gitee.com/godLei6/things/src/dcsvr/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

type SendActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendActionLogic {
	return &SendActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 同步调用设备行为
func (l *SendActionLogic) SendAction(in *dc.SendActionReq) (*dc.SendActionResp, error) {
	l.Infof("SendAction|req=%+v", in)
	deviceMemberID := in.ProductID + ":" + in.DeviceName
	ok,err := l.svcCtx.DcDB.CheckMemeberWithGoupID(in.MemberID,in.MemberType,deviceMemberID,1)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	if ok != true {//如果设备和发送者不在同一个组
		return nil,errors.Permissions
	}
	resp,err := l.svcCtx.Dmsvr.SendAction(l.ctx,&dm.SendActionReq{
		ProductID:in.ProductID,
		DeviceName:in.DeviceName,
		ActionId:in.ActionId,
		InputParams:in.InputParams,
	})
	if err != nil {
		return nil, err
	}
	return &dc.SendActionResp{
		ClientToken:resp.ClientToken,
		OutputParams:resp.OutputParams,
		Status:resp.Status,
		Code:resp.Code,
	}, nil
}
