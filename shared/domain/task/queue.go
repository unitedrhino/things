// 消息队列通知类型
package task

type Queue struct {
	Topic   string `json:"topic"`
	Payload string `json:"payload"`
}
