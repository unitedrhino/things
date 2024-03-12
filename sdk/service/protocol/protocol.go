package protocol

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/client/timedmanage"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/eventBus"
	"gitee.com/i-Things/share/events/topics"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/client/deviceinteract"
	"github.com/i-Things/things/service/dmsvr/client/devicemanage"
	"github.com/i-Things/things/service/dmsvr/client/productmanage"
	"github.com/i-Things/things/service/dmsvr/client/protocolmanage"
	"github.com/i-Things/things/service/dmsvr/client/schemamanage"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/mitchellh/mapstructure"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"time"
)

type SvrClient struct {
	ProtocolM      protocolmanage.ProtocolManage
	ProductM       productmanage.ProductManage
	SchemaM        schemamanage.SchemaManage
	DeviceM        devicemanage.DeviceManage
	DeviceInteract deviceinteract.DeviceInteract
	TimedM         timedmanage.TimedManage
}

type ConfImp interface {
	GenKey() string
	Equal(imp ConfImp) bool
	Init() error //配置初始化
}

type ConfigOpt int64

const (
	ConfigOptCreate ConfigOpt = iota
	ConfigOptUpdate
	ConfigOptDelete
)

type GetProductsFunc[pConf ConfImp] func(ctx context.Context, conf pConf) ([]*ProductInfo[pConf], error)
type GetDevicesFunc[pConf ConfImp] func(ctx context.Context, conf pConf, productID string) ([]*DeviceInfo, error)

type Protocol[pConf ConfImp] struct {
	FastEvent  *eventBus.FastEvent
	Pi         *dm.ProtocolInfo
	ServerName string
	SvrClient

	ConfMap map[string]ConfInfo[pConf]

	GetProducts GetProductsFunc[pConf]
	GetDevices  GetDevicesFunc[pConf]

	InitFunc func(conf pConf) (close func(), err error)

	productMap map[string]*ProductInfo[pConf]
}

type ConfInfo[pConf ConfImp] struct {
	Conf       pConf
	Close      func()
	ProductMap map[string]*ProductInfo[pConf]
}

type DeviceInfo struct {
	DeviceName string
}

type ProductInfo[pConf ConfImp] struct {
	ProductID   string
	ProductName string
	DeviceMap   map[string]*DeviceInfo
	Conf        pConf
}

type ProtocolConf[pConf ConfImp] struct {
	ServerName string
	DmClient   zrpc.Client
	TimedM     zrpc.Client
}

func NewProtocol[pConf ConfImp](c conf.EventConf, pi *dm.ProtocolInfo, pc *ProtocolConf[pConf]) (*Protocol[pConf], error) {
	e, err := eventBus.NewFastEvent(c, pc.ServerName)
	if err != nil {
		return nil, err
	}
	return &Protocol[pConf]{
		FastEvent:  e,
		Pi:         pi,
		ServerName: pc.ServerName,
		SvrClient: SvrClient{
			ProtocolM:      protocolmanage.NewProtocolManage(pc.DmClient),
			ProductM:       productmanage.NewProductManage(pc.DmClient),
			SchemaM:        schemamanage.NewSchemaManage(pc.DmClient),
			DeviceM:        devicemanage.NewDeviceManage(pc.DmClient),
			DeviceInteract: deviceinteract.NewDeviceInteract(pc.DmClient),
			TimedM:         timedmanage.NewTimedManage(pc.TimedM),
		},
		productMap: make(map[string]*ProductInfo[pConf]),
		ConfMap:    map[string]ConfInfo[pConf]{},
	}, nil
}

func (p *Protocol[pConf]) DefaultGetProducts() GetProductsFunc[pConf] {
	return func(ctx context.Context, conf pConf) (ret []*ProductInfo[pConf], err error) {
		list, err := p.ProductM.ProductInfoIndex(ctx, &dm.ProductInfoIndexReq{ProtocolCode: p.Pi.Code})
		if err != nil {
			return nil, err
		}
		for _, v := range list.List {
			ret = append(ret, &ProductInfo[pConf]{ProductID: v.ProductID, ProductName: v.ProductName})
		}
		return
	}
}

func (p *Protocol[pConf]) DefaultGetDevices() GetDevicesFunc[pConf] {
	return func(ctx context.Context, conf pConf, productID string) (ret []*DeviceInfo, err error) {
		list, err := p.DeviceM.DeviceInfoIndex(ctx, &dm.DeviceInfoIndexReq{ProductID: productID})
		if err != nil {
			return nil, err
		}
		for _, v := range list.List {
			ret = append(ret, &DeviceInfo{DeviceName: v.DeviceName})
		}
		return
	}
}

