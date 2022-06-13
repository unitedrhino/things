package mysql

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	DmModel interface {
		FindByProductInfo(ctx context.Context, page def.PageInfo) ([]*ProductInfo, error)
		FindByProductID(ctx context.Context, productID string, page def.PageInfo) ([]*DeviceInfo, error)
		GetCountByProductID(ctx context.Context, productID string) (size int64, err error)
		GetCountByProductInfo(ctx context.Context) (size int64, err error)
		Insert(ctx context.Context, pi *ProductInfo, pt *ProductTemplate) error
		Delete(ctx context.Context, productID string) error
	}

	defaultDmModel struct {
		sqlc.CachedConn
		cache.CacheConf
		productInfo     string
		deviceInfo      string
		productTemplate string
		deviceLog       string
		ProductInfoModel
	}
)

func NewDmModel(conn sqlx.SqlConn, c cache.CacheConf) DmModel {
	return &defaultDmModel{
		CachedConn:       sqlc.NewConn(conn, c),
		CacheConf:        c,
		productInfo:      "`product_info`",
		deviceInfo:       "`device_info`",
		productTemplate:  "`product_template`",
		deviceLog:        "`device_log`",
		ProductInfoModel: NewProductInfoModel(conn, c),
	}
}

func (m *defaultDmModel) FindByProductID(ctx context.Context, productID string, page def.PageInfo) ([]*DeviceInfo, error) {
	var resp []*DeviceInfo
	query := fmt.Sprintf("select %s from %s where `productID` = ? limit %d offset %d ",
		deviceInfoRows, m.deviceInfo, page.GetLimit(), page.GetOffset())
	err := m.CachedConn.QueryRowsNoCache(&resp, query, productID)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultDmModel) GetCountByProductID(ctx context.Context, productID string) (size int64, err error) {
	query := fmt.Sprintf("select count(1) from %s where `productID` = ?",
		m.deviceInfo)
	err = m.CachedConn.QueryRowNoCache(&size, query, productID)

	switch err {
	case nil:
		return size, nil
	default:
		return 0, err
	}
}

func (m *defaultDmModel) FindByProductInfo(ctx context.Context, page def.PageInfo) ([]*ProductInfo, error) {
	var resp []*ProductInfo
	query := fmt.Sprintf("select %s from %s  limit %d offset %d",
		productInfoRows, m.productInfo, page.GetLimit(), page.GetOffset())
	err := m.QueryRowsNoCache(&resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultDmModel) GetCountByProductInfo(ctx context.Context) (size int64, err error) {
	query := fmt.Sprintf("select count(1)  from %s ",
		m.productInfo)
	err = m.CachedConn.QueryRowNoCache(&size, query)

	switch err {
	case nil:
		return size, nil
	default:
		return 0, err
	}
}

func (m *defaultDmModel) Delete(ctx context.Context, productID string) error {
	data, err := m.FindOne(ctx, productID)
	if err != nil {
		return err
	}
	dmProductTemplateProductIDKey := fmt.Sprintf("%s%v", cacheThingsDmProductTemplateProductIDPrefix, productID)
	dmProductInfoProductIDKey := fmt.Sprintf("%s%v", cacheThingsDmProductInfoProductIDPrefix, productID)
	dmProductInfoProductNameKey := fmt.Sprintf("%s%v", cacheThingsDmProductInfoProductNamePrefix, data.ProductName)
	if err := m.DelCache(dmProductTemplateProductIDKey, dmProductInfoProductIDKey, dmProductInfoProductNameKey); err != nil {
		return err
	}
	return m.Transact(func(session sqlx.Session) error {
		query := fmt.Sprintf("delete from %s where `productID` = ?", m.productInfo)
		_, err := session.Exec(query, productID)
		if err != nil {
			return err
		}
		query = fmt.Sprintf("delete from %s where `productID` = ?", m.productTemplate)
		_, err = session.Exec(query, productID)
		return err
	})
}

func (m *defaultDmModel) Insert(ctx context.Context, pi *ProductInfo, pt *ProductTemplate) error {
	return m.Transact(func(session sqlx.Session) error {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.productInfo, productInfoRowsExpectAutoSet)
		_, err := session.Exec(query, pi.ProductID, pi.ProductName, pi.ProductType, pi.AuthMode, pi.DeviceType, pi.CategoryID, pi.NetType, pi.DataProto, pi.AutoRegister, pi.Secret, pi.Description, pi.CreatedTime, pi.UpdatedTime, pi.DeletedTime, pi.DevStatus)
		if err != nil {
			return err
		}
		query = fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.productTemplate, productTemplateRowsExpectAutoSet)
		_, err = session.Exec(query, pt.ProductID, pt.Template, pt.CreatedTime, pt.UpdatedTime, pt.DeletedTime)
		return err
	})
}
