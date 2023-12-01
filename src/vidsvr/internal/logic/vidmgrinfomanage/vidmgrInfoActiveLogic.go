package vidmgrinfomanagelogic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/logic/zlmedia/index"
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
	fmt.Println("[*****test2*****]", utils.FuncName())
	fmt.Println("[*****test2*****]", utils.FuncName(), infoData)
	//STEP2  修改流媒体服务  //set default
	if infoData.VidmgrID != "" {
		fmt.Println("[*****test2.5*****]", utils.FuncName())
		getServerConfig := index.NewGetServerConfigLogic(l.ctx, l.svcCtx)
		setServerConfig := index.NewSetServerConfigLogic(l.ctx, l.svcCtx)
		fmt.Println("[*****test2.6*****]", utils.FuncName())
		currentConf, error := getServerConfig.GetServerConfig(&types.IndexApiReq{
			VidmgrID: infoData.VidmgrID,
			Data:     "",
		})
		fmt.Println("[*****test2.7*****]", utils.FuncName())
		if error != nil {
			fmt.Println("Server Activer failed")
			return nil, error
		}
		fmt.Println("dataConfig:", currentConf)
		fmt.Println("[*****test3*****]", utils.FuncName())
		//STEP3  配置流服务
		if len(currentConf.Data) > 0 {
			if infoData.MediasvrType == 1 { //docker模式
				index.SetDefaultConfig(l.svcCtx.Config.Mediakit.Host, int64(l.svcCtx.Config.Restconf.Port), &currentConf.Data[0])
			} else { //独立主机模式
				index.SetDefaultConfig(utils.InetNtoA(infoData.ServerIP), infoData.ServerPort, &currentConf.Data[0])
			}
			currentConf.Data[0].GeneralMediaServerId = infoData.VidmgrID
			byteConfig, _ := json.Marshal(currentConf.Data[0])
			//STEP3 配置流服务
			configReq := &types.IndexApiSetServerConfigReq{
				VidmgrID: infoData.VidmgrID,
				Data:     string(byteConfig),
			}
			setConfig, error := setServerConfig.SetServerConfig(configReq)
			if setConfig.Code != 0 {
				fmt.Println("setServerConfig  配置流服务出错")
				return nil, error
			}
			//STEP3  insert配置到数据库
			fmt.Println("[*****test4*****]", utils.FuncName())
			confRepo := relationDB.NewVidmgrConfigRepo(l.ctx)
			confRepo.Insert(l.ctx, ToVidmgrConfigRpc(&currentConf.Data[0]))
			//STEP4 更新状态
			fmt.Println("[*****test4*****]", utils.FuncName())
			if infoData.VidmgrStatus != def.DeviceStatusOnline {
				//UPDATE
				infoData.VidmgrStatus = def.DeviceStatusOnline
				infoData.FirstLogin.Time = time.Now()
				infoData.FirstLogin.Valid = true
				infoData.LastLogin.Time = time.Now()
				infoData.LastLogin.Valid = true

				err := infoRepo.Insert(l.ctx, infoData)
				if err != nil {
					er := errors.Fmt(err)
					l.Errorf("%s.rpc.ManageVidmgr req=%v err=%v", utils.FuncName(), infoData, er)
					return nil, er
				}
				fmt.Println("[*****test4*****] success", utils.FuncName())
				return &vid.Response{}, nil
			}
		}
		fmt.Println("[*****test5*****]", utils.FuncName())
	}
	return nil, errors.MediaNotfoundError.AddDetailf("The VidmgrID not found")
}
