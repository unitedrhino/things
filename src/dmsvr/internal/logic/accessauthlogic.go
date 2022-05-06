package logic

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"
	"strings"

	"github.com/i-Things/things/src/dmsvr/dm"
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

/*
物理型topic:
$thing/up/property/${productID}/${deviceName}	发布	属性上报
$thing/down/property/${productID}/${deviceName}	订阅	属性下发与属性上报响应
$thing/up/event/${productID}/${deviceName}	发布	事件上报
$thing/down/event/${productID}/${deviceName}	订阅	事件上报响应
$thing/up/action/${productID}/${deviceName}	发布	设备响应行为执行结果
$thing/down/action/${productID}/${deviceName}	订阅	应用调用设备行为
系统级topic:
$ota/report/${productID}/${deviceName}	发布	固件升级消息上行
$ota/update/${productID}/${deviceName}	订阅	固件升级消息下行
$broadcast/rxd/${productID}/${deviceName}	订阅	广播消息下行
自定义topic:
${productID}/${deviceName}/control	订阅	编辑删除
${productID}/${deviceName}/data	订阅和发布	编辑删除
${productID}/${deviceName}/event	发布
${productID}/${deviceName}/xxxxx	订阅和发布   //自定义 暂不做支持
*/
//key 为topic的第一个 value 为该key下的正则表达式
var TopicSub map[string][]string = map[string][]string{
	"$thing": //物理型topic:
	[]string{
		"$thing/down/property/%s/%s", //订阅   属性下发与属性上报响应
		"$thing/down/event/%s/%s",    //订阅	事件上报响应
		"$thing/down/action/%s/%s",   //订阅	应用调用设备行为
	},
	"$ota": //系统级topic:
	[]string{
		"$ota/update/%s/%s",    //	订阅	固件升级消息下行
		"$broadcast/rxd/%s/%s", //	订阅	广播消息下行
	},
	"$broadcast": //系统级广播topic:
	[]string{
		"$broadcast/rxd/%s/%s", //订阅	广播消息下行
	},
}
var TopicPub map[string][]string = map[string][]string{
	"$thing":                            //物理型topic:
	[]string{"$thing/up/property/%s/%s", //发布	属性上报
		"$thing/up/event/%s/%s",  //发布	事件上报
		"$thing/up/action/%s/%s", //发布	设备响应行为执行结果
	},
	"$ota": //系统级topic:
	[]string{
		"$ota/report/%s/%s", //	发布	固件升级消息上行
	},
}

func (l *AccessAuthLogic) CompareTopic(in *dm.AccessAuthReq) error {
	var Topic map[string][]string
	switch in.Access {
	case def.PUB:
		Topic = TopicPub
	case def.SUB:
		Topic = TopicSub
	default:
	}
	lg, err := device.GetClientIDInfo(in.ClientID)
	if err != nil {
		return err
	}
	keys := strings.Split(in.Topic, "/")
	topics, ok := Topic[keys[0]]
	if ok != true {
		//自定义topic
		if keys[0] == lg.ProductID && keys[1] == lg.DeviceName { //用户自定义先不判断后续的字段 undo
			return nil
		}
	}
	for _, v := range topics { //把所有topic的前缀比较一次
		prefix := fmt.Sprintf(v, lg.ProductID, lg.DeviceName)
		if !strings.HasPrefix(in.Topic, prefix) { //如果没有比较到该前缀
			continue
		}
		if len(in.Topic) == len(prefix) {
			return nil
		}
		if in.Topic[len(prefix)] != '/' {
			return errors.Permissions
		}
		return nil
	}
	return errors.Permissions
}

func (l *AccessAuthLogic) AccessAuth(in *dm.AccessAuthReq) (*dm.Response, error) {
	l.Infof("AccessAuth|req=%+v", in)
	return &dm.Response{}, l.CompareTopic(in)
}
