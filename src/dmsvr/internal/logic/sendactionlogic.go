package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/src/dmsvr/device"
	"gitee.com/godLei6/things/src/dmsvr/internal/repo/model/mysql"
	"time"

	"gitee.com/godLei6/things/src/dmsvr/dm"
	"gitee.com/godLei6/things/src/dmsvr/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

type SendActionLogic struct {
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	pi       *mysql.ProductInfo
	template *device.Template
	logx.Logger
}

func NewSendActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendActionLogic {
	return &SendActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *SendActionLogic) initMsg(productID string) error {
	var err error
	l.pi, err = l.svcCtx.ProductInfo.FindOneByProductID(productID)
	if err != nil {
		return err
	}
	l.template, err = device.NewTemplate([]byte(l.pi.Template))
	if err != nil {
		return err
	}
	return nil
}

func (l *SendActionLogic) SendAction(in *dm.SendActionReq) (*dm.SendActionResp, error) {
	l.Infof("SendAction|req=%+v", in)
	err := l.initMsg(in.ProductID)
	if err != nil {
		return nil, err
	}
	param := map[string]interface{}{}
	err = json.Unmarshal([]byte(in.InputParams), &param)
	if err != nil {
		return nil, errors.Parameter.AddDetail("SendAction|InputParams not right:", in.InputParams)
	}
	//uuid,err := uuid.GenerateUUID()
	//if err != nil{
	//	l.Errorf("SendAction|GenerateUUID err:%v",err)
	//	return nil, errors.System.AddDetail(err)
	//}
	req := device.DeviceReq{
		Method: device.ACTION,
		//ClientToken: uuid,
		ClientToken: "de65377c-4041-565d-0b5e-67b664a06be8", //这个是测试代码
		Timestamp:   time.Now().UnixMilli(),
		Params:      param}
	l.template.VerifyReqParam(req, device.ACTION_INPUT)
	pubTopic := fmt.Sprintf("$thing/down/action/%s/%s", in.ProductID, in.DeviceName)
	subTopic := fmt.Sprintf("$thing/up/action/%s/%s", in.ProductID, in.DeviceName)

	resp, err := l.svcCtx.DevClient.DeviceReq(l.ctx, req, pubTopic, subTopic)
	if err != nil {
		return nil, err
	}
	respParam, err := json.Marshal(resp.Response)
	if err != nil {
		return nil, errors.RespParam.AddDetailf("SendAction|get device resp not right:%+v", resp.Response)
	}
	return &dm.SendActionResp{
		ClientToken:  resp.ClientToken,
		Status:       resp.Status,
		Code:         resp.Code,
		OutputParams: string(respParam),
	}, nil
}
