package scene

type When struct {
	TermCondType TermCondType `json:"termCondType"` //触发条件类型 and: 与 or: 或
	ValidRange   *WhenRange   `json:"validRange"`   //生效时间段
	InvalidRange *WhenRange   `json:"invalidRange"` //无效时间段(最高优先级)
	Terms        Terms        `json:"terms"`        //条件
}
type WhenRange struct {
	DateRange DateRange `json:"dateRange"`
	TimeRange TimeRange `json:"timeRange"`
	Repeat    Timer     `json:"repeat"`
}

func (w *When) Validate() error {
	if w == nil {
		return nil
	}
	if err := w.ValidRange.Validate(); err != nil {
		return err
	}
	if err := w.Terms.Validate(); err != nil {
		return err
	}
	return nil
}

func (w WhenRange) Validate() error {
	if err := w.TimeRange.Validate(); err != nil {
		return err
	}
	if err := w.Repeat.Validate(); err != nil {
		return err
	}
	return nil
}
