// 消息队列通知类型
package job

type Queue struct {
	Topic   string `json:"topic"`
	Payload string `json:"payload"`
}
