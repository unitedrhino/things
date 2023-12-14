package startup

import (
	"context"
	"github.com/i-Things/things/src/vidsvr/internal/event/serverEvent"
	"github.com/i-Things/things/src/vidsvr/internal/repo/event/server"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

func Subscribe(svcCtx *svc.ServiceContext) {
	{
		cli, err := server.NewServer(svcCtx.Config.Event)
		logx.Must(err)
		err = cli.Subscribe(func(ctx context.Context) server.ServerHandle {
			return serverEvent.NewServerHandle(ctx, svcCtx)
		})

		logx.Must(err)
	}
}
