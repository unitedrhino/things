package devicemanagelogic

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/spf13/cast"
	"time"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceInfoCreateLogic {
	return &DeviceInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

/*
发现返回true 没有返回false
*/
func (l *DeviceInfoCreateLogic) CheckDevice(in *dm.DeviceInfo) (bool, error) {
	_, err := l.svcCtx.DeviceInfo.FindOneByProductIDDeviceName(l.ctx, in.ProductID, in.DeviceName)
	switch err {
	case mysql.ErrNotFound:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}

/*
发现返回true 没有返回false
*/
func (l *DeviceInfoCreateLogic) CheckProduct(in *dm.DeviceInfo) (bool, error) {
	_, err := l.svcCtx.ProductInfo.FindOne(l.ctx, in.ProductID)
	switch err {
	case mysql.ErrNotFound:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}

// 新增设备
func (l *DeviceInfoCreateLogic) DeviceInfoCreate(in *dm.DeviceInfo) (*dm.Response, error) {
	find, err := l.CheckDevice(in)
	if err != nil {
		l.Errorf("%s.CheckDevice in=%v\n", utils.FuncName(), in)
		return nil, errors.Database.AddDetail(err)
	} else if find == true {
		return nil, errors.Duplicate.AddDetail("DeviceName:" + in.DeviceName)
	}
	find, err = l.CheckProduct(in)
	if err != nil {
		l.Errorf("%s.CheckProduct in=%v", utils.FuncName(), in)
		return nil, errors.Database.AddDetail(err)
	} else if find == false {
		return nil, errors.Parameter.AddDetail("not find product id:" + cast.ToString(in.ProductID))
	}
	err = l.InitDevice(in)
	if err != nil {
		return nil, err
	}
	di := mysql.DeviceInfo{
		ProductID:   in.ProductID,  // 产品id
		DeviceName:  in.DeviceName, // 设备名称
		Secret:      utils.GetPwdBase64(20),
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	if in.Tags != nil {
		tags, err := json.Marshal(in.Tags)
		if err == nil {
			di.Tags = sql.NullString{
				String: string(tags),
				Valid:  true,
			}
		}
	} else {
		di.Tags = sql.NullString{
			String: "{}",
			Valid:  true,
		}
	}
	if in.LogLevel != def.Unknown {
		di.LogLevel = def.LogClose
	}
	_, err = l.svcCtx.DeviceInfo.Insert(l.ctx, &di)
	if err != nil {
		l.Errorf("AddDevice.DeviceInfo.Insert err=%+v", err)
		return nil, errors.System.AddDetail(err)
	}
	return &dm.Response{}, nil
}

func (l *DeviceInfoCreateLogic) InitDevice(in *dm.DeviceInfo) error {
	pt, err := l.svcCtx.SchemaRepo.GetSchemaModel(l.ctx, in.ProductID)
	if err != nil {
		return errors.System.AddDetail(err)
	}
	err = l.svcCtx.SchemaManaRepo.InitDevice(l.ctx, pt, in.ProductID, in.DeviceName)
	if err != nil {
		return errors.Database.AddDetail(err)
	}
	err = l.svcCtx.SDKLogRepo.InitDevice(l.ctx, in.ProductID, in.DeviceName)
	if err != nil {
		return errors.Database.AddDetail(err)
	}
	return nil
}
