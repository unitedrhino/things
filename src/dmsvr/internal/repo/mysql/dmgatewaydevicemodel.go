package mysql

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ DmGatewayDeviceModel = (*customDmGatewayDeviceModel)(nil)

type (
	// DmGatewayDeviceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDmGatewayDeviceModel.
	DmGatewayDeviceModel interface {
		dmGatewayDeviceModel
		CreateList(ctx context.Context, gateway *devices.Core, subDevices []*devices.Core) error
		DeleteList(ctx context.Context, gateway *devices.Core, subDevices []*devices.Core) error
		FindByFilter(ctx context.Context, filter GatewayDeviceFilter, page *def.PageInfo) ([]*DmDeviceInfo, error)
		CountByFilter(ctx context.Context, filter GatewayDeviceFilter) (size int64, err error)
	}

	customDmGatewayDeviceModel struct {
		*defaultDmGatewayDeviceModel
		deviceInfoTable string
	}
	GatewayDeviceFilter struct {
		//网关和子设备 至少要有一个填写
		Gateway *devices.Core
		//网关和子设备 至少要有一个填写
		SubDevice *devices.Core
	}
)

func (c customDmGatewayDeviceModel) FmtSql(ctx context.Context, f GatewayDeviceFilter, sql sq.SelectBuilder) sq.SelectBuilder {
	if f.Gateway != nil { //通过网关获取旗下子设备列表
		sql = sql.LeftJoin(fmt.Sprintf("%s as di on di.productID=gd.productID and di.deviceName=gd.deviceName", c.deviceInfoTable)).
			Where("`gatewayProductID`=? and `gatewayDeviceName`=? and di.id IS NOT NULL", f.Gateway.ProductID, f.Gateway.DeviceName)
	} else {
		sql = sql.LeftJoin(fmt.Sprintf("%s as di on di.productID=gd.gatewayProductID and di.deviceName=gd.gatewayDeviceName", c.deviceInfoTable)).
			Where("gd.`productID`=? and gd.`deviceName`=? and di.id IS NOT NULL", f.SubDevice.ProductID, f.SubDevice.DeviceName)
	}

	//数据权限条件（企业版功能）
	if uc := ctxs.GetUserCtxOrNil(ctx); uc != nil && !uc.IsAllData { //存在用户态&&无所有数据权限
		mdProjectID := ctxs.GetMetaProjectID(ctx)
		if mdProjectID != 0 {
			sql = sql.Where("di.`ProjectID` = ?", mdProjectID)
		}
	}

	return sql
}

// NewDmGatewayDeviceModel returns a model for the database table.
func NewDmGatewayDeviceModel(conn sqlx.SqlConn) DmGatewayDeviceModel {
	return &customDmGatewayDeviceModel{
		defaultDmGatewayDeviceModel: newDmGatewayDeviceModel(conn),
		deviceInfoTable:             "`dm_device_info`",
	}
}

func (c customDmGatewayDeviceModel) FindByFilter(ctx context.Context, f GatewayDeviceFilter, page *def.PageInfo) ([]*DmDeviceInfo, error) {
	var resp []*DmDeviceInfo
	sql := sq.Select("di.*").From(c.table + "as gd").
		Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset()))
	sql = c.FmtSql(ctx, f, sql)
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

func (c customDmGatewayDeviceModel) CountByFilter(ctx context.Context, f GatewayDeviceFilter) (size int64, err error) {
	sql := sq.Select("count(1)").From(c.table + "as gd")
	sql = c.FmtSql(ctx, f, sql)
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

func (c customDmGatewayDeviceModel) CreateList(ctx context.Context, gateway *devices.Core, subDevices []*devices.Core) error {
	return c.conn.Transact(func(session sqlx.Session) error {
		for _, v := range subDevices {
			sql := sq.Select("count(1)").
				Where("`productID` = ? and `deviceName` = ?", v.ProductID, v.DeviceName).
				From(c.deviceInfoTable)
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
			query = fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?) ON duplicate KEY UPDATE id = id", c.table, dmGatewayDeviceRowsExpectAutoSet)
			_, err = session.ExecCtx(ctx, query, gateway.ProductID, gateway.DeviceName, v.ProductID, v.DeviceName)
			if err != nil {
				logx.WithContext(ctx).Errorf("customGatewayDeviceModel.CreateList data:%v err:%v", v, err)
				return err
			}
		}
		return nil
	})
}

func (c customDmGatewayDeviceModel) DeleteList(ctx context.Context, gateway *devices.Core, subDevices []*devices.Core) error {
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
