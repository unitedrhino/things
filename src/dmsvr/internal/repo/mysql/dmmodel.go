package mysql

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	DmModel interface {
		FindByProductInfo(page def.PageInfo) ([]*ProductInfo, error)
		FindByProductID(productID string, page def.PageInfo) ([]*DeviceInfo, error)
		GetCountByProductID(productID string) (size int64, err error)
		GetCountByProductInfo() (size int64, err error)
		GetDeviceLog(productID, deviceName string, page def.PageInfo2) ([]*DeviceLog, error)
	}

	defaultDmModel struct {
		sqlc.CachedConn
		cache.CacheConf
		productInfo     string
		deviceInfo      string
		productTemplate string
		deviceLog       string
	}
)

func NewDmModel(conn sqlx.SqlConn, c cache.CacheConf) DmModel {
	return &defaultDmModel{
		CachedConn:      sqlc.NewConn(conn, c),
		CacheConf:       c,
		productInfo:     "`product_info`",
		deviceInfo:      "`device_info`",
		productTemplate: "`product_template`",
		deviceLog:       "`device_log`",
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
func (m *defaultDmModel) GetDeviceLog(productID, deviceName string, page def.PageInfo2) ([]*DeviceLog, error) {
	sql := sq.Select(deviceLogFieldNames...).From(m.deviceLog).
		Where(sq.Eq{"productID": productID, "deviceName": deviceName}).
		Limit(uint64(page.GetLimit())).OrderBy("CreatedTime desc")
	if page.TimeStart != 0 {
		sql = sql.Where(sq.GtOrEq{"timestamp": page.TimeStart})
	}
	if page.TimeEnd != 0 {
		sql = sql.Where(sq.LtOrEq{"timestamp": page.TimeEnd})
	}
	sqlStr, value, err := sql.ToSql()
	fmt.Println(sqlStr, value, err)
	var resp []*DeviceLog
	err = m.CachedConn.QueryRowsNoCache(&resp, sqlStr, value...)
	fmt.Println(resp, err)
	return resp, err
}
