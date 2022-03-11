package deviceTemplate

type TempParam struct {
	ID       string `json:"id"`       //标识符
	Name     string `json:"name"`     //功能名称
	Desc     string `json:"gesc"`     //描述
	Mode     string `json:"mode"`     //读写乐行:rw(可读可写) r(只读)
	Required bool   `json:"required"` //是否必须
	Type     string `json:"type"`     //事件类型: 信息:info  告警alert  故障:fault
	Value    struct {
		Type   string            `json:"type"`              //参数类型:bool int string struct float timestamp array enum
		Maping map[string]string `json:"mapping,omitempty"` //枚举及bool类型:bool enum
		Min    string            `json:"min,omitempty"`     //数值最小值:int string float
		Max    string            `json:"max,omitempty"`     //数值最大值:int string float
		Start  string            `json:"start,omitempty"`   //初始值:int float
		Step   string            `json:"step,omitempty"`    //步长:int float
		Unit   string            `json:"unit,omitempty"`    //单位:int float
		Value  interface{}       `json:"Value"`
		/*
			读到的数据  如果是是数组则类型为[]interface{}  如果是结构体类型则为map[id]TempParam
				interface 为数据内容  					string为结构体的key value 为数据内容
		*/
	} `json:"Value"` //数据定义
}

func (t *TempParam) AddDefine(d *Define, val interface{}) (err error) {
	t.Value.Type = d.Type
	t.Value.Type = d.Type
	t.Value.Maping = make(map[string]string)
	for k, v := range d.Maping {
		t.Value.Maping[k] = v
	}
	t.Value.Maping = d.Maping
	t.Value.Min = d.Min
	t.Value.Max = d.Max
	t.Value.Start = d.Start
	t.Value.Step = d.Step
	t.Value.Unit = d.Unit
	t.Value.Value, err = d.GetVal(val)
	return err
}

func ToVal(tp map[string]TempParam) map[string]interface{} {
	ret := make(map[string]interface{}, len(tp))
	for k, v := range tp {
		ret[k] = v.ToVal()
	}
	return ret
}

func (t *TempParam) ToVal() interface{} {
	if t == nil {
		panic("TempParam is nil")
	}

	switch t.Value.Type {
	case STRUCT:
		v, ok := t.Value.Value.(map[string]TempParam)
		if ok == false {
			return nil
		}
		val := make(map[string]interface{}, len(v)+1)
		for _, tp := range v {
			val[tp.ID] = tp.ToVal()
		}
		return val
	case ARRAY:
		array, ok := t.Value.Value.([]interface{})
		if ok == false {
			return nil
		}
		val := make([]interface{}, 0, len(array)+1)
		for _, value := range array {
			switch value.(type) {
			case map[string]TempParam:
				valMap := make(map[string]interface{}, len(array)+1)
				for _, tp := range value.(map[string]TempParam) {
					valMap[tp.ID] = tp.ToVal()
				}
				val = append(val, valMap)
			default:
				val = append(val, value)
			}
		}
		return val
	default:
		return t.Value.Value
	}
}
