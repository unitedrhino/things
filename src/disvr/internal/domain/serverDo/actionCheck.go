package serverDo

import "encoding/json"

type SendOption struct {
	TimeoutToFail  int64 `json:"timeoutToFail,optional"`  //超时失败时间
	RequestTimeout int64 `json:"requestTimeout,optional"` //请求超时,超时后会进行重试
	RetryInterval  int64 `json:"retryInterval,optional"`  //重试间隔
}

func (s SendOption) String() string {
	ret, _ := json.Marshal(s)
	return string(ret)
}
