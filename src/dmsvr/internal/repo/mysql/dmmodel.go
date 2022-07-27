package mysql

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	DmModel interface {
		FindByProductInfo(ctx context.Context, deviceType int64, productName string, page def.PageInfo) ([]*ProductInfo, error)
		FindByProductID(ctx context.Context, productID string, page def.PageInfo) ([]*DeviceInfo, error)
		GetCountByProductID(ctx context.Context, productID string) (size int64, err error)
		GetCountByProductInfo(ctx context.Context, deviceType int64, productName string) (size int64, err error)
		Insert(ctx context.Context, pi *ProductInfo, pt *ProductSchema) error
		Delete(ctx context.Context, productID string) error
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

func (m *defaultDmModel) FindByProductID(ctx context.Context, productID string, page def.PageInfo) ([]*DeviceInfo, error) {
	var resp []*DeviceInfo
	sql := sq.Select(deviceInfoRows).From(m.deviceInfo).Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset()))
	if productID != "" {
		sql = sql.Where("`productID` = ?", productID)
	}
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

func (m *defaultDmModel) GetCountByProductID(ctx context.Context, productID string) (size int64, err error) {
	sql := sq.Select("count(1)").From(m.deviceInfo)
	if productID != "" {
		sql = sql.Where("`productID`=?", productID)
	}
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

func (m *defaultDmModel) FindByProductInfo(ctx context.Context, deviceType int64, productName string, page def.PageInfo) ([]*ProductInfo, error) {
	var resp []*ProductInfo
	sql := sq.Select(productInfoRows).From(m.productInfo).Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset()))
	if deviceType != 0 {
		sql = sql.Where("deviceType=?", deviceType)
	}
	if productName != "" {
		sql = sql.Where("productName like ?", "%"+productName+"%")
	}
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

func (m *defaultDmModel) GetCountByProductInfo(ctx context.Context, deviceType int64, productName string) (size int64, err error) {
	sql := sq.Select("count(1)").From(m.productInfo)
	if deviceType != 0 {
		sql = sql.Where("deviceType=?", deviceType)
	}
	if productName != "" {
		sql = sql.Where("productName like ?", "%"+productName+"%")
	}
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
		query := fmt.Sprintf("delete from %s where `productID` = ?", m.productInfo)
		_, err := session.Exec(query, productID)
		if err != nil {
			return err
		}
		query = fmt.Sprintf("delete from %s where `productID` = ?", m.deviceInfo)
		_, err = session.Exec(query, productID)
		if err != nil {
			return err
		}
		query = fmt.Sprintf("delete from %s where `productID` = ?", m.productSchema)
		_, err = session.Exec(query, productID)
		return err
	})
}

func (m *defaultDmModel) Insert(ctx context.Context, pi *ProductInfo, pt *ProductSchema) error {
	return m.Transact(func(session sqlx.Session) error {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.productInfo, productInfoRowsExpectAutoSet)
		_, err := session.Exec(query, pi.ProductID, pi.ProductName, pi.ProductType, pi.AuthMode, pi.DeviceType, pi.CategoryID, pi.NetType, pi.DataProto, pi.AutoRegister, pi.Secret, pi.Description, pi.CreatedTime, pi.UpdatedTime, pi.DeletedTime, pi.DevStatus)
		if err != nil {
			return err
		}
		query = fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.productSchema, productSchemaRowsExpectAutoSet)
		_, err = session.Exec(query, pt.ProductID, pt.Schema, pt.CreatedTime, pt.UpdatedTime, pt.DeletedTime)
		return err
	})
}
