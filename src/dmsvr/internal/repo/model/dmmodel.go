package model

import (
	"fmt"
	"gitee.com/godLei6/things/shared/def"
	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

type (
	DmModel interface {
		FindByProductInfo(page def.PageInfo) ([]*ProductInfo, error)
		FindByProductID(productID string, page def.PageInfo) ([]*DeviceInfo, error)
		GetCountByProductID(productID string) (size int64, err error)
		GetCountByProductInfo() (size int64, err error)
	}

	defaultDmModel struct {
		sqlc.CachedConn
		productInfo string
		deviceInfo  string
	}
)

func NewDmModel(conn sqlx.SqlConn, c cache.CacheConf) DmModel {
	return &defaultDmModel{
		CachedConn:  sqlc.NewConn(conn, c),
		productInfo: "`product_info`",
		deviceInfo:  "device_info",
	}
}

func (m *defaultDmModel) FindByProductID(productID string, page def.PageInfo) ([]*DeviceInfo, error) {
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

func (m *defaultDmModel) GetCountByProductID(productID string) (size int64, err error) {
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

func (m *defaultDmModel) FindByProductInfo(page def.PageInfo) ([]*ProductInfo, error) {
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

func (m *defaultDmModel) GetCountByProductInfo() (size int64, err error) {
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
