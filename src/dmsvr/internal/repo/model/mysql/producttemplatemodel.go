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
	productTemplateFieldNames          = builder.RawFieldNames(&ProductTemplate{})
	productTemplateRows                = strings.Join(productTemplateFieldNames, ",")
	productTemplateRowsExpectAutoSet   = strings.Join(stringx.Remove(productTemplateFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	productTemplateRowsWithPlaceHolder = strings.Join(stringx.Remove(productTemplateFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheDmProductTemplateIdPrefix        = "cache:dm:productTemplate:id:"
	cacheDmProductTemplateProductIDPrefix = "cache:dm:productTemplate:productID:"
)

type (
	ProductTemplateModel interface {
		Insert(data *ProductTemplate) (sql.Result, error)
		FindOne(id int64) (*ProductTemplate, error)
		FindOneByProductID(productID string) (*ProductTemplate, error)
		Update(data *ProductTemplate) error
		Delete(id int64) error
	}

	defaultProductTemplateModel struct {
		sqlc.CachedConn
		table string
	}

	ProductTemplate struct {
		Id          int64
		ProductID   string         // 产品id
		Template    sql.NullString // 数据模板
		CreatedTime time.Time
		UpdatedTime sql.NullTime
		DeletedTime sql.NullTime
	}
)

func NewProductTemplateModel(conn sqlx.SqlConn, c cache.CacheConf) ProductTemplateModel {
	return &defaultProductTemplateModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`product_template`",
	}
}

func (m *defaultProductTemplateModel) Insert(data *ProductTemplate) (sql.Result, error) {
	dmProductTemplateIdKey := fmt.Sprintf("%s%v", cacheDmProductTemplateIdPrefix, data.Id)
	dmProductTemplateProductIDKey := fmt.Sprintf("%s%v", cacheDmProductTemplateProductIDPrefix, data.ProductID)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, productTemplateRowsExpectAutoSet)
		return conn.Exec(query, data.ProductID, data.Template, data.CreatedTime, data.UpdatedTime, data.DeletedTime)
	}, dmProductTemplateIdKey, dmProductTemplateProductIDKey)
	return ret, err
}

func (m *defaultProductTemplateModel) FindOne(id int64) (*ProductTemplate, error) {
	dmProductTemplateIdKey := fmt.Sprintf("%s%v", cacheDmProductTemplateIdPrefix, id)
	var resp ProductTemplate
	err := m.QueryRow(&resp, dmProductTemplateIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", productTemplateRows, m.table)
		return conn.QueryRow(v, query, id)
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

func (m *defaultProductTemplateModel) FindOneByProductID(productID string) (*ProductTemplate, error) {
	dmProductTemplateProductIDKey := fmt.Sprintf("%s%v", cacheDmProductTemplateProductIDPrefix, productID)
	var resp ProductTemplate
	err := m.QueryRowIndex(&resp, dmProductTemplateProductIDKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `productID` = ? limit 1", productTemplateRows, m.table)
		if err := conn.QueryRow(&resp, query, productID); err != nil {
			return nil, err
		}
		return resp.Id, nil
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

func (m *defaultProductTemplateModel) Update(data *ProductTemplate) error {
	dmProductTemplateIdKey := fmt.Sprintf("%s%v", cacheDmProductTemplateIdPrefix, data.Id)
	dmProductTemplateProductIDKey := fmt.Sprintf("%s%v", cacheDmProductTemplateProductIDPrefix, data.ProductID)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, productTemplateRowsWithPlaceHolder)
		return conn.Exec(query, data.ProductID, data.Template, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.Id)
	}, dmProductTemplateIdKey, dmProductTemplateProductIDKey)
	return err
}

func (m *defaultProductTemplateModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	dmProductTemplateIdKey := fmt.Sprintf("%s%v", cacheDmProductTemplateIdPrefix, id)
	dmProductTemplateProductIDKey := fmt.Sprintf("%s%v", cacheDmProductTemplateProductIDPrefix, data.ProductID)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, dmProductTemplateIdKey, dmProductTemplateProductIDKey)
	return err
}

func (m *defaultProductTemplateModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheDmProductTemplateIdPrefix, primary)
}

func (m *defaultProductTemplateModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", productTemplateRows, m.table)
	return conn.QueryRow(v, query, primary)
}