func (p *Protocol[pConf]) Start(GetProducts GetProductsFunc[pConf], GetDevices GetDevicesFunc[pConf]) error {
	p.GetProducts = p.DefaultGetProducts()
	p.GetDevices = p.DefaultGetDevices()
	if GetProducts != nil {
		p.GetProducts = GetProducts
	}
	if GetDevices != nil {
		p.GetDevices = GetDevices
	}

	ctx := context.Background()
	_, err := p.ProtocolM.ProtocolInfoCreate(ctx, p.Pi) //初始化协议
	if err != nil && !errors.Cmp(errors.Fmt(err), errors.Duplicate) {
		logx.Must(err)
	}

	pi, err := p.ProtocolM.ProtocolInfoRead(ctx, &dm.WithIDCode{Code: p.Pi.Code})
	if err != nil {
		return err
	}
	var cs []pConf
	for _, cMap := range pi.ConfigInfos {
		var c pConf
		err := mapstructure.Decode(cMap.Config, &c)
		if err != nil {
			return err
		}
		cs = append(cs, c)
	}
	err = p.UpdateConfig(ctx, cs)
	if err != nil {
		return err
	}

	err = p.FastEvent.Start()
	if err != nil {
		return err
	}
	return nil
}

func (p *Protocol[pConf]) ConfigChange(ctx context.Context, opt ConfigOpt, c pConf) error {
	key := c.GenKey()
	switch opt {
	case ConfigOptCreate, ConfigOptUpdate:
		err := c.Init()
		if err != nil {
			return err
		}
		productList, err := p.GetProducts(ctx, c)
		if err != nil {
			return err
		}
		var ProductMap = map[string]*ProductInfo[pConf]{}
		for _, v := range productList {
			devices, err := p.GetDevices(ctx, c, v.ProductID)
			if err != nil {
				return err
			}
			var deviceMap = map[string]*DeviceInfo{}
			for _, v := range devices {
				deviceMap[v.DeviceName] = v
			}
			v.DeviceMap = deviceMap
			v.Conf = c
			ProductMap[v.ProductID] = v
			p.productMap[v.ProductID] = v
		}
		var Close func()
		if p.InitFunc != nil {
			Close, err = p.InitFunc(c)
			if err != nil {
				return err
			}
		}
		p.ConfMap[key] = ConfInfo[pConf]{ProductMap: ProductMap, Conf: c, Close: Close}
		for key := range ProductMap {
			err := p.ProductInit(ctx, key)
			if err != nil {
				return err
			}
		}

	case ConfigOptDelete:
		if p.ConfMap[key].Close != nil {
			p.ConfMap[key].Close()
		}
		for k := range p.ConfMap[key].ProductMap {
			delete(p.productMap, k)
		}
		delete(p.ConfMap, key)
	}
	return nil
}
func (p *Protocol[pConf]) ProductInit(ctx context.Context, productID string) error {
	pi := p.productMap[productID]
	_, err := p.ProductM.ProductInfoRead(ctx, &dm.ProductInfoReadReq{
		ProductID: productID,
	})
	if err != nil {
		if !errors.Cmp(errors.Fmt(err), errors.NotFind) {
			return err
		}
		_, err := p.ProductM.ProductInfoCreate(ctx, &dm.ProductInfo{
			ProductID:    productID,
			ProductName:  pi.ProductName,
			ProtocolCode: p.Pi.Code,
			Desc:         utils.ToRpcNullString(fmt.Sprintf("%s自动生成", p.Pi.Name)),
		})
		if err != nil {
			return nil
		}
	}
	list, err := p.DeviceM.DeviceInfoIndex(ctx, &dm.DeviceInfoIndexReq{
		ProductID: productID,
	})
	if err != nil {
		return nil
	}
	var deviceMap = map[string]*DeviceInfo{}
	for k, v := range pi.DeviceMap {
		deviceMap[k] = v
	}
	for _, v := range list.List {
		delete(deviceMap, v.DeviceName)
	}
	for k := range deviceMap {
		_, err := p.DeviceM.DeviceInfoCreate(ctx, &dm.DeviceInfo{
			ProductID:  productID,
			DeviceName: k,
		})
		if err != nil {
			logx.WithContext(ctx).Error(err)
			return err
		}
	}
	return nil
}

