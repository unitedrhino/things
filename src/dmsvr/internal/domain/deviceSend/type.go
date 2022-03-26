package deviceSend

//Elements kafka publish elements
type Elements struct {
	ProductID  string
	DeviceName string
	ClientID   string
	Username   string
	Address    string

	Topic     string
	Payload   []byte
	Timestamp int64
	Action    string
	Reason    string
}
