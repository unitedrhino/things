package protocol

type DefaultConf struct {
}

func (p DefaultConf) GenKey() string {
	return ""
}

func (p DefaultConf) Equal(in DefaultConf) bool {
	return p == in
}
