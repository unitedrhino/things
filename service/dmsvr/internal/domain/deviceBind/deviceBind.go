package deviceBind

type TokenInfo struct {
	Token  string `json:"token"`
	Status Status `json:"status"`
	UserID int64  `json:"userID,string"`
}

type Status = int64

const (
	StatusInit   Status = 1
	StatusReport Status = 2
)
