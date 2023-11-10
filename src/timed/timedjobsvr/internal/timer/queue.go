package timer

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/timed/internal/domain"
	"github.com/i-Things/things/src/timed/internal/repo/relationDB"
	"github.com/zeromicro/go-zero/core/logx"
)

func (t Timed) Queue(ctx context.Context, task *domain.TaskInfo) error {
	err := t.SvcCtx.PubJob.Publish(ctx, task.GroupSubType, task.Queue.Topic, []byte(task.Queue.Payload))
	e := errors.Fmt(err)
	er := relationDB.NewJobLogRepo(ctx).Insert(ctx, &relationDB.TimedTaskLog{
		GroupCode:  task.GroupCode,
		TaskCode:   task.Code,
		Params:     task.Params,
		ResultCode: e.GetCode(),
		ResultMsg:  e.GetMsg(),
	})
	if er != nil {
		logx.WithContext(ctx).Errorf("Queue.Publish.Insert err:%v", er)
	}
	return err
}
