package timer

import (
	"context"
	"github.com/i-Things/things/shared/domain/task"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/timedjobsvr/internal/repo/relationDB"
	"github.com/zeromicro/go-zero/core/logx"
)

func (t Timed) Queue(ctx context.Context, jb *task.Info) error {
	err := t.SvcCtx.PubJob.Publish(ctx, jb.SubType, jb.Queue.Topic, []byte(jb.Queue.Payload))
	e := errors.Fmt(err)
	er := relationDB.NewJobLogRepo(ctx).Insert(ctx, &relationDB.TimedJobLog{
		Group:      jb.Group,
		Type:       jb.Type,
		SubType:    jb.SubType,
		Name:       jb.Name,
		Code:       jb.Code,
		ResultCode: e.GetCode(),
		ResultMsg:  e.GetMsg(),
	})
	if er != nil {
		logx.WithContext(ctx).Errorf("Queue.Publish.Insert err:%v", er)
	}
	return err
}
