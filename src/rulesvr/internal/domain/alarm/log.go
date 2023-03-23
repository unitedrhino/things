package alarm

import (
	"context"
	"github.com/i-Things/things/shared/def"
)

type LogFilter struct {
	Time def.TimeRange
}

type LogRepo interface {
	Create(ctx context.Context, dto *LogCreateDto) error
}
type LogCreateDto struct {
	AlarmID   int64  //告警记录ID
	Serial    string //告警流水
	SceneName string //场景名称
	SceneID   int64  //场景ID
	Desc      string //告警说明
}
