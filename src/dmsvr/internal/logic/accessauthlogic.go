package logic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/dm"
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"
	"github.com/i-Things/things/src/dmsvr/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AccessAuthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAccessAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccessAuthLogic {
	return &AccessAuthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

var (
	AccessMap = map[string]devices.DIRECTION{
		def.PUB: devices.UP,
		def.SUB: devices.DOWN,
	}
)

func (l *AccessAuthLogic) CompareTopic(in *dm.AccessAuthReq) (err error) {
	l.Infof("%s|in:%v", utils.FuncName(), utils.GetJson(in))
	lg, err := device.GetClientIDInfo(in.ClientID)
	if err != nil {
		return err
	}
	topicInfo, err := devices.GetTopicInfo(in.Topic)
	if err != nil {
		return errors.Permissions
	}
	/*
		系统topic及物模型topic都是
			第一个表示大的功能(如$thing,$ota)
			第二个表示上行还是下行
			中间为自定义字段
			以产品id/设备名结尾
	*/
	if access, ok := AccessMap[in.Access]; !ok {
		return errors.Permissions
	} else if access != topicInfo.Direction {
		return errors.Permissions
	}
	if topicInfo.ProductID != lg.ProductID || topicInfo.DeviceName != lg.DeviceName {
		return errors.Permissions
	}
	return nil
}

func (l *AccessAuthLogic) AccessAuth(in *dm.AccessAuthReq) (*dm.Response, error) {
	l.Infof("%s|req=%+v", utils.FuncName(), utils.GetJson(in))
	err := l.CompareTopic(in)
	if err != nil {
		l.Infof("%s|auth failure=%v", utils.FuncName(), err)
	}
	return &dm.Response{}, err
}
