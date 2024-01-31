package vidmgrstreammanagelogic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsvr/internal/types"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrStreamDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVidmgrStreamDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrStreamDeleteLogic {
	return &VidmgrStreamDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除流
func (l *VidmgrStreamDeleteLogic) VidmgrStreamDelete(in *vid.VidmgrStreamDeleteReq) (*vid.Response, error) {
	// todo: add your logic here and delete this line
	//如果是Pull流，需要删除ZLMediakit配置的流

	streamRepo := relationDB.NewVidmgrStreamRepo(l.ctx)
	filter := &relationDB.VidmgrStreamFilter{
		StreamIDs: []int64{in.StreamID},
	}
	//先查询流 判断流类型
	vidStreamInfo, err := streamRepo.FindOneByFilter(l.ctx, *filter)
	if vidStreamInfo == nil {
		return nil, errors.MediaStreamDeleteError.AddDetail("流数据不存在")
	}
	infoRepo := relationDB.NewVidmgrInfoRepo(l.ctx)
	infoData, err := infoRepo.FindOneByFilter(l.ctx, relationDB.VidmgrFilter{
		VidmgrIDs: []string{vidStreamInfo.VidmgrID},
	})
	if infoData == nil {
		return nil, errors.MediaStreamDeleteError.AddDetail("流服务不存在")
	}
	mgr := &clients.SvcZlmedia{
		Secret: infoData.VidmgrSecret,
		Port:   infoData.VidmgrPort,
		IP:     utils.InetNtoA(infoData.VidmgrIpV4),
	}
	//需要流服务同时在线和是拉流类型才会去删除ZLMediakit的配置
	if vidStreamInfo.OriginType == clients.PULL && infoData.VidmgrStatus == def.DeviceStatusOnline {
		//拉流需要先删除ZLMediakit的配置
		zlmData := &types.IndexApiDelStreamProxy{
			Key: vidStreamInfo.PullKey,
		}
		//docker模式的连接方式
		if infoData.MediasvrType == 1 {
			mgr.Secret = l.svcCtx.Config.Mediakit.Secret
			mgr.Port = l.svcCtx.Config.Mediakit.Port
			mgr.IP = l.svcCtx.Config.Mediakit.Host
		}
		bytetmp, _ := json.Marshal(zlmData)
		mdata, _ := clients.ProxyMediaServer(clients.DELSTREAMPROXY, mgr, bytetmp)
		zlmResp := new(types.IndexApiDelStreamProxyResp)
		json.Unmarshal(mdata, zlmResp)
		fmt.Println("DeleteStream ZlmResp:", zlmResp)
		//if zlmResp.Code == 0
	}
	err = streamRepo.DeleteByFilter(l.ctx, *filter)
	if err != nil {
		l.Errorf("%s.Delete err=%v", utils.FuncName(), utils.Fmt(err))
		return nil, err
	}
	return &vid.Response{}, nil
}
