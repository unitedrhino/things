package vidmgrstreammanagelogic

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/i-Things/core/shared/clients"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/vidsvr/internal/common"
	"github.com/i-Things/things/service/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/vidsvr/internal/types"
	"time"

	"github.com/i-Things/things/service/vidsvr/internal/svc"
	"github.com/i-Things/things/service/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrStreamCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.VidmgrStreamRepo
}

func NewVidmgrStreamCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrStreamCreateLogic {
	return &VidmgrStreamCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewVidmgrStreamRepo(ctx),
	}
}

// 流添加 拉流添加接口
func (l *VidmgrStreamCreateLogic) VidmgrStreamCreate(in *vid.VidmgrStreamCreateReq) (*vid.Response, error) {
	// todo: add your logic here and delete this line
	//查询流服务
	infoRepo := relationDB.NewVidmgrInfoRepo(l.ctx)
	infoData, err := infoRepo.FindOneByFilter(l.ctx, relationDB.VidmgrFilter{
		VidmgrIDs: []string{in.StreamInfo.VidmgrID},
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ManageVidmgr req=%v err=%v", utils.FuncName(), in.StreamInfo.VidmgrID, er)
		return nil, er
	}
	//配置流服务 返回Key，由Key可以回删流  将Key写入数据库
	//默认数据初始化成独立主机流服务
	mgr := &clients.SvcZlmedia{
		Secret: infoData.VidmgrSecret,
		Port:   infoData.VidmgrPort,
		IP:     utils.InetNtoA(infoData.VidmgrIpV4),
	}

	if infoData.VidmgrID != "" {
		//docker模式的连接方式
		if infoData.MediasvrType == 1 {
			mgr.Secret = l.svcCtx.Config.Mediakit.Secret
			mgr.Port = l.svcCtx.Config.Mediakit.Port
			mgr.IP = l.svcCtx.Config.Mediakit.Host
		}
		//初始化发送的数据
		zlmData := &types.IndexApiAddStreamProxy{
			Stream:  in.StreamInfo.Stream,
			Vhost:   in.StreamInfo.Vhost,
			App:     in.StreamInfo.App,
			Url:     in.StreamInfo.OriginUrl,
			RtpType: in.RtpType,
		}
		bytetmp, _ := json.Marshal(zlmData)
		//发送添加流数据
		mdata, err := clients.ProxyMediaServer(clients.ADDSTREAMPROXY, mgr, bytetmp)
		if err != nil {
			fmt.Println("Server Activer failed")
			return nil, err
		}
		zlmResp := new(types.IndexApiAddStreamProxyResp)
		json.Unmarshal(mdata, zlmResp)
		fmt.Println("Server Activer IndexApiAddStreamProxyResp:", zlmResp)
		if zlmResp.Code != 0 {
			return nil, fmt.Errorf("AddStreamProxy:%s", zlmResp.Msg)
		}
		//PullKey := zlmResp.Data.Key
		//PullUrl := in.OriginUrl
		//插入数据库
		streamRepo := relationDB.NewVidmgrStreamRepo(l.ctx)
		streamData := common.ToVidmgrStreamDB(in.StreamInfo)
		streamData.PullKey = zlmResp.Data.Key
		streamData.RtpType = in.RtpType
		streamData.OriginType = clients.PULL
		streamData.FirstLogin = time.Now()

		//filter
		filter := &relationDB.VidmgrStreamFilter{
			VidmgrID:   streamData.VidmgrID,
			App:        streamData.App,
			Stream:     streamData.Stream,
			OriginType: streamData.OriginType,
			OriginUrl:  streamData.OriginUrl,
		}
		vidStreamInfo, err := streamRepo.FindOneByFilter(l.ctx, *filter)
		if vidStreamInfo == nil {
			streamRepo.Insert(l.ctx, streamData)
		} else {
			return nil, errors.MediaPullCreateError.AddDetail("流数据已经存在")
		}
	}
	return &vid.Response{}, nil
}
