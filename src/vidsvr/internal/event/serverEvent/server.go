package serverEvent

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/domain/deviceAuth"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/common"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"
	"github.com/i-Things/things/src/vidsvr/pb/vid"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type ServerHandle struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
	PiDB *relationDB.VidmgrInfoRepo
}

func NewServerHandle(ctx context.Context, svcCtx *svc.ServiceContext) *ServerHandle {
	return &ServerHandle{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		PiDB:   relationDB.NewVidmgrInfoRepo(ctx),
	}
}

func (l *ServerHandle) ActionCheck() error {
	//l.Infof("ActionCheck req:%v", in)
	fmt.Println("[****] func (l *ServerHandle) ActionCheck() error ")
	//需要做的操作，查旬数据库
	now := time.Now().Unix()
	//过滤条件为：在线设备且超时时间为60秒
	filter := relationDB.VidmgrFilter{LastLoginTime: struct {
		Start int64
		End   int64
	}{Start: 0, End: now - clients.VIDMGRTIMEOUT}, VidmgrStatus: def.DeviceStatusOnline}
	di, err := l.PiDB.FindAllFilter(l.ctx, filter)
	if err != nil {
		return err
	}
	if len(di) > 0 {
		for _, v := range di {
			v.VidmgrStatus = def.DeviceStatusOffline
			l.PiDB.Update(l.ctx, v) //更新数据库
		}
	} else {
		//do nothing
	}
	//判断当前时间与最后login时间，是否超过30s
	//1分钟会执行一次
	return nil
}

func (l *ServerHandle) ActionInit() error {
	//l.Infof("ActionCheck req:%v", in)
	var vidInfo *relationDB.VidmgrInfo
	//fmt.Println("[**ActionInit**]0 ", utils.FuncName())
	//查找流服务的数据库：根据IP和端口确定一个流服务
	var (
		c      = l.svcCtx.Config
		filter = relationDB.VidmgrFilter{VidmgrIpV4: utils.InetAtoN(c.Mediakit.Host), VidmgrPort: c.Mediakit.Port}
	)
	//fmt.Println("[**ActionInit**]1 ", utils.FuncName())
	size, err := l.PiDB.CountByFilter(l.ctx, filter)
	if err != nil {
		fmt.Errorf("MediaServer init data countfilter error")
		return err
	}
	//找到存在一条流服务 更新这条服务
	if size > 0 {
		//update
		fmt.Println("[**ActionInit**]3 ", utils.FuncName())
		page := vid.PageInfo{}
		di, err := l.PiDB.FindByFilter(l.ctx, filter, common.ToPageInfoWithDefault(&page, &def.PageInfo{
			Page: 1, Size: 20,
			Orders: []def.OrderBy{{"created_time", def.OrderDesc}, {"vidmgr_id", def.OrderDesc}},
		}))
		if err != nil {
			fmt.Errorf("MediaServer init data find filter error")
			return err
		}
		if di[0].VidmgrSecret != c.Mediakit.Secret {
			di[0].VidmgrSecret = c.Mediakit.Secret
			err = l.PiDB.Update(l.ctx, di[0])
		}
		vidInfo = di[0]
	} else {
		//流服务还未存在的情况就插入这条服务
		fmt.Println("[**ActionInit**]4 ", utils.FuncName())
		dbDocker := &relationDB.VidmgrInfo{
			VidmgrID:     deviceAuth.GetStrProductID(l.svcCtx.VidmgrID.GetSnowflakeId()),
			VidmgrName:   "default Docker",
			VidmgrIpV4:   utils.InetAtoN(c.Mediakit.Host),
			VidmgrPort:   c.Mediakit.Port,
			VidmgrSecret: c.Mediakit.Secret,
			VidmgrStatus: 2, //默认设置离线状态
			VidmgrType:   1, //ZLmediakit
			MediasvrType: 1, //docker模式
			Desc:         "",
			Tags:         map[string]string{},
		}
		err = l.PiDB.Insert(l.ctx, dbDocker)
		if err != nil {
			l.Errorf("%s.Insert err=%+v", utils.FuncName(), err)
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
	//仅考虑docker的模式
	//STEP3  配置流服务
	if len(currentConf.Data) > 0 {
		currentConf.Data[0].GeneralMediaServerId = vidInfo.VidmgrID
		//docker通信IP用eth0 从zlmediakit->vidsvr
		common.SetDefaultConfig(l.svcCtx.Config.Restconf.Host, int64(l.svcCtx.Config.Restconf.Port), &currentConf.Data[0])
		fmt.Println("[****setting****] ", l.svcCtx.Config.Mediakit.Host, int64(l.svcCtx.Config.Restconf.Port))
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
		confRepo := relationDB.NewVidmgrConfigRepo(l.ctx)
		//查找config配置
		confRepo.FindOneByFilter(l.ctx, relationDB.VidmgrConfigFilter{
			//VidmgrIPv4: vidInfo.VidmgrIpV4,
			//VidmgrPort: vidInfo.VidmgrPort,
			Secret: vidInfo.VidmgrSecret,
			//VidmgrIDs: []string{vidInfo.VidmgrID},
		})
		if err != nil {
			l.Errorf("%s.Can find vidmgr config err=%v", utils.FuncName(), utils.Fmt(err))
			confRepo.Insert(l.ctx, common.ToVidmgrConfigDB1(&currentConf.Data[0]))
		} else {
			//update
			confRepo.Update(l.ctx, common.ToVidmgrConfigDB1(&currentConf.Data[0]))
		}

		//STEP4 更新状态
		fmt.Println("[*****test7*****]", utils.FuncName())
		if vidInfo.VidmgrStatus != def.DeviceStatusOnline {
			//UPDATE
			vidInfo.VidmgrStatus = def.DeviceStatusOnline
			vidInfo.FirstLogin = time.Now()
			vidInfo.LastLogin = time.Now()

			err := l.PiDB.Update(l.ctx, vidInfo)
			if err != nil {
				er := errors.Fmt(err)
				l.Errorf("%s.rpc.ManageVidmgr req=%v err=%v", utils.FuncName(), vidInfo, er)
				return er
			}
			fmt.Println("[*****test8*****] success", utils.FuncName())
			return nil
		}
	}
	fmt.Println("[*****test9*****]error:配置错误", utils.FuncName())
	return nil
}
