package deviceSend

//Elements kafka publish elements
type Elements struct {
	ClientID string
	Username string
	Address  string

	Topic     string
	Payload   string
	Timestamp int64
	Action    string
	Reason    string
}
