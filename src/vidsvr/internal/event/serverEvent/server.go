package serverEvent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/domain/deviceAuth"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/logic"
	"github.com/i-Things/things/src/vidsvr/internal/logic/zlmedia/index"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"
	"github.com/i-Things/things/src/vidsvr/pb/vid"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
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
	fmt.Println("[****] func (l *ServerHandle) ActionInit database() error ")
	var (
		c      = l.svcCtx.Config
		filter = relationDB.VidmgrFilter{VidmgrIpV4: utils.InetAtoN(c.Mediakit.Host), VidmgrPort: c.Mediakit.Port}
	)
	size, err := l.PiDB.CountByFilter(l.ctx, filter)
	if err != nil {
		fmt.Errorf("MediaServer init data countfilter error")
		return err
	}
	if size > 0 {
		//update
		page := vid.PageInfo{}
		di, err := l.PiDB.FindByFilter(l.ctx, filter, logic.ToPageInfoWithDefault(&page, &def.PageInfo{
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
	} else {
		//create
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
	}
	//config dockerServer
	config := new(types.ServerConfig)
	index.SetDefaultConfig(c.Mediakit.Host, int64(c.Restconf.Port), config)
	byte4, err := json.Marshal(config)
	var tdata map[string]interface{}
	err = json.Unmarshal(byte4, &tdata)
	tdata["secret"] = c.Mediakit.Secret
	byte4, err = json.Marshal(tdata)
	if err != nil {
		er := errors.Fmt(err)
		fmt.Print("%s map string phares failed  err=%+v", utils.FuncName(), er)
		return er
	}
	preUrl := fmt.Sprintf("http://%s:%s/index/api/setServerConfig", c.Mediakit.Host, c.Mediakit.Port)
	request, error := http.NewRequest("POST", preUrl, bytes.NewBuffer(byte4))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		fmt.Println(err)
		return err
	}
	defer response.Body.Close()
	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	//body, _ := ioutil.ReadAll(response.Body)
	body, err := io.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))

	return nil
}
