package protocol

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/client/timedmanage"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/eventBus"
	"gitee.com/unitedrhino/share/events/topics"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/client/deviceinteract"
	"gitee.com/unitedrhino/things/service/dmsvr/client/devicemanage"
	"gitee.com/unitedrhino/things/service/dmsvr/client/productmanage"
	"gitee.com/unitedrhino/things/service/dmsvr/client/protocolmanage"
	"gitee.com/unitedrhino/things/service/dmsvr/dmExport"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/netx"
	"github.com/zeromicro/go-zero/zrpc"
	"sync"
	"time"
)

type CoreSvrClient struct {
	ProtocolM      protocolmanage.ProtocolManage
	ProductM       productmanage.ProductManage
	ProductCache   dmExport.ProductCacheT
	DeviceCache    dmExport.DeviceCacheT
	SchemaCache    dmExport.DeviceSchemaCacheT
	DeviceM        devicemanage.DeviceManage
	DeviceInteract deviceinteract.DeviceInteract
	TimedM         timedmanage.TimedManage
	TimerHandles   []func(ctx context.Context, t time.Time) error
}

type CoreProtocol struct {
	FastEvent             *eventBus.FastEvent
	Pi                    *dm.ProtocolInfo
	ServerName            string
	ProductIDMap          map[string]string //key 是外部的产品ID,value是内部的产品ID
	UnitedRhinoProductIDs []string          //iThings 的产品ID列表
	ProductIDMapMutex     sync.RWMutex
	CoreSvrClient
	ThirdProductIDFieldName devices.ProtocolKey
	taskCreateOnce          sync.Once
}

type CoreProtocolConf struct {
	ServerName string
	DmClient   zrpc.Client
	TimedM     zrpc.Client
	NodeID     int64
	Port       int64
}

func NewCoreProtocol(c conf.EventConf, pi *dm.ProtocolInfo, pc *CoreProtocolConf) (*CoreProtocol, error) {
	e, err := eventBus.NewFastEvent(c, pc.ServerName, pc.NodeID)
	if err != nil {
		return nil, err
	}
	pm := productmanage.NewProductManage(pc.DmClient)
	di := devicemanage.NewDeviceManage(pc.DmClient)
	ps, err := dmExport.NewProductSchemaCache(pm, e)
	if err != nil {
		return nil, err
	}
	sc, err := dmExport.NewDeviceSchemaCache(di, ps, e)
	if err != nil {
		return nil, err
	}
	dc, err := dmExport.NewDeviceInfoCache(di, e)
	if err != nil {
		return nil, err
	}
	pic, err := dmExport.NewProductInfoCache(pm, e)
	if err != nil {
		return nil, err
	}

	var timedM timedmanage.TimedManage
	if pc.TimedM != nil {
		timedM = timedmanage.NewTimedManage(pc.TimedM)
	}
	return &CoreProtocol{
		FastEvent:  e,
		Pi:         pi,
		ServerName: pc.ServerName,
		CoreSvrClient: CoreSvrClient{
			ProtocolM:      protocolmanage.NewProtocolManage(pc.DmClient),
			ProductM:       pm,
			SchemaCache:    sc,
			DeviceCache:    dc,
			ProductCache:   pic,
			DeviceM:        di,
			TimedM:         timedM,
			DeviceInteract: deviceinteract.NewDeviceInteract(pc.DmClient),
		},
	}, nil
}

func (p *CoreProtocol) Start() error {
	ctx := ctxs.WithRoot(context.Background())
	_, err := p.ProtocolM.ProtocolInfoCreate(ctx, p.Pi) //初始化协议
	if err != nil && !errors.Cmp(errors.Fmt(err), errors.Duplicate) {
		logx.Must(err)
	}
	utils.Go(ctx, func() {
		run := func() {
			_, err := p.ProtocolM.ProtocolServiceUpdate(ctx, &dm.ProtocolService{
				Code:   p.Pi.Code,
				Ip:     netx.InternalIp(),
				Port:   0,
				Status: def.True,
			})
			if err != nil {
				logx.WithContext(ctx).Error(err)
			}
		}
		run()
		tick := time.Tick(time.Minute)
		for _ = range tick {
			run()
		}
	})

	err = p.FastEvent.Start()
	if err != nil {
		return err
	}
	return nil
}

func (p *CoreProtocol) RunTimerHandles() {
	for _, f := range p.TimerHandles {
		err := f(ctxs.WithRoot(context.Background()), time.Now())
		if err != nil {
			logx.Error(err)
		}
	}
}

