package media

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/utils"
	sip "github.com/i-Things/things/src/vidsvr/gosip/sip"
	"github.com/i-Things/things/src/vidsvr/internal/config"
	"strings"
	"sync"
	"time"
)

// todoWFJ 注册一个流媒体，也会跟着注册一个gb服务
type SipServer struct {
	Conf config.Config
	Srv  *sip.Server
}

var (
	ServerOnce sync.Once
	SipSrv     *SipServer
	NotifyMap  map[string]string
)

func NewSipServer(config config.Config) *SipServer {
	ServerOnce.Do(func() {

		SipSrv = &SipServer{
			Conf: config,
			Srv:  sip.NewServer(context.Background()),
		}
		if config.Notify != nil {
			for k, v := range config.Notify {
				if v != "" {
					NotifyMap[strings.ReplaceAll(k, "_", ".")] = v
				}
			}
		}
		_recordList = &sync.Map{}

		uri, _ := sip.ParseSipURI(fmt.Sprintf("sip:%s@%s", config.GbsipConf.Lid, config.GbsipConf.Region))
		_serverDevices = &GbSipDevice{
			DeviceID: config.GbsipConf.Lid,
			Region:   config.GbsipConf.Region,
			addr: &sip.Address{
				DisplayName: sip.String{Str: "sipserver"},
				URI:         &uri,
				Params:      sip.NewParams(),
			},
		}
		_serverDevices.RandomStr = utils.NewSnowFlake(time.Now().Unix())
		SipInfo = &GbSipInfo{
			DID: Config.GbsipConf.Did,
			CID: Config.GbsipConf.Cid,
			LID: Config.GbsipConf.Lid,
		}
		//初始化一个携程监听UDP SIP
		utils.Go(context.Background(), func() {
			SipSrv.Srv.RegistHandler(sip.REGISTER, handlerRegister)
			SipSrv.Srv.RegistHandler(sip.MESSAGE, handlerMessage)

			SipSrv.Srv.ListenUDPServer(fmt.Sprintf("0.0.0.0:%d", config.GbsipConf.UDP))
		})
	})
	//
	fmt.Println("_______airgens read gb28181_____", config.GbsipConf)
	//init Gb28181 sipInfo Data
	return SipSrv
}

func SipPlay() error {

	return nil
}
