package media

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/utils"
	sip "github.com/i-Things/things/src/vidsvr/gosip/sip"
	"github.com/i-Things/things/src/vidsvr/internal/config"
	db "github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"strings"
	"sync"
)

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
		uri, _ := sip.ParseSipURI(fmt.Sprintf("sip:%s@%s", config.GbsipConf.Lid, config.GbsipConf.Region))
		_serverDevices = &db.VidmgrDevices{
			DeviceID: config.GbsipConf.Lid,
			Region:   config.GbsipConf.Region,
			Addr: &sip.Address{
				DisplayName: sip.String{Str: "sipserver"},
				URI:         &uri,
				Params:      sip.NewParams(),
			},
		}

		SipInfo = &db.VidmgrSipInfo{
			DID: Config.GbsipConf.Did,
			CID: Config.GbsipConf.Cid,
		}

		//初始化一个携程监听UDP SIP
		utils.Go(context.Background(), func() {
			SipSrv.Srv.RegistHandler(sip.REGISTER, handlerRegister)
			SipSrv.Srv.RegistHandler(sip.MESSAGE, handlerMessage)
			SipSrv.Srv.ListenUDPServer(config.GbsipConf.UDP)
		})
	})

	//
	fmt.Println("_______airgens read gb28181_____", config.GbsipConf)
	//init Gb28181 sipInfo Data
	sipInfoRepo := db.NewVidmgrSipInfoRepo(Ctx)
	filter := db.VidmgrSipInfoFilter{}
	size, err := sipInfoRepo.CountByFilter(Ctx, filter)
	if err != nil {
		fmt.Sprintf("---NewSipServer ERROR----has no data")
	}
	if size > 0 {
		diInfo, _ := sipInfoRepo.FindByFilter(Ctx, filter, &def.PageInfo{
			Page: 1, Size: 20,
			Orders: []def.OrderBy{{"created_time", def.OrderDesc}, {"id", def.OrderDesc}},
		})
		//GetData  取第一个组数据
		// = diInfo[0].DID
		SipInfo = diInfo[0]
	} else {
		//SipInfo = config.GbsipConf
		//insert Data
		SipInfo.Region = config.GbsipConf.Region
		SipInfo.LID = config.GbsipConf.Lid
		SipInfo.DID = config.GbsipConf.Did
		SipInfo.CID = config.GbsipConf.Cid
		SipInfo.IsOpenServer = true
		SipInfo.MediaServerRtpIP = utils.InetAtoN(config.Restconf.Host)
		SipInfo.MediaServerRtpPort = 10000
		sipInfoRepo.Insert(Ctx, SipInfo)
	}

	return SipSrv
}