func (p *Protocol[pConf]) GetDeviceConf(ctx context.Context, productID string, deviceName string) (pConf, error) {
	return p.productMap[productID].Conf, nil
}

func (p *Protocol[pConf]) UpdateConfig(ctx context.Context, c []pConf) error {
	var KeySet = map[string]struct{}{}
	//新增配置
	for _, v := range c {
		key := v.GenKey()
		KeySet[key] = struct{}{}
		conf, ok := p.ConfMap[key]
		if ok {
			if !conf.Conf.Equal(v) { //配置项做了调整
				err := p.ConfigChange(ctx, ConfigOptUpdate, v)
				if err != nil {
					return err
				}
			}
			continue
		}
		err := p.ConfigChange(ctx, ConfigOptCreate, v)
		if err != nil {
			return err
		}
	}
	//删除配置
	for key, v := range p.ConfMap {
		_, ok := KeySet[key]
		if ok {
			continue
		}
		p.ConfigChange(ctx, ConfigOptDelete, v.Conf)
	}
	return nil
}

func (p *Protocol[pConf]) RegisterDeviceMsgDownHandler(
	handle func(ctx context.Context, info *devices.InnerPublish) error) error {
	err := p.FastEvent.QueueSubscribe(fmt.Sprintf(topics.DeviceDownAll, p.Pi.Code),
		func(ctx context.Context, t time.Time, body []byte) error {
			info := devices.GetPublish(body)
			logx.WithContext(ctx).Infof("mqtt.SubDevMsg protocolCode:%v Handle:%s Type:%v Payload:%v",
				info.ProtocolCode, info.Handle, info.Type, string(info.Payload))
			err := handle(ctx, info)
			return err
		})
	return err
}

func (p *Protocol[pConf]) RegisterInitHandler(
	handle func(conf pConf) (close func(), err error)) error {
	p.InitFunc = handle
	return nil
}

func (p *Protocol[pConf]) RegisterGetProductWithDevice() error {
	err := p.FastEvent.Subscribe(fmt.Sprintf(eventBus.DmProtocolInfoUpdate, p.Pi.Code),
		func(ctx context.Context, t time.Time, body []byte) error {
			var conf []pConf
			err := json.Unmarshal(body, &conf)
			if err != nil {
				return err
			}
			err = p.UpdateConfig(ctx, conf)
			return err
		})
	return err
}

func (p *Protocol[pConf]) RegisterConfigChange() error {
	err := p.FastEvent.Subscribe(fmt.Sprintf(eventBus.DmProtocolInfoUpdate, p.Pi.Code),
		func(ctx context.Context, t time.Time, body []byte) error {
			var conf []pConf
			err := json.Unmarshal(body, &conf)
			if err != nil {
				return err
			}
			err = p.UpdateConfig(ctx, conf)
			return err
		})
	return err
}

func (p *Protocol[pConf]) genCode() string {
	return fmt.Sprintf("protocol-%s-timer", p.Pi.Code)
}

func (p *Protocol[pConf]) genTimerTopic() string {
	return fmt.Sprintf("server.things.%s.protocol.timer", p.ServerName)
}

// 定时同步设备信息,产品信息 如果不需要可以不注册
func (p *Protocol[pConf]) RegisterTimerHandler(f func(ctx context.Context, t time.Time) error) error {
	ctx := context.Background()
	_, err := p.TimedM.TaskInfoCreate(ctx, &timedmanage.TaskInfo{
		GroupCode: def.TimedIThingsQueueGroupCode,                                //组编码
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
	err = p.FastEvent.Subscribe(p.genTimerTopic(), func(ctx context.Context, t time.Time, body []byte) error {
		err = f(ctx, t)
		return err
	})
	return err
}
func (p *Protocol[pConf]) DevPubMsg(ctx context.Context, publishMsg *devices.DevPublish) error {
	publishMsg.ProtocolCode = p.Pi.Code
	err := p.FastEvent.Publish(ctx, fmt.Sprintf(topics.DeviceUpMsg, publishMsg.Handle, publishMsg.ProductID, publishMsg.DeviceName), publishMsg)
	if err != nil {
		logx.Errorf("%s.publish  err:%v", utils.FuncName(), err)
		return err
	}
	return nil
}
