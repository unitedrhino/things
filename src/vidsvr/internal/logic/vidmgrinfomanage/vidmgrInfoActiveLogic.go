package vidmgrinfomanagelogic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/common"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsvr/internal/types"
	"time"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrInfoActiveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVidmgrInfoActiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrInfoActiveLogic {
	return &VidmgrInfoActiveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 激活服务
func (l *VidmgrInfoActiveLogic) VidmgrInfoActive(in *vid.VidmgrInfoActiveReq) (*vid.Response, error) {
	// todo: add your logic here and delete this line
	//STEP1  需要获取流服务的配置
	fmt.Println("[*****test1*****]", utils.FuncName())
	infoRepo := relationDB.NewVidmgrInfoRepo(l.ctx)
	infoData, err := infoRepo.FindOneByFilter(l.ctx, relationDB.VidmgrFilter{
		VidmgrIDs: []string{in.VidmgrID},
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ManageVidmgr req=%v err=%v", utils.FuncName(), in.VidmgrID, er)
		return nil, er
	}
	//STEP2  修改流媒体服务  //set default
	bytetmp := make([]byte, 0)
	mgr := &clients.SvcZlmedia{
		Secret: infoData.VidmgrSecret,
		Port:   infoData.VidmgrPort,
		IP:     utils.InetNtoA(infoData.VidmgrIpV4),
		//docker服务时
		//Secret: l.svcCtx.Config.Mediakit.Secret,
		//Port:   l.svcCtx.Config.Mediakit.Port,
		//IP:     l.svcCtx.Config.Mediakit.Host,
	}
	if infoData.VidmgrID != "" {
		fmt.Println("[*****test2.6*****]", utils.FuncName())
		mdata, err := clients.ProxyMediaServer(clients.GETSERVERCONFIG, mgr, bytetmp)
		currentConf := new(types.IndexApiServerConfigResp)
		json.Unmarshal(mdata, currentConf)
		fmt.Println("Server Activer getServerConfig:", currentConf)
		//fmt.Println("[*****test2.7*****]", utils.FuncName())
		if err != nil {
			fmt.Println("Server Activer failed")
			return nil, err
		}
		if currentConf.Code != 0 {
			fmt.Println("Server Activer failed")
			return nil, fmt.Errorf("Get MediaServer Config error：%s", currentConf.Msg)
		}
		//STEP3  配置流服务
		var hostIP string
		hostIP = l.svcCtx.Config.Restconf.Host
		common.SetDefaultConfig(hostIP, int64(l.svcCtx.Config.Restconf.Port), &currentConf.Data[0])
		currentConf.Data[0].GeneralMediaServerId = infoData.VidmgrID
		byteConfig, _ := json.Marshal(currentConf.Data[0])
		//STEP3 配置流服务
		respConfig := &types.IndexApiSetServerConfigResp{}
		mdata, err = clients.ProxyMediaServer(clients.SETSERVERCONFIG, mgr, byteConfig)
		json.Unmarshal(mdata, respConfig)
		if err != nil {
			fmt.Println("Server Activer failed")
			return nil, err
		}
		if respConfig.Code != 0 {
			fmt.Println("setServerConfig  配置流服务出错")
			return nil, fmt.Errorf("setServerConfig  配置流服务出错 %s", respConfig.Msg)
		}
		//STEP3  insert配置到数据库
		//STEP3  insert配置到数据库
		confRepo := relationDB.NewVidmgrConfigRepo(l.ctx)

		confRepo.FindOneByFilter(l.ctx, relationDB.VidmgrConfigFilter{
			//VidmgrIPv4: infoData.VidmgrIpV4,
			//VidmgrPort: infoData.VidmgrPort,
			Secret: infoData.VidmgrSecret,
			//VidmgrIDs: []string{vidInfo.VidmgrID},
		})
		if err != nil {
			l.Errorf("%s.Can find vidmgr config err=%v", utils.FuncName(), utils.Fmt(err))
			confRepo.Insert(l.ctx, common.ToVidmgrConfigDB1(&currentConf.Data[0]))
		} else {
			//查询配置数据库，未找到旰做
			confRepo.Insert(l.ctx, common.ToVidmgrConfigDB1(&currentConf.Data[0]))
		}

		//STEP4 更新状态
		//fmt.Println("[*****test4*****]", utils.FuncName())
		if infoData.VidmgrStatus != def.DeviceStatusOnline {
			//UPDATE
			infoData.VidmgrStatus = def.DeviceStatusOnline
			infoData.FirstLogin = time.Now()
			infoData.LastLogin = time.Now()

			err := infoRepo.Update(l.ctx, infoData)
			if err != nil {
				er := errors.Fmt(err)
				l.Errorf("%s.rpc.ManageVidmgr req=%v err=%v", utils.FuncName(), infoData, er)
				return nil, er
			}
			//fmt.Println("[*****test4*****] success", utils.FuncName())
			return &vid.Response{}, nil
		}
		//fmt.Println("[*****test5*****]", utils.FuncName())
	}
	return nil, errors.MediaNotfoundError.AddDetailf("The VidmgrID not found")
}
