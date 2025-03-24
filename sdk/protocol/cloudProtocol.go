package protocol

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/topics"
	"github.com/mitchellh/mapstructure"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"sync"
	"time"
)

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

type SyncDevicesFunc[pConf ConfImp] func(ctx context.Context, conf pConf, product *dm.ProductInfo) error

type CloudProtocol[pConf ConfImp] struct {
	ConfMap      map[string]ConfInfo[pConf]
	ConfMapMutex sync.RWMutex
	SyncDevices  SyncDevicesFunc[pConf]
	InitFunc     func(conf pConf) (close func(), err error)
	*CoreProtocol
}

type ConfInfo[pConf ConfImp] struct {
	Conf  pConf
	Close func()
}

type CloudProtocolConf[pConf ConfImp] struct {
	ServerName string
	DmClient   zrpc.Client
	TimedM     zrpc.Client
	NodeID     int64
}

func NewCloudProtocol[pConf ConfImp](c conf.EventConf, pi *dm.ProtocolInfo, pc *CloudProtocolConf[pConf]) (*CloudProtocol[pConf], error) {
	lp, err := NewCoreProtocol(c, pi, &CoreProtocolConf{
		ServerName: pc.ServerName,
		DmClient:   pc.DmClient,
		TimedM:     pc.TimedM,
		NodeID:     pc.NodeID,
	})
	if err != nil {
		return nil, err
	}
	return &CloudProtocol[pConf]{
		CoreProtocol: lp,
		ConfMap:      map[string]ConfInfo[pConf]{},
	}, nil
}

func (p *CloudProtocol[pConf]) Start() error {
	ctx := context.Background()
	err := p.CoreProtocol.Start()
	if err != nil {
		return err
	}
	pi, err := p.ProtocolM.ProtocolInfoRead(ctx, &dm.WithIDCode{Code: p.Pi.Code})
	if err != nil {
		return err
	}
	var cs []pConf
	for _, cMap := range pi.ConfigInfos {
		var c pConf
		if len(cMap.Config) == 0 {
			continue
		}
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
	p.RunTimerHandles()
	return nil
}

/*
{
"accessKeyID":"LTAI5tHECH5pcapPoam3vdLK",
"accessKeySecret":"LTAI5tHECH5pcapPoam3vdLK",
"regionID":"cn-shanghai",
"iotInstanceID":"iot-060aate5",
"uid":"209182205384655582",
"consumerGroupID":"DEFAULT_GROUP"
}
*/

func (p *CloudProtocol[pConf]) ConfigChange(ctx context.Context, opt ConfigOpt, c pConf) error {
	key := c.GenKey()
	switch opt {
	case ConfigOptCreate, ConfigOptUpdate:
		err := c.Init()
		if err != nil {
			return err
		}
		var Close func()
		if p.InitFunc != nil {
			Close, err = p.InitFunc(c)
			if err != nil {
				return err
			}
		}
		p.ConfMap[key] = ConfInfo[pConf]{Conf: c, Close: Close}
	case ConfigOptDelete:
		if p.ConfMap[key].Close != nil {
			p.ConfMap[key].Close()
		}
		delete(p.ConfMap, key)
	}
	return nil
}

func (p *CloudProtocol[pConf]) GetConf(key string) *pConf {
	p.ConfMapMutex.RLock()
	defer p.ConfMapMutex.RUnlock()
	c, ok := p.ConfMap[key]
	if !ok {
		return nil
	}
	return &c.Conf
}
func (p *CloudProtocol[pConf]) GetAllConf() (ret []pConf) {
	p.ConfMapMutex.RLock()
	defer p.ConfMapMutex.RUnlock()
	for _, v := range p.ConfMap {
		ret = append(ret, v.Conf)
	}
	return ret
}

func (p *CloudProtocol[pConf]) RegisterDeviceSync(fieldName string /*自定义协议的对应协议code的字段名*/, f SyncDevicesFunc[pConf]) error {
	err := p.RegisterTimerHandler(func(ctx context.Context, t time.Time) error {
		pis, err := p.ProductM.ProductInfoIndex(ctx, &dm.ProductInfoIndexReq{
			ProtocolCode: p.Pi.Code,
		})
		if err != nil {
			return err
		}
		for _, pi := range pis.List {
			pc := p.GetProtocolConf(pi)
			key := pc[fieldName]
			if key == "" {
				continue
			}
			c := p.GetConf(key)
			if c == nil {
				continue
			}
			err := f(ctx, *c, pi)
			if err != nil {
				logx.WithContext(ctx).Error(err)
			}
		}
		return nil
	})
	return err
}

func (p *CloudProtocol[pConf]) UpdateConfig(ctx context.Context, c []pConf) error {
	var KeySet = map[string]struct{}{}
	//新增配置
	p.ConfMapMutex.Lock()
	defer p.ConfMapMutex.Unlock()
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

func (p *CloudProtocol[pConf]) RegisterInitHandler(
	handle func(conf pConf) (close func(), err error)) error {
	p.InitFunc = handle
	return nil
}

func (p *CloudProtocol[pConf]) RegisterConfigChange() error {
	err := p.FastEvent.Subscribe(fmt.Sprintf(topics.DmProtocolInfoUpdate, p.Pi.Code),
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
