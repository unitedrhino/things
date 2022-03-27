package eventDevSub

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceSend"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceTemplate"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"time"
)

//设备回复消息处理

type DevPublishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	topics []string
	dresp  deviceTemplate.DeviceResp
}

func NewDevPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DevPublishLogic {
	return &DevPublishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DevPublishLogic) initMsg(msg *deviceSend.Elements) error {
	var err error
	if err != nil {
		return err
	}
	err = utils.Unmarshal(msg.Payload, &l.dresp)
	if err != nil {
		return errors.Parameter.AddDetail("things topic is err:" + msg.Topic)
	}
	l.topics = strings.Split(msg.Topic, "/")
	return nil
}

func (l *DevPublishLogic) HandleProperty(msg *deviceSend.Elements) error {
	l.Slowf("DevPublishLogic|HandleProperty")
	switch l.dresp.Method {
	case deviceTemplate.CONTROL_REPLY:
		return l.HandleResp(msg)
	default:
		return errors.Method
	}
	return nil
}

func (l *DevPublishLogic) HandleResp(msg *deviceSend.Elements) error {
	l.Slowf("DevPublishLogic|HandleResp")
	//todo 这里后续需要处理异步获取消息的情况
	return nil
}

func (l *DevPublishLogic) HandleThing(msg *deviceSend.Elements) error {
	l.Slowf("DevPublishLogic|HandleThing")
	if len(l.topics) < 5 || l.topics[1] != "up" {
		return errors.Parameter.AddDetail("things topic is err:" + msg.Topic)
	}
	switch l.topics[2] {
	case def.PROPERTY_METHOD: //属性上报
		return l.HandleProperty(msg)
	case def.ACTION_METHOD: //设备响应行为执行结果
		return l.HandleResp(msg)
	default:
		return errors.Method
	}
	return nil
}

func (l *DevPublishLogic) Handle(msg *deviceSend.Elements) (err error) {
	l.Infof("DevPublishLogic|req=%+v", msg)
	err = l.initMsg(msg)
	if err != nil {
		return err
	}
	if len(l.topics) > 1 {
		switch l.topics[0] {
		case "$thing":
			err = l.HandleThing(msg)
		default:
			err = errors.Parameter.AddDetailf("not suppot topic :%s", msg.Topic)
		}
	} else {
		err = errors.Parameter.AddDetailf("need topic :%s", msg.Topic)
	}
	l.svcCtx.DeviceLog.Insert(&mysql.DeviceLog{
		ProductID:   msg.ProductID,
		Action:      "publish",
		Timestamp:   l.dresp.GetTimeStamp(time.UnixMilli(msg.Timestamp)), // 操作时间
		DeviceName:  msg.DeviceName,
		TranceID:    utils.TraceIdFromContext(l.ctx),
		RequestID:   l.dresp.ClientToken,
		Content:     string(msg.Payload),
		Topic:       msg.Topic,
		ResultType:  errors.Fmt(err).GetCode(),
		CreatedTime: time.Now(),
	})
	return err
}
