package mysql

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ GatewayDeviceModel = (*customGatewayDeviceModel)(nil)

type (
	// GatewayDeviceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGatewayDeviceModel.
	GatewayDeviceModel interface {
		gatewayDeviceModel
		CreateList(ctx context.Context, gateway *device.Core, subDevices []*device.Core) error
		DeleteList(ctx context.Context, gateway *device.Core, subDevices []*device.Core) error
		FindByFilter(ctx context.Context, filter GatewayDeviceFilter, page def.PageInfo) ([]*DeviceInfo, error)
		CountByFilter(ctx context.Context, filter GatewayDeviceFilter) (size int64, err error)
	}

	customGatewayDeviceModel struct {
		*defaultGatewayDeviceModel
		deviceInfoTable string
	}
	GatewayDeviceFilter struct {
		Gateway device.Core //必填
	}
)

func (g GatewayDeviceFilter) FmtSql(sql sq.SelectBuilder) sq.SelectBuilder {
	sql = sql.Where("`gatewayProductID`=? and `gatewayDeviceName`=?", g.Gateway.ProductID, g.Gateway.DeviceName)
	return sql
}

func (c customGatewayDeviceModel) FindByFilter(ctx context.Context, f GatewayDeviceFilter, page def.PageInfo) ([]*DeviceInfo, error) {
	var resp []*DeviceInfo
	sql := sq.Select("di.*").From(c.table + "as gd").
		LeftJoin(fmt.Sprintf("%s as di on di.productID=gd.productID and di.deviceName=gd.deviceName", c.deviceInfoTable)).
		Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset()))
	sql = f.FmtSql(sql)
	query, arg, err := sql.ToSql()
	if err != nil {
		return nil, err
	}
	err = c.conn.QueryRowsCtx(ctx, &resp, query, arg...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (c customGatewayDeviceModel) CountByFilter(ctx context.Context, f GatewayDeviceFilter) (size int64, err error) {
	sql := sq.Select("count(1)").From(c.table)
	sql = f.FmtSql(sql)
	query, arg, err := sql.ToSql()
	if err != nil {
		return 0, err
	}
	err = c.conn.QueryRowCtx(ctx, &size, query, arg...)

	switch err {
	case nil:
		return size, nil
	default:
		return 0, err
	}
}

// NewGatewayDeviceModel returns a model for the database table.
func NewGatewayDeviceModel(conn sqlx.SqlConn) GatewayDeviceModel {
	return &customGatewayDeviceModel{
		defaultGatewayDeviceModel: newGatewayDeviceModel(conn),
		deviceInfoTable:           "`device_info`",
	}
}

func (c customGatewayDeviceModel) CreateList(ctx context.Context, gateway *device.Core, subDevices []*device.Core) error {
	return c.conn.Transact(func(session sqlx.Session) error {
		for _, v := range subDevices {
			sql := sq.Select("count(1)").From(c.deviceInfoTable)
			//f := GatewayDeviceFilter{ProductID: v.ProductID, DeviceName: v.DeviceName}
			//sql = f.FmtSql(sql)
			query, arg, err := sql.ToSql()
			if err != nil {
				logx.WithContext(ctx).Errorf("customGatewayDeviceModel.GatewayDeviceFilter.ToSql data:%v err:%v", v, err)
				return err
			}
			var size int64
			err = session.QueryRowCtx(ctx, &size, query, arg...)
			if err != nil {
				logx.WithContext(ctx).Errorf("customGatewayDeviceModel.deviceInfoTable.QueryRowCtx data:%v err:%v", v, err)
				return err
			}
			if size == 0 {
				return errors.Parameter.WithMsgf("设备不存在:产品ID:%v,设备名:%", v.ProductID, v.DeviceName)
			}
			query = fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", c.table, gatewayDeviceRowsExpectAutoSet)
			_, err = session.ExecCtx(ctx, query, gateway.ProductID, gateway.DeviceName, v.ProductID, v.DeviceName)
			if err != nil {
				logx.WithContext(ctx).Errorf("customGatewayDeviceModel.CreateList data:%v err:%v", v, err)
				return err
			}
		}
		return nil
	})
}

func (c customGatewayDeviceModel) DeleteList(ctx context.Context, gateway *device.Core, subDevices []*device.Core) error {
	return c.conn.Transact(func(session sqlx.Session) error {
		for _, v := range subDevices {
			query := fmt.Sprintf("delete from %s where `gatewayProductID` = ? and `GatewayDeviceName`=? and `productID`=? and `deviceName`=?", c.table)
			_, err := session.ExecCtx(ctx, query, gateway.ProductID, gateway.DeviceName, v.ProductID, v.DeviceName)
			if err != nil {
				logx.WithContext(ctx).Errorf("customGatewayDeviceModel.DeleteList gateway:%v data:%v err:%v", gateway, v, err)
				return err
			}
		}
		return nil
	})
}
