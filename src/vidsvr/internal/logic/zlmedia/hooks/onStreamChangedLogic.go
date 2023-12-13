package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"regexp"
	"time"
)

const StreamUpTimeOut = 2 * time.Second

type OnStreamChangedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	//cache  cache.Cache
}

//var (
//	// can't use one SharedCalls per conn, because multiple conns may share the same cache key.
//	exclusiveCalls = syncx.NewSingleFlight()
//	stats          = cache.NewStat("streamupdate")
//	cachePrefix    = "streamupdate"
//)

func NewOnStreamChangedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnStreamChangedLogic {
	return &OnStreamChangedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		//cache:  cache.New(svcCtx.Config.CacheRedis, exclusiveCalls, stats, errors.NotFind),
	}
}

// [fix-bug by wfj]   不同IP推上来的相同路径的流，可能会有问题
func (l *OnStreamChangedLogic) OnStreamChanged(req *types.HooksApiStreamChangedRep) (resp *types.HooksApiResp, err error) {
	// todo: add your logic here and delete this line
	reqStr, _ := json.Marshal(*req)
	fmt.Println("---------OnStreamChanged--------------:", string(reqStr))

	//需要先判断该流服务是否有注册过，未注册过忽略消息
	infoRepo := relationDB.NewVidmgrInfoRepo(l.ctx)
	vidInfo, err := infoRepo.FindOneByFilter(l.ctx, relationDB.VidmgrFilter{
		VidmgrIDs: []string{req.MediaServerId},
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s rpc.ManageVidmgr req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	if vidInfo != nil {
		//查找要素：vidmgr_id  app  stream    peerIP
		streamRepo := relationDB.NewVidmgrStreamRepo(l.ctx)
		vidStreamInfo, err := streamRepo.FindOneByFilter(l.ctx, relationDB.VidmgrStreamFilter{
			VidmgrID:   req.MediaServerId,
			App:        req.App,
			Stream:     req.Stream,
			OriginType: req.OriginType,
			PeerIP:     utils.InetAtoN(req.OriginSock.PeerIp),
		})
		//未找到流信息
		if err != nil {
			//如何未查询到插入该流
			erros := &types.IndexApiResp{}
			json.Unmarshal([]byte(err.Error()), erros)
			//未找到记录和注册回调同时满足时登录该流
			if req.Regist && erros.Code == 100009 {
				vidStreamInfo = ToVidmgrStreamRpc(req)
				vidStreamInfo.IsOnline = true //设置状态为在线
				vidStreamInfo.FirstLogin = time.Now()
				vidStreamInfo.LastLogin = time.Now()
				//不关心流类型了
				//SetProtocol(req.Schema, vidStreamInfo)
				if vidStreamInfo.OriginType == RTMP_PUSH {
					re := regexp.MustCompile(vidStreamInfo.Vhost)
					if vidInfo.MediasvrType == 1 { //docker 模式
						//l.svcCtx.Config.Restconf.Host
						vidStreamInfo.OriginUrl =
							re.ReplaceAllString(vidStreamInfo.OriginUrl, l.svcCtx.Config.Restconf.Host)
					} else {
						vidStreamInfo.OriginUrl =
							re.ReplaceAllString(vidStreamInfo.OriginUrl, req.OriginSock.PeerIp)
					}
				}
				err := streamRepo.Insert(l.ctx, vidStreamInfo)
				if err != nil {
					l.Errorf("%s rpc.OnStreamChanged  err=%+v", utils.FuncName(), err)
					return nil, err
				}
			} else { //ignore message
				l.Errorf("ignore req=%v err=%+v", utils.FuncName(), err)
				return nil, err
			}
		} else { //找到了一条流就直需要修改状态就可以了
			//判断Sock相同为同一流  	update
			if req.Regist {
				//对应位(bit)  置1
				//SetProtocol(req.Schema, vidStreamIndex)
				vidStreamInfo.IsOnline = true
				vidStreamInfo.LastLogin = time.Now()
			} else {
				//对应位(bit) 置0
				vidStreamInfo.IsOnline = false
				vidStreamInfo.LastLogin = time.Now()
				//UnSetProtocol(req.Schema, vidStreamInfo)
			}
			err := streamRepo.Update(l.ctx, vidStreamInfo)
			if err != nil {
				l.Errorf("%s rpc.VidmgrStreamUpdate  err=%+v", utils.FuncName(), err)
				return nil, err
			}
		}
	}

	return &types.HooksApiResp{
		Code: 0,
		Msg:  "success",
	}, nil
}
