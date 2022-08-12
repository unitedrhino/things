package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/domain/schema"
	"github.com/i-Things/things/src/dmsvr/internal/domain/service/deviceSend"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"time"

	"github.com/i-Things/things/src/dmsvr/dmclient"
	"github.com/i-Things/things/src/dmsvr/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	schema *schema.Model
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
	l.schema, err = l.svcCtx.SchemaRepo.GetSchemaModel(l.ctx, productID)
	if err != nil {
		return errors.System.AddDetail(err)
	}
	return nil
}

func (l *SendActionLogic) SendAction(in *dmclient.SendActionReq) (*dmclient.SendActionResp, error) {
	l.Infof("SendAction|req=%+v", in)
	err := l.initMsg(in.ProductID)
	if err != nil {
		return nil, err
	}
	param := map[string]any{}
	err = json.Unmarshal([]byte(in.InputParams), &param)
	if err != nil {
		return nil, errors.Parameter.AddDetail("SendAction|InputParams not right:", in.InputParams)
	}
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		l.Errorf("SendAction|GenerateUUID err:%v", err)
		return nil, errors.System.AddDetail(err)
	}
	req := deviceSend.DeviceReq{
		Method:      deviceSend.ACTION,
		ClientToken: uuid,
		//ClientToken: "de65377c-4041-565d-0b5e-67b664a06be8", //这个是测试代码
		Timestamp: time.Now().UnixMilli(),
		Params:    param}
	_, err = req.VerifyReqParam(l.schema, schema.ACTION_INPUT)
	if err != nil {
		return nil, err
	}
	pubTopic := fmt.Sprintf("$thing/down/action/%s/%s", in.ProductID, in.DeviceName)
	subTopic := fmt.Sprintf("$thing/up/action/%s/%s", in.ProductID, in.DeviceName)

	resp, err := l.svcCtx.PubDev.ReqToDeviceSync(l.ctx, pubTopic, subTopic, &req, in.ProductID, in.DeviceName)
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
