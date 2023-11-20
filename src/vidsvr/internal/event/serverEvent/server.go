package serverEvent

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type ServerHandle struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
	PiDB *relationDB.VidmgrInfoRepo
}

func NewServerHandle(ctx context.Context, svcCtx *svc.ServiceContext) *ServerHandle {
	return &ServerHandle{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		PiDB:   relationDB.NewVidmgrInfoRepo(ctx),
	}
}

func (l *ServerHandle) ActionCheck() error {
	//l.Infof("ActionCheck req:%v", in)
	fmt.Println("[****] func (l *ServerHandle) ActionCheck() error ")
	//需要做的操作，查旬数据库
	now := time.Now().Unix()
	//过滤条件为：在线设备且超时时间为60秒
	filter := relationDB.VidmgrFilter{LastLoginTime: struct {
		Start int64
		End   int64
	}{Start: 0, End: now - clients.VIDMGRTIMEOUT}, VidmgrStatus: def.DeviceStatusOnline}
	di, err := l.PiDB.FindAllFilter(l.ctx, filter)
	if err != nil {
		return err
	}
	if len(di) > 0 {
		for _, v := range di {
			v.VidmgrStatus = def.DeviceStatusOffline
			l.PiDB.Update(l.ctx, v) //更新数据库
		}
	} else {
		//do nothing
	}
	//判断当前时间与最后login时间，是否超过30s
	//1分钟会执行一次
	return nil
}
