package def

type (
	// MsgHead 消息队列的头
	//todo 后续考虑用proto重构这个头
	MsgHead struct {
		Trace     string //追踪tid
		Timestamp int64  //发送时毫秒级时间戳
		Data      []byte //传送的内容
	}
)
