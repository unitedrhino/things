package mysql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/stores/builder"
	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
)

var (
	productInfoFieldNames          = builder.RawFieldNames(&ProductInfo{})
	productInfoRows                = strings.Join(productInfoFieldNames, ",")
	productInfoRowsExpectAutoSet   = strings.Join(stringx.Remove(productInfoFieldNames, "`create_time`", "`update_time`"), ",")
	productInfoRowsWithPlaceHolder = strings.Join(stringx.Remove(productInfoFieldNames, "`productID`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheDmProductInfoProductIDPrefix   = "cache:dm:productInfo:productID:"
	cacheDmProductInfoProductNamePrefix = "cache:dm:productInfo:productName:"
)

type (
	ProductInfoModel interface {
		Insert(data *ProductInfo) (sql.Result, error)
		FindOne(productID string) (*ProductInfo, error)
		FindOneByProductName(productName string) (*ProductInfo, error)
		Update(data *ProductInfo) error
		Delete(productID string) error
	}

	defaultProductInfoModel struct {
		sqlc.CachedConn
		table string
	}

	ProductInfo struct {
		ProductID    string // 产品id
		ProductName  string // 产品名称
		ProductType  int64  // 产品状态:0:开发中,1:审核中,2:已发布
		AuthMode     int64  // 认证方式:0:账密认证,1:秘钥认证
		DeviceType   int64  // 设备类型:0:设备,1:网关,2:子设备
		CategoryID   int64  // 产品品类
		NetType      int64  // 通讯方式:0:其他,1:wi-fi,2:2G/3G/4G,3:5G,4:BLE,5:LoRaWAN
		DataProto    int64  // 数据协议:0:自定义,1:数据模板
		AutoRegister int64  // 动态注册:0:关闭,1:打开,2:打开并自动创建设备
		Secret       string // 动态注册产品秘钥
		Description  string // 描述
		CreatedTime  time.Time
		UpdatedTime  sql.NullTime
		DeletedTime  sql.NullTime
		DevStatus    string // 产品状态
	}
)

func NewProductInfoModel(conn sqlx.SqlConn, c cache.CacheConf) ProductInfoModel {
	return &defaultProductInfoModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`product_info`",
	}
}

func (m *defaultProductInfoModel) Insert(data *ProductInfo) (sql.Result, error) {
	dmProductInfoProductNameKey := fmt.Sprintf("%s%v", cacheDmProductInfoProductNamePrefix, data.ProductName)
	dmProductInfoProductIDKey := fmt.Sprintf("%s%v", cacheDmProductInfoProductIDPrefix, data.ProductID)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, productInfoRowsExpectAutoSet)
		return conn.Exec(query, data.ProductID, data.ProductName, data.ProductType, data.AuthMode, data.DeviceType, data.CategoryID, data.NetType, data.DataProto, data.AutoRegister, data.Secret, data.Description, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.DevStatus)
	}, dmProductInfoProductIDKey, dmProductInfoProductNameKey)
	return ret, err
}

func (m *defaultProductInfoModel) FindOne(productID string) (*ProductInfo, error) {
	dmProductInfoProductIDKey := fmt.Sprintf("%s%v", cacheDmProductInfoProductIDPrefix, productID)
	var resp ProductInfo
	err := m.QueryRow(&resp, dmProductInfoProductIDKey, func(conn sqlx.SqlConn, v interface{}) error {
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
	dmProductInfoProductNameKey := fmt.Sprintf("%s%v", cacheDmProductInfoProductNamePrefix, productName)
	var resp ProductInfo
	err := m.QueryRowIndex(&resp, dmProductInfoProductNameKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
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

func (m *defaultProductInfoModel) Update(data *ProductInfo) error {
	dmProductInfoProductIDKey := fmt.Sprintf("%s%v", cacheDmProductInfoProductIDPrefix, data.ProductID)
	dmProductInfoProductNameKey := fmt.Sprintf("%s%v", cacheDmProductInfoProductNamePrefix, data.ProductName)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `productID` = ?", m.table, productInfoRowsWithPlaceHolder)
		return conn.Exec(query, data.ProductName, data.ProductType, data.AuthMode, data.DeviceType, data.CategoryID, data.NetType, data.DataProto, data.AutoRegister, data.Secret, data.Description, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.DevStatus, data.ProductID)
	}, dmProductInfoProductIDKey, dmProductInfoProductNameKey)
	return err
}

func (m *defaultProductInfoModel) Delete(productID string) error {
	data, err := m.FindOne(productID)
	if err != nil {
		return err
	}

	dmProductInfoProductIDKey := fmt.Sprintf("%s%v", cacheDmProductInfoProductIDPrefix, productID)
	dmProductInfoProductNameKey := fmt.Sprintf("%s%v", cacheDmProductInfoProductNamePrefix, data.ProductName)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `productID` = ?", m.table)
		return conn.Exec(query, productID)
	}, dmProductInfoProductIDKey, dmProductInfoProductNameKey)
	return err
}

func (m *defaultProductInfoModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheDmProductInfoProductIDPrefix, primary)
}

func (m *defaultProductInfoModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `productID` = ? limit 1", productInfoRows, m.table)
	return conn.QueryRow(v, query, primary)
}
