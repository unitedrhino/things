package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/godLei6/things/src/dmsvr/device/model"
	"gitee.com/godLei6/things/src/dmsvr/dm"
	"gitee.com/godLei6/things/src/dmsvr/internal/exchange/types"
	"gitee.com/godLei6/things/src/dmsvr/internal/svc"
	"github.com/tal-tech/go-zero/core/logx"
)

type PublishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) LogicHandle {
	return LogicHandle(&PublishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	})
}

func (l *PublishLogic) Handle(msg *types.Elements) error {
	l.Infof("PublishLogic|req=%+v", msg)
	ld, err :=dm.GetClientIDInfo(msg.ClientID)
	if err != nil {
		return err
	}
	pi,err := l.svcCtx.ProductInfo.FindOneByProductID(ld.ProductID)
	if err != nil {
		return err
	}
	var template model.Template
	err = json.Unmarshal([]byte(pi.Template),&template)
	if err != nil {
		return err
	}
	var deviceData model.DeviceReq
	err = json.Unmarshal([]byte(msg.Payload),&deviceData)
	if err != nil {
		return err
	}
	fmt.Printf("template=%+v|req=%+v\n",template,deviceData)
	return nil
}
