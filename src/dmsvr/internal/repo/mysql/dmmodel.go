package mysql

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/store"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	DmModel interface {
		FindProductsByFilter(ctx context.Context, filter ProductFilter, page def.PageInfo) ([]*ProductInfo, error)
		FindDevicesByFilter(ctx context.Context, filter DeviceFilter, page def.PageInfo) ([]*DeviceInfo, error)
		GetDevicesCountByFilter(ctx context.Context, filter DeviceFilter) (size int64, err error)
		GetProductsCountByFilter(ctx context.Context, filter ProductFilter) (size int64, err error)
		Insert(ctx context.Context, pi *ProductInfo, pt *ProductSchema) error
		Delete(ctx context.Context, productID string) error
	}
	ProductFilter struct {
		DeviceType  int64
		ProductName string
		ProductIDs  []string
	}
	DeviceFilter struct {
		ProductID  string
		DeviceName string
		Tags       map[string]string
	}
	defaultDmModel struct {
		sqlx.SqlConn
		productInfo   string
		deviceInfo    string
		productSchema string
		deviceLog     string
		ProductInfoModel
	}
)

func NewDmModel(conn sqlx.SqlConn) DmModel {
	return &defaultDmModel{
		SqlConn:          conn,
		productInfo:      "`product_info`",
		deviceInfo:       "`device_info`",
		productSchema:    "`product_schema`",
		deviceLog:        "`device_log`",
		ProductInfoModel: NewProductInfoModel(conn),
	}
}

func (d *DeviceFilter) FmtSql(sql sq.SelectBuilder) sq.SelectBuilder {
	if d.ProductID != "" {
		sql = sql.Where("`ProductID` = ?", d.ProductID)
	}
	if d.DeviceName != "" {
		sql = sql.Where("`DeviceName` like ?", "%"+d.DeviceName+"%")
	}
	if d.Tags != nil {
		for k, v := range d.Tags {
			sql = sql.Where("JSON_CONTAINS(`tags`, JSON_OBJECT(?,?))",
				k, v)
		}
	}
	return sql
}

func (m *defaultDmModel) FindDevicesByFilter(ctx context.Context, f DeviceFilter, page def.PageInfo) ([]*DeviceInfo, error) {
	var resp []*DeviceInfo
	sql := sq.Select(deviceInfoRows).From(m.deviceInfo).Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset()))
	sql = f.FmtSql(sql)
	query, arg, err := sql.ToSql()
	if err != nil {
		return nil, err
	}
	err = m.QueryRowsCtx(ctx, &resp, query, arg...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultDmModel) GetDevicesCountByFilter(ctx context.Context, f DeviceFilter) (size int64, err error) {
	sql := sq.Select("count(1)").From(m.deviceInfo)
	sql = f.FmtSql(sql)
	query, arg, err := sql.ToSql()
	if err != nil {
		return 0, err
	}
	err = m.QueryRowCtx(ctx, &size, query, arg...)

	switch err {
	case nil:
		return size, nil
	default:
		return 0, err
	}
}

func (p *ProductFilter) FmtSql(sql sq.SelectBuilder) sq.SelectBuilder {
	if p.DeviceType != 0 {
		sql = sql.Where("DeviceType=?", p.DeviceType)
	}
	if p.ProductName != "" {
		sql = sql.Where("ProductName like ?", "%"+p.ProductName+"%")
	}
	if len(p.ProductIDs) != 0 {
		sql = sql.Where(fmt.Sprintf("ProductID in (%v)", store.ArrayToSql(p.ProductIDs)))
	}
	return sql
}

func (m *defaultDmModel) FindProductsByFilter(ctx context.Context, f ProductFilter, page def.PageInfo) ([]*ProductInfo, error) {
	var resp []*ProductInfo
	sql := sq.Select(productInfoRows).From(m.productInfo).Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset()))
	sql = f.FmtSql(sql)
	query, arg, err := sql.ToSql()
	if err != nil {
		return nil, err
	}
	err = m.QueryRowsCtx(ctx, &resp, query, arg...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultDmModel) GetProductsCountByFilter(ctx context.Context, f ProductFilter) (size int64, err error) {
	sql := sq.Select("count(1)").From(m.productInfo)
	sql = f.FmtSql(sql)
	query, arg, err := sql.ToSql()
	if err != nil {
		return 0, err
	}
	err = m.QueryRowCtx(ctx, &size, query, arg...)

	switch err {
	case nil:
		return size, nil
	default:
		return 0, err
	}
}

func (m *defaultDmModel) Delete(ctx context.Context, productID string) error {
	return m.Transact(func(session sqlx.Session) error {
		query := fmt.Sprintf("delete from %s where `ProductID` = ?", m.productInfo)
		_, err := session.Exec(query, productID)
		if err != nil {
			return err
		}
		query = fmt.Sprintf("delete from %s where `ProductID` = ?", m.deviceInfo)
		_, err = session.Exec(query, productID)
		if err != nil {
			return err
		}
		query = fmt.Sprintf("delete from %s where `ProductID` = ?", m.productSchema)
		_, err = session.Exec(query, productID)
		return err
	})
}

func (m *defaultDmModel) Insert(ctx context.Context, pi *ProductInfo, pt *ProductSchema) error {
	return m.Transact(func(session sqlx.Session) error {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.productInfo, productInfoRowsExpectAutoSet)
		_, err := session.ExecCtx(ctx, query, pi.ProductID, pi.ProductName, pi.ProductType, pi.AuthMode, pi.DeviceType, pi.CategoryID, pi.NetType, pi.DataProto, pi.AutoRegister, pi.Secret, pi.Desc, pi.CreatedTime, pi.UpdatedTime, pi.DeletedTime, pi.DevStatus)
		if err != nil {
			return err
		}
		query = fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.productSchema, productSchemaRowsExpectAutoSet)
		_, err = session.ExecCtx(ctx, query, pt.ProductID, pt.Schema, pt.CreatedTime, pt.UpdatedTime, pt.DeletedTime)
		return err
	})
}
