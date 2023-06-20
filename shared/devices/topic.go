package devices

import (
	"github.com/i-Things/things/shared/errors"
	"strings"
)

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
$shadow/operation/up/{productID}/${deviceName}	发布	设备影子消息上行
$shadow/operation/down/{productID}/${deviceName}	订阅	设备影子消息下行
$rrpc/txd/{productID}/${deviceName}/${MessageId}	发布	RRPC消息上行，MessageId为RRPC消息ID
$rrpc/rxd/{productID}/${deviceName}/+	订阅	RRPC消息下行
$sys/operation/up/{productID}/${deviceName}	发布	系统topic：ntp服务消息上行
$sys/operation/down/{productID}/${deviceName}/+	订阅	系统topic：ntp服务消息下行
log topic
$log/up/operation/${productID}/${deviceName} //设备查询是否需要上传调试日志及日志级别，上行
$log/down/operation/${productID}/${deviceName}
$log/up/report/${productID}/${deviceName} //设备上传调试日志内容，上行
$log/down/report/${productID}/${deviceName}
$log/down/update/${productID}/${deviceName} //服务器端下发调试日志配置，下行

自定义topic:
${productID}/${deviceName}/control	订阅	编辑删除
${productID}/${deviceName}/data	订阅和发布	编辑删除
${productID}/${deviceName}/event	发布
${productID}/${deviceName}/xxxxx	订阅和发布   //自定义 暂不做支持
*/

const (
	TopicHeadThing   = "$thing"
	Thing            = "thing"
	TopicHeadOta     = "$ota"
	Ota              = "ota"
	TopicHeadConfig  = "$config"
	Config           = "config"
	TopicHeadLog     = "$log"
	Log              = "log"
	TopicHeadShadow  = "$shadow"
	Shadow           = "shadow"
	TopicHeadGateway = "$gateway"
	Gateway          = "gateway"
	TopicHeadExt     = "$ext"
	Ext              = "ext"
)

//设备通信流向
type Direction int

const (
	Unknown Direction = iota //设备通信流向：未知
	Up                       //设备通信流向：上行
	Down                     //设备通信流向：下行
)

type TopicInfo struct {
	ProductID  string
	DeviceName string
	Direction  Direction
	Types      []string
	TopicHead  string
}

func GetTopicInfo(topic string) (topicInfo *TopicInfo, err error) {
	keys := strings.Split(topic, "/")
	return parseTopic(keys)
}

// 通过topic的第一个字段来获取处理函数
func parseTopic(topics []string) (topicInfo *TopicInfo, err error) {
	if len(topics) < 2 {
		return nil, errors.Parameter.AddDetail("topic is err")
	}
	switch topics[0] {
	case TopicHeadThing, TopicHeadOta, TopicHeadShadow, TopicHeadLog, TopicHeadConfig, TopicHeadGateway, TopicHeadExt:
		return parseLast(topics)
	default: //自定义消息
		return parsePose(0, topics)
	}
}

func parsePose(productPos int, topics []string) (topicInfo *TopicInfo, err error) {
	return nil, errors.Parameter.AddDetail("topic is err")
	//先不考虑自定义消息
	//if len(topics) < (productPos + 2) {
	//	return nil, errors.Parameter.AddDetail("topic is err")
	//}
	//return &TopicInfo{
	//	ProductID:  topics[productPos],
	//	DeviceName: topics[productPos+1],
	//	TopicHead:  topics[0],
	//}, err
}

func parseLast(topics []string) (topicInfo *TopicInfo, err error) {
	if len(topics) < 4 {
		return nil, errors.Parameter.AddDetail("topic is err")
	}
	return &TopicInfo{
		ProductID:  topics[len(topics)-2],
		DeviceName: topics[len(topics)-1],
		Direction:  getDirection(topics[1]),
		Types:      topics[2 : len(topics)-2],
		TopicHead:  topics[0],
	}, err
}

func getDirection(dir string) Direction {
	switch dir {
	case "up":
		return Up
	case "down":
		return Down
	default:
		return Unknown
	}
}
