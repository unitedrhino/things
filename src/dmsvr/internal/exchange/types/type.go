package types

import "time"

//Elements kafka publish elements
type Elements struct {
	ClientID  string `json:"clientid"`
	Username  string `json:"username"`
	Topic     string `json:"topic"`
	Payload   string `json:"payload"`
	Timestamp int64  `json:"ts"`
	Size      int32  `json:"size"`
	Action    string `json:"action"`
}

type Info struct {
	Timeout time.Time
	ClientID  string `json:"clientid"`
	Msg  chan *Elements
}

func NewInfo(timeout time.Time,ClientID string) *Info{
	return &Info{
		Timeout: timeout,
		ClientID:ClientID,
		Msg:  make(chan *Elements,1),
	}
}

func (i *Info)IsTimeOut()bool{
	return time.Now().Before(i.Timeout)
}