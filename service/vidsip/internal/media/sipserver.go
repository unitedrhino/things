package media

import (
	"context"
	"fmt"
	sip2 "github.com/i-Things/things/service/vidsip/gosip/sip"
	"github.com/i-Things/things/service/vidsip/internal/config"
	db "github.com/i-Things/things/service/vidsip/internal/repo/relationDB"
	"github.com/robfig/cron/v3"
	"strings"
	"sync"
)

// todoWFJ 注册一个流媒体
type SipServer struct {
	Conf config.Config
	Srv  *sip2.Server
}

var (
	Ctx           context.Context
	ServerOnce    sync.Once
	SipSrv        *SipServer
	SipInfo       *db.GbSipInfo
	NotifyMap     map[string]string
	activeDevices *ActiveDevices
	recordList    *sync.Map
	ssrcLock      *sync.Mutex
)

func init() {
	fmt.Println("+++++++++++++++vidsip init function+++++++++++++++")
	//LoadSysInfo()
	CronTask()
}

//func LoadSysInfo() {
//
//}

func CronTask() {
	c := cron.New()
	c.Start()
}

func NewSipServer(config config.Config) *SipServer {
	ServerOnce.Do(func() {
		Ctx = context.Background()
		SipSrv = &SipServer{
			Conf: config,
			Srv:  sip2.NewServer(context.Background()),
		}
		if config.Notify != nil {
			for k, v := range config.Notify {
				if v != "" {
					NotifyMap[strings.ReplaceAll(k, "_", ".")] = v
				}
			}
		}
		ssrcLock = &sync.Mutex{}
		recordList = &sync.Map{}
		activeDevices = &ActiveDevices{sync.Map{}}

		SipInfo = &db.GbSipInfo{
			DID:     config.SipConf.Did,
			CID:     config.SipConf.Cid,
			LID:     config.SipConf.Lid,
			SipIp:   config.SipConf.Host,
			SipPort: config.SipConf.Port,
			CNUM:    config.SipConf.Cnum,
			DNUM:    config.SipConf.Dnum,
			NetType: config.SipConf.NetT,
		}
		//初始化一个携程监听UDP SIP
		//utils.Go(context.Background(), func() {
		SipSrv.Srv.RegistHandler(sip2.REGISTER, handlerRegister)
		SipSrv.Srv.RegistHandler(sip2.MESSAGE, handlerMessage)
		//ListenTCPServer
		//go SipSrv.Srv.ListenTCPServer(fmt.Sprintf("%s:%d", SipInfo.SipIp, SipInfo.SipPort))
		go SipSrv.Srv.ListenUDPServer(fmt.Sprintf("%s:%d", SipInfo.SipIp, SipInfo.SipPort))
		//if SipInfo.NetType == "udp" {
		//	go SipSrv.Srv.ListenUDPServer(fmt.Sprintf("%s:%d", SipInfo.SipIp, SipInfo.SipPort))
		//} else {
		//	//has not coding
		//	//go SipSrv.Srv.ListenTCPServer(fmt.Sprintf("%s:%d", SipInfo.SipIp, SipInfo.SipPort))
		//}
		//})
	})
	//
	fmt.Println("_______airgens read gb28181_____", config.SipConf)
	//init Gb28181 sipInfo Data
	return SipSrv
}
