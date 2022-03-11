package deviceSend

type LogicHandle interface {
	Handle(msg *Elements) error
}
