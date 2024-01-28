package scene

import (
	"testing"
	"time"
)

func TestTimeRange_IsHit(t1 *testing.T) {
	var st = "123"
	var aaa = "423425"
	if st > aaa {

	}

	type fields struct {
		Type string
		Cron string
	}
	type args struct {
		tim time.Time
	}
	now := time.Now()
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			fields: fields{Cron: "5 5 * * * *"}, want: false,
			args: args{tim: time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 5, now.Nanosecond(), time.Local)},
		},
		{
			fields: fields{Cron: "* 5 1 2 3 *"}, want: true,
			args: args{tim: time.Date(now.Year(), 3, 2, 1, 5, now.Second(), now.Nanosecond(), time.Local)},
		},
		{
			fields: fields{Cron: "0 0 12 ? * WED"}, want: true, //表示每个星期三中午12点
			args: args{tim: time.Date(2023, 3, 29, 12, 0, 0, now.Nanosecond(), time.Local)},
		},
		{
			fields: fields{Cron: "* 10,44 14 ? 3 WED"}, want: true, //每年三月的星期三的下午2:10和2:44触发
			args: args{tim: time.Date(now.Year(), 3, 29, 14, 10, now.Second(), now.Nanosecond(), time.Local)},
		},
		{
			fields: fields{Cron: "* 10,44 14 ? 3 WED"}, want: false, //每年三月的星期三的下午2:10和2:44触发
			args: args{tim: time.Date(now.Year(), 4, 29, 14, 10, now.Second(), now.Nanosecond(), time.Local)},
		},
		{
			fields: fields{Cron: "*/2 10,44 14 ? 3 WED"}, want: false, //每年三月的星期三的下午2:10和2:44触发
			args: args{tim: time.Date(now.Year(), 4, 29, 14, 10, now.Second(), now.Nanosecond(), time.Local)},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TimeRange{
				Type: tt.fields.Type,
			}
			if got := t.IsHit(tt.args.tim); got != tt.want {
				t1.Errorf("IsHit() = %v, want %v", got, tt.want)
			}
		})
	}
}
