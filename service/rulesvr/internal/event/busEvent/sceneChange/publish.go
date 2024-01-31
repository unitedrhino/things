package sceneChange

import (
	"context"
	"gitee.com/i-Things/core/shared/events"
	"github.com/i-Things/things/service/rulesvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

func EventsHandle(svcCtx *svc.ServiceContext, topic string) any {
	return func(ctx context.Context, id int64) {
		err := svcCtx.DataUpdate.UpdateWithTopic(
			ctx, topic, &events.ChangeInfo{ID: id})
		if err != nil {
			logx.WithContext(ctx).Errorf("EventsHandle id:%v err:%v", id, err)
		}
		logx.WithContext(ctx).Infof("EventsHandle id:%v", id)
	}
}
