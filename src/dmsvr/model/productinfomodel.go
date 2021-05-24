package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	productInfoFieldNames          = builderx.RawFieldNames(&ProductInfo{})
	productInfoRows                = strings.Join(productInfoFieldNames, ",")
	productInfoRowsExpectAutoSet   = strings.Join(stringx.Remove(productInfoFieldNames, "`create_time`", "`update_time`"), ",")
	productInfoRowsWithPlaceHolder = strings.Join(stringx.Remove(productInfoFieldNames, "`productID`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheProductInfoProductIDPrefix   = "cache#productInfo#productID#"
	cacheProductInfoProductNamePrefix = "cache#productInfo#productName#"
)

type (
	ProductInfoModel interface {
		Insert(data ProductInfo) (sql.Result, error)
		FindOne(productID int64) (*ProductInfo, error)
		FindOneByProductName(productName string) (*ProductInfo, error)
		Update(data ProductInfo) error
		Delete(productID int64) error
	}

	defaultProductInfoModel struct {
		sqlc.CachedConn
		table string
	}

	ProductInfo struct {
		ProductID   int64        `db:"productID"`   // 产品id
		ProductName string       `db:"productName"` // 产品名称
		AuthMode    int64        `db:"authMode"`    // 认证方式:0:账密认证,1:秘钥认证
		CreatedTime time.Time    `db:"createdTime"`
		UpdatedTime sql.NullTime `db:"updatedTime"`
		DeletedTime sql.NullTime `db:"deletedTime"`
	}
)

func NewProductInfoModel(conn sqlx.SqlConn, c cache.CacheConf) ProductInfoModel {
	return &defaultProductInfoModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`product_info`",
	}
}

func (m *defaultProductInfoModel) Insert(data ProductInfo) (sql.Result, error) {
	productInfoProductNameKey := fmt.Sprintf("%s%v", cacheProductInfoProductNamePrefix, data.ProductName)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, productInfoRowsExpectAutoSet)
		return conn.Exec(query, data.ProductID, data.ProductName, data.AuthMode, data.CreatedTime, data.UpdatedTime, data.DeletedTime)
	}, productInfoProductNameKey)
	return ret, err
}

func (m *defaultProductInfoModel) FindOne(productID int64) (*ProductInfo, error) {
	productInfoProductIDKey := fmt.Sprintf("%s%v", cacheProductInfoProductIDPrefix, productID)
	var resp ProductInfo
	err := m.QueryRow(&resp, productInfoProductIDKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `productID` = ? limit 1", productInfoRows, m.table)
		return conn.QueryRow(v, query, productID)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultProductInfoModel) FindOneByProductName(productName string) (*ProductInfo, error) {
	productInfoProductNameKey := fmt.Sprintf("%s%v", cacheProductInfoProductNamePrefix, productName)
	var resp ProductInfo
	err := m.QueryRowIndex(&resp, productInfoProductNameKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `productName` = ? limit 1", productInfoRows, m.table)
		if err := conn.QueryRow(&resp, query, productName); err != nil {
			return nil, err
		}
		return resp.ProductID, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultProductInfoModel) Update(data ProductInfo) error {
	productInfoProductIDKey := fmt.Sprintf("%s%v", cacheProductInfoProductIDPrefix, data.ProductID)
	productInfoProductNameKey := fmt.Sprintf("%s%v", cacheProductInfoProductNamePrefix, data.ProductName)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `productID` = ?", m.table, productInfoRowsWithPlaceHolder)
		return conn.Exec(query, data.ProductName, data.AuthMode, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.ProductID)
	}, productInfoProductIDKey, productInfoProductNameKey)
	return err
}

func (m *defaultProductInfoModel) Delete(productID int64) error {
	data, err := m.FindOne(productID)
	if err != nil {
		return err
	}

	productInfoProductIDKey := fmt.Sprintf("%s%v", cacheProductInfoProductIDPrefix, productID)
	productInfoProductNameKey := fmt.Sprintf("%s%v", cacheProductInfoProductNamePrefix, data.ProductName)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `productID` = ?", m.table)
		return conn.Exec(query, productID)
	}, productInfoProductIDKey, productInfoProductNameKey)
	return err
}

func (m *defaultProductInfoModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheProductInfoProductIDPrefix, primary)
}

func (m *defaultProductInfoModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `productID` = ? limit 1", productInfoRows, m.table)
	return conn.QueryRow(v, query, primary)
}
