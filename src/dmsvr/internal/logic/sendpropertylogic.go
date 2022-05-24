package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/domain/service/deviceSend"
	"github.com/i-Things/things/src/dmsvr/internal/domain/thing"
	"time"

	"github.com/i-Things/things/src/dmsvr/dm"
	"github.com/i-Things/things/src/dmsvr/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendPropertyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	template *thing.Template
}

func NewSendPropertyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendPropertyLogic {
	return &SendPropertyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendPropertyLogic) initMsg(productID string) error {
	var err error
	l.template, err = l.svcCtx.TemplateRepo.GetTemplate(l.ctx, productID)
	if err != nil {
		return errors.System.AddDetail(err.Error())
	}
	return nil
}

func (l *SendPropertyLogic) SendProperty(in *dm.SendPropertyReq) (*dm.SendPropertyResp, error) {
	l.Infof("SendProperty|req=%+v", in)
	err := l.initMsg(in.ProductID)
	if err != nil {
		return nil, err
	}
	param := map[string]interface{}{}
	err = json.Unmarshal([]byte(in.Data), &param)
	if err != nil {
		return nil, errors.Parameter.AddDetail(
			"SendProperty|data not right:", in.Data)
	}
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		l.Errorf("SendProperty|GenerateUUID err:%v", err)
		return nil, errors.System.AddDetail(err)
	}
	req := deviceSend.DeviceReq{
		Method:      deviceSend.CONTROL,
		ClientToken: uuid,
		//ClientToken:"de65377c-4041-565d-0b5e-67b664a06be8",//这个是测试代码
		Timestamp: time.Now().UnixMilli(),
		Params:    param}
	_, err = req.VerifyReqParam(l.template, thing.ACTION_INPUT)
	if err != nil {
		return nil, err
	}
	pubTopic := fmt.Sprintf("$thing/down/property/%s/%s", in.ProductID, in.DeviceName)
	subTopic := fmt.Sprintf("$thing/up/property/%s/%s", in.ProductID, in.DeviceName)

	resp, err := l.svcCtx.InnerLink.ReqToDeviceSync(l.ctx, pubTopic, subTopic, &req, in.ProductID, in.DeviceName)
	if err != nil {
		return nil, err
	}

	return &dm.SendPropertyResp{
		ClientToken: resp.ClientToken,
		Status:      resp.Status,
		Code:        resp.Code,
	}, nil

	return &dm.SendPropertyResp{}, nil
}
