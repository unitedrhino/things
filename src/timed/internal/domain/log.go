package domain

type ScriptLog struct {
	Level       string `json:"level"`   //info warn error
	Content     string `json:"content"` //上下文
	CreatedTime int64  `json:"createdTime"`
}
