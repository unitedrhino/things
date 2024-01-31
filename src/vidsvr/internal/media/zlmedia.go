package media

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/common"
	"github.com/i-Things/things/src/vidsvr/internal/config"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsvr/internal/types"
	"github.com/i-Things/things/src/vidsvr/pb/vid"
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

func SrvInfoStatusCheck() error {
	//l.Infof("ActionCheck req:%v", in)
	fmt.Println("[****] func (l *ServerHandle) ActionCheck() error ")
	//需要做的操作，查旬数据库
	now := time.Now().Unix()
	//过滤条件为：在线设备且超时时间为60秒
	filter := relationDB.VidmgrFilter{LastLoginTime: struct {
		Start int64
		End   int64
	}{Start: 0, End: now - clients.VIDMGRTIMEOUT}, VidmgrStatus: def.DeviceStatusOnline}
	InfoRepo := relationDB.NewVidmgrInfoRepo(Ctx)
	di, err := InfoRepo.FindAllFilter(Ctx, filter)
	if err != nil {
		return err
	}
	if len(di) > 0 {
		for _, v := range di {
			v.VidmgrStatus = def.DeviceStatusOffline
			InfoRepo.Update(Ctx, v) //更新数据库
		}
	} else {
		//do nothing
	}
	//判断当前时间与最后login时间，是否超过30s
	//1分钟会执行一次
	return nil
}

func InitDockerSrv(c config.Config, vidmgrID string) error {
	//l.Infof("ActionCheck req:%v", in)
	var vidInfo *relationDB.VidmgrInfo
	fmt.Println("[**ActionInit**]0 ", utils.FuncName())
	//查找流服务的数据库：根据IP和端口确定一个流服务
	var (
		filter = relationDB.VidmgrFilter{VidmgrIpV4: utils.InetAtoN(c.Mediakit.Host), VidmgrPort: c.Mediakit.Port}
	)
	InfoRepo := relationDB.NewVidmgrInfoRepo(Ctx)
	//fmt.Println("[**ActionInit**]1 ", utils.FuncName())
	size, err := InfoRepo.CountByFilter(Ctx, filter)
	if err != nil {
		fmt.Errorf("MediaServer init data countfilter error")
		return err
	}
	//找到存在一条流服务 更新这条服务
	if size > 0 {
		//update
		fmt.Println("[**ActionInit**]3 ", utils.FuncName())
		page := vid.PageInfo{}

		di, err := InfoRepo.FindByFilter(Ctx, filter, common.ToPageInfoWithDefault(&page, &def.PageInfo{
			Page: 1, Size: 20,
			Orders: []def.OrderBy{{"created_time", def.OrderDesc}, {"vidmgr_id", def.OrderDesc}},
		}))
		if err != nil {
			fmt.Errorf("MediaServer init data find filter error")
			return err
		}
		if di[0].VidmgrSecret != c.Mediakit.Secret {
			di[0].VidmgrSecret = c.Mediakit.Secret
			err = InfoRepo.Update(Ctx, di[0])
		}
		vidInfo = di[0]
	} else {
		//流服务还未存在的情况就插入这条服务
		fmt.Println("[**ActionInit**]4 ", utils.FuncName())
		dbDocker := &relationDB.VidmgrInfo{
			VidmgrID:     vidmgrID,
			VidmgrName:   "default Docker",
			VidmgrIpV4:   utils.InetAtoN(c.Mediakit.Host),
			VidmgrPort:   c.Mediakit.Port,
			VidmgrSecret: c.Mediakit.Secret,
			RtpPort:      10000,
			IsOpenGbSip:  true,
			VidmgrStatus: 2, //默认设置离线状态
			VidmgrType:   1, //ZLmediakit
			MediasvrType: 1, //docker模式
			Desc:         "",
			Tags:         map[string]string{},
		}

		err = InfoRepo.Insert(Ctx, dbDocker)
		if err != nil {
			fmt.Printf("%s.Insert err=%+v", utils.FuncName(), err)
			return err
		}
		vidInfo = dbDocker
		fmt.Println("[**ActionInit**]5 ", utils.FuncName())
	}
	bytetmp := make([]byte, 0)
	//vidsvr->zlmediakit
	mgr := &clients.SvcZlmedia{
		Secret: c.Mediakit.Secret,
		Port:   c.Mediakit.Port,
		IP:     c.Mediakit.Host,
	}
	mdata, err := clients.ProxyMediaServer(clients.GETSERVERCONFIG, mgr, bytetmp)
	currentConf := new(types.IndexApiServerConfigResp)
	json.Unmarshal(mdata, currentConf)

	//config dockerServer
	fmt.Println("[**ActionInit**]6 ", utils.FuncName())
	fmt.Println("[**ActionInit**]6.1 ProxyMediaServer: ", string(mdata))
	//仅考虑docker的模式
	//STEP3  配置流服务
	if len(currentConf.Data) > 0 {
		currentConf.Data[0].GeneralMediaServerId = vidInfo.VidmgrID
		//docker通信IP用eth0 从zlmediakit->vidsvr
		common.SetDefaultConfig(c.Restconf.Host, int64(c.Restconf.Port), &currentConf.Data[0])
		fmt.Println("[****setting****] ", c.Mediakit.Host, int64(c.Restconf.Port))
		byteConfig, _ := json.Marshal(currentConf.Data[0])
		//STEP3 配置流服务
		mdata, err = clients.ProxyMediaServer(clients.SETSERVERCONFIG, mgr, byteConfig)
		dataRecv := new(types.IndexApiSetServerConfigResp)
		err = json.Unmarshal(mdata, dataRecv)
		if err != nil {
			fmt.Println("parse Json failed:", err)
			return err
		}
		//STEP3  insert配置到数据库
		fmt.Println("[*****test6*****]", utils.FuncName())
		confRepo := relationDB.NewVidmgrConfigRepo(Ctx)
		//查找config配置
		confRepo.FindOneByFilter(Ctx, relationDB.VidmgrConfigFilter{
			//VidmgrIPv4: vidInfo.VidmgrIpV4,
			//VidmgrPort: vidInfo.VidmgrPort,
			Secret: vidInfo.VidmgrSecret,
			//VidmgrIDs: []string{vidInfo.VidmgrID},
		})
		if err != nil {
			fmt.Errorf("%s.Can find vidmgr config err=%v", utils.FuncName(), utils.Fmt(err))
			confRepo.Insert(Ctx, common.ToVidmgrConfigDB1(&currentConf.Data[0]))
		} else {
			//update
			confRepo.Update(Ctx, common.ToVidmgrConfigDB1(&currentConf.Data[0]))
		}

		//STEP4 更新状态
		fmt.Println("[*****test7*****]", utils.FuncName())
		if vidInfo.VidmgrStatus != def.DeviceStatusOnline {
			//UPDATE
			vidInfo.VidmgrStatus = def.DeviceStatusOnline
			vidInfo.FirstLogin = time.Now()
			vidInfo.LastLogin = time.Now()

			err := InfoRepo.Update(Ctx, vidInfo)
			if err != nil {
				er := errors.Fmt(err)
				fmt.Errorf("%s.rpc.ManageVidmgr req=%v err=%v", utils.FuncName(), vidInfo, er)
				return er
			}
			fmt.Println("[*****test8*****] success", utils.FuncName())
			return nil
		}
	}
	fmt.Println("[*****test9*****]error:配置错误", utils.FuncName())
	return nil
}
