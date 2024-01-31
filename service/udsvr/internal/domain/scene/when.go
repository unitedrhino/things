package scene

import "gitee.com/i-Things/share/errors"

type When struct {
	ValidRanges   WhenRanges `json:"validRanges"`   //生效时间段
	InvalidRanges WhenRanges `json:"invalidRanges"` //无效时间段(最高优先级)
	Conditions    Conditions `json:"conditions"`    //条件
}

type WhenRanges []WhenRange

const (
	WhenRangeTypeDate = "date"
	WhenRangeTypeTime = "time"
)

type WhenRange struct {
	Type      string    `json:"type"` //范围类型 date: 日期范围 time: 时间范围
	DateRange DateRange `json:"dateRange"`
	TimeRange TimeRange `json:"timeRange"`
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

func (w *WhenRange) Validate() error {
	switch w.Type {
	case WhenRangeTypeDate:
		if err := w.DateRange.Validate(); err != nil {
			return err
		}
	case WhenRangeTypeTime:
		if err := w.TimeRange.Validate(); err != nil {
			return err
		}
	default:
		return errors.Parameter.AddMsg("WhenRange type not right")
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
