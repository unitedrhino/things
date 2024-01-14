package scene

type When struct {
	ValidRanges   WhenRanges `json:"validRanges"`   //生效时间段
	InvalidRanges WhenRanges `json:"invalidRanges"` //无效时间段(最高优先级)
	Conditions    Conditions `json:"conditions"`    //条件
}

type WhenRanges []WhenRange

type WhenRange struct {
	DateRange DateRange `json:"dateRange"`
	TimeRange TimeRange `json:"timeRange"`
	Repeat    Timer     `json:"repeat"`
}

func (w *When) Validate() error {
	if w == nil {
		return nil
	}
	if err := w.ValidRanges.Validate(); err != nil {
		return err
	}
	if err := w.InvalidRanges.Validate(); err != nil {
		return err
	}
	if err := w.Conditions.Validate(); err != nil {
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
func (w WhenRanges) Validate() error {
	if len(w) == 0 {
		return nil
	}
	for _, v := range w {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	return nil
}
