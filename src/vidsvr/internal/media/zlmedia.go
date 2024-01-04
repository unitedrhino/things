package media

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/common"
	"github.com/i-Things/things/src/vidsvr/internal/config"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsvr/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"regexp"
	"sync"
	"time"
)

// 用于流改变状态记录的结构
type LastStreamInfo struct {
	LoginTime time.Time
	Req       types.HooksApiStreamChangedRep
}

type MediaChan struct {
	logx.Logger
	ChangeStream chan LastStreamInfo //channel
}

var (
	InitOnce     sync.Once
	SvcMediaChan *MediaChan
	Config       config.Config
	//保存上一次请求
	Old LastStreamInfo
	Ctx context.Context
)

func NewMediaChan(config config.Config) *MediaChan {
	InitOnce.Do(func() {
		Config = config
		Ctx := context.Background()
		Old = LastStreamInfo{}
		SvcMediaChan = &MediaChan{
			Logger:       logx.WithContext(Ctx),
			ChangeStream: make(chan LastStreamInfo, clients.STREAMCHANGESIZE),
		}
		//初始化一个协程
		utils.Go(Ctx, MediaThreadRun) //后台执行
	})
	return SvcMediaChan
}

func GetMediaChan() *MediaChan {
	return SvcMediaChan
}

// // thread
func MediaThreadRun() {
	for {
		select {
		//查询执行函数
		case v := <-SvcMediaChan.ChangeStream:
			fmt.Println("-------取出数据 start---------:", v)
			QueueOnChangeStream(v)
			fmt.Println("-------取出数据 end  ---------:", v)
		default:
			//fmt.Println("*********************MediaThreadRun wait-300ms************************")
			time.Sleep(300 * time.Microsecond)
		}
	}
}

func CheckTimeSecond(old time.Time, now time.Time) int {
	return now.Second() - old.Second()
}

func QueueOnChangeStream(qreq LastStreamInfo) {
	//过滤掉多多余的数据[不解析协议类注册]
	if CheckTimeSecond(qreq.LoginTime, time.Now()) < 2 &&
		Old.Req.Regist == qreq.Req.Regist &&
		Old.Req.OriginSock.Identifier == qreq.Req.OriginSock.Identifier &&
		Old.Req.MediaServerId == qreq.Req.MediaServerId {
		return
	}
	Old = qreq

	//先查询流服务器
	infoRepo := relationDB.NewVidmgrInfoRepo(Ctx)
	vidInfo, err := infoRepo.FindOneByFilter(Ctx, relationDB.VidmgrFilter{
		VidmgrIDs: []string{qreq.Req.MediaServerId},
	})
	if err != nil {
		er := errors.Fmt(err)
		SvcMediaChan.Errorf("%s rpc.ManageVidmgr req=%v err=%+v", utils.FuncName(), qreq.Req, er)
		return
	}
	if vidInfo != nil {
		//查找要素：vidmgr_id  app  stream    peerIP OriginUrl(push时依据)
		streamRepo := relationDB.NewVidmgrStreamRepo(Ctx)
		filter := &relationDB.VidmgrStreamFilter{
			VidmgrID:   qreq.Req.MediaServerId,
			App:        qreq.Req.App,
			Stream:     qreq.Req.Stream,
			OriginType: qreq.Req.OriginType,
		}
		if qreq.Req.OriginType == clients.PULL {
			filter.OriginUrl = qreq.Req.OriginUrl
		} else { //Push stream
			filter.PeerIP = utils.InetAtoN(qreq.Req.OriginSock.PeerIp)
		}

		vidStreamInfo, err := streamRepo.FindOneByFilter(Ctx, *filter)
		//未找到流信息
		if err != nil {
			//如何未查询到插入该流
			erros := &types.IndexApiResp{}
			json.Unmarshal([]byte(err.Error()), erros)
			//未找到记录和注册回调同时满足时登录该流
			if qreq.Req.Regist && erros.Code == 100009 {
				vidStreamInfo = common.ToVidmgrStreamRpc1(&qreq.Req)
				vidStreamInfo.IsOnline = true //设置状态为在线
				vidStreamInfo.FirstLogin = time.Now()
				vidStreamInfo.LastLogin = time.Now()
				//根据流类型，确定
				if vidStreamInfo.OriginType == clients.RTMP_PUSH || vidStreamInfo.OriginType == clients.RTSP_PUSH ||
					vidStreamInfo.OriginType == clients.RTP_PUSH {
					re := regexp.MustCompile(vidStreamInfo.Vhost)
					if vidInfo.MediasvrType == 1 { //docker 模式
						vidStreamInfo.OriginUrl =
							re.ReplaceAllString(vidStreamInfo.OriginUrl, Config.Restconf.Host)
					} else {
						//LocalIP为流服务IP，PeerIP为推流源地址
						vidStreamInfo.OriginUrl =
							re.ReplaceAllString(vidStreamInfo.OriginUrl, qreq.Req.OriginSock.LocalIp)
					}
				}
				//如果是拉流方式，不需要修改
				err := streamRepo.Insert(Ctx, vidStreamInfo)
				if err != nil {
					SvcMediaChan.Errorf("%s rpc.OnStreamChanged  err=%+v", utils.FuncName(), err)
					return
				}
			} else { //ignore message
				SvcMediaChan.Errorf("ignore req=%v err=%+v", utils.FuncName(), err)
				return
			}
		} else { //找到了一条流就直需要修改状态就可以了
			if qreq.Req.Regist {
				vidStreamInfo.IsOnline = true
				vidStreamInfo.LastLogin = time.Now()
			} else {
				vidStreamInfo.IsOnline = false
				vidStreamInfo.LastLogin = time.Now()
			}
			err := streamRepo.Update(Ctx, vidStreamInfo)
			if err != nil {
				SvcMediaChan.Errorf("%s rpc.VidmgrStreamUpdate  err=%+v", utils.FuncName(), err)
				return
			}
		}
	}
	return
}
