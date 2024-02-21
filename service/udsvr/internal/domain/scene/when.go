package scene

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"time"
)

type When struct {
	ValidRanges   WhenRanges `json:"validRanges,omitempty"`   //生效时间段
	InvalidRanges WhenRanges `json:"invalidRanges,omitempty"` //无效时间段(最高优先级)
	Conditions    Conditions `json:"conditions,omitempty"`    //条件 todo(暂不支持)
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

func (w *When) IsHit(ctx context.Context, t time.Time, repo WhenRepo) bool {
	if len(w.InvalidRanges) != 0 {
		if w.InvalidRanges.IsHit(ctx, t, repo) { //禁止的优先级最高
			return false
		}
	}
	if len(w.ValidRanges) != 0 {
		if !w.ValidRanges.IsHit(ctx, t, repo) { //如果没有命中执行范围则返回未命中
			return false
		}
	}
	return true

}

func (w WhenRanges) IsHit(ctx context.Context, t time.Time, repo WhenRepo) bool {
	if len(w) == 0 {
		return true
	}
	for _, v := range w {
		if v.IsHit(ctx, t, repo) {
			return true
		}
	}
	return false
}
func (w *WhenRange) IsHit(ctx context.Context, t time.Time, repo WhenRepo) bool {
	if w.Type == WhenRangeTypeDate {
		return w.DateRange.IsHit(ctx, t, repo)
	}
	return w.TimeRange.IsHit(ctx, t, repo)
}