func (p *CoreProtocol) RegisterDeviceMsgDownHandler(
	handle func(ctx context.Context, info *devices.InnerPublish) error) error {
	err := p.FastEvent.QueueSubscribe(fmt.Sprintf(topics.DeviceDownAll, p.Pi.Code),
		func(ctx context.Context, t time.Time, body []byte) error {
			info := devices.GetPublish(body)
			logx.WithContext(ctx).Infof("mqtt.SubDevMsg protocolCode:%v Handle:%s Type:%v Payload:%v",
				info.ProtocolCode, info.Handle, info.Type, string(info.Payload))
			err := handle(ctxs.WithRoot(ctx), info)
			if err != nil {
				logx.WithContext(ctx).Errorf("mqtt.SubDevMsg protocolCode:%v Handle:%s Type:%v Payload:%v err:%v",
					info.ProtocolCode, info.Handle, info.Type, string(info.Payload), err)
			}

			return err
		})
	return err
}

func (p *CoreProtocol) DevPubMsg(ctx context.Context, publishMsg *devices.DevPublish) error {
	publishMsg.ProtocolCode = p.Pi.Code
	UpdateDeviceActivity(ctx, devices.Core{
		ProductID:  publishMsg.ProductID,
		DeviceName: publishMsg.DeviceName,
	})
	err := p.FastEvent.Publish(ctx, fmt.Sprintf(topics.DeviceUpMsg, publishMsg.Handle, publishMsg.ProductID, publishMsg.DeviceName), utils.MarshalNoErr(publishMsg))
	if err != nil {
		logx.WithContext(ctx).Errorf("devPublishToCloud  err:%v", err)
		return err
	} else {
		logx.WithContext(ctx).Infof("devPublishToCloud msg:%v", publishMsg)
	}

	return nil
}

func (p *CoreProtocol) genCode() string {
	return fmt.Sprintf("protocol-%s-timer", p.Pi.Code)
}

func (p *CoreProtocol) genTimerTopic() string {
	return fmt.Sprintf("server.things.%s.protocol.timer", p.ServerName)
}

// 定时同步设备信息,产品信息 如果不需要可以不注册
func (p *CoreProtocol) RegisterTimerHandler(f func(ctx context.Context, t time.Time) error) error {
	if p.TimedM == nil {
		return errors.Panic.AddMsg("timed not init")
	}
	ctx := context.Background()
	p.taskCreateOnce.Do(func() {
		_, err := p.TimedM.TaskInfoCreate(ctx, &timedmanage.TaskInfo{
			GroupCode: def.TimedUnitedRhinoQueueGroupCode,                            //组编码
			Type:      1,                                                             //任务类型 1 定时任务 2 延时任务
			Name:      fmt.Sprintf("自定义协议-%s-定时任务-数据同步", p.Pi.Name),                  // 任务名称
			Code:      p.genCode(),                                                   //任务编码
			Params:    fmt.Sprintf(`{"topic":"%s","payload":""}`, p.genTimerTopic()), // 任务参数,延时任务如果没有传任务参数会拿数据库的参数来执行
			CronExpr:  "@every 10m",                                                  // cron执行表达式
			Status:    def.StatusWaitRun,                                             // 状态
			Priority:  3,                                                             //优先级: 10:critical 最高优先级  3: default 普通优先级 1:low 低优先级
		})
		if err != nil && !errors.Cmp(errors.Fmt(err), errors.Duplicate) {
			logx.Must(err)
		}
	})
	p.TimerHandles = append(p.TimerHandles, f)
	err := p.FastEvent.Subscribe(p.genTimerTopic(), func(ctx context.Context, t time.Time, body []byte) error {
		if t.Before(time.Now().Add(-time.Second * 2)) { //2秒之前的跳过
			return nil
		}
		err := f(ctxs.WithRoot(context.Background()), t)
		return err
	})
	return err
}
func (p *CoreProtocol) RegisterProductIDSync() error {
	p.ThirdProductIDFieldName = devices.ProtocolKeyProductID
	err := p.RegisterTimerHandler(func(ctx context.Context, t time.Time) error {
		pis, err := p.ProductM.ProductInfoIndex(ctx, &dm.ProductInfoIndexReq{
			ProtocolCode: p.Pi.Code,
		})
		if err != nil {
			return err
		}
		p.ProductIDMapMutex.Lock()
		defer p.ProductIDMapMutex.Unlock()
		p.ProductIDMap = map[string]string{}
		p.UnitedRhinoProductIDs = nil
		for _, pi := range pis.List {
			id := pi.ProtocolConf[p.ThirdProductIDFieldName]
			if id == "" {
				continue
			}
			p.UnitedRhinoProductIDs = append(p.UnitedRhinoProductIDs, pi.ProductID)
			p.ProductIDMap[id] = pi.ProductID
		}
		return nil
	})
	return err
}

// 通过外部的产品iD查询联犀的产品iD
func (p *CoreProtocol) GetProductID(productID string) string {
	p.ProductIDMapMutex.RLock()
	defer p.ProductIDMapMutex.RUnlock()
	return p.ProductIDMap[productID]
}

func (p *CoreProtocol) GetUnitedRhinoProductIDs() []string {
	p.ProductIDMapMutex.RLock()
	defer p.ProductIDMapMutex.RUnlock()
	return p.UnitedRhinoProductIDs
}
