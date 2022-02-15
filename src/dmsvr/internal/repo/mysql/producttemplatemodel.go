package mysql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	productTemplateFieldNames          = builder.RawFieldNames(&ProductTemplate{})
	productTemplateRows                = strings.Join(productTemplateFieldNames, ",")
	productTemplateRowsExpectAutoSet   = strings.Join(stringx.Remove(productTemplateFieldNames, "`create_time`", "`update_time`"), ",")
	productTemplateRowsWithPlaceHolder = strings.Join(stringx.Remove(productTemplateFieldNames, "`productID`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheDmProductTemplateProductIDPrefix = "cache:dm:productTemplate:productID:"
)

type (
	ProductTemplateModel interface {
		Insert(data *ProductTemplate) (sql.Result, error)
		FindOne(productID string) (*ProductTemplate, error)
		Update(data *ProductTemplate) error
		Delete(productID string) error
	}

	defaultProductTemplateModel struct {
		sqlc.CachedConn
		table string
	}

	ProductTemplate struct {
		ProductID   string       `db:"productID"` // 产品id
		Template    string       `db:"template"`  // 数据模板
		CreatedTime time.Time    `db:"createdTime"`
		UpdatedTime sql.NullTime `db:"updatedTime"`
		DeletedTime sql.NullTime `db:"deletedTime"`
	}
)

func NewProductTemplateModel(conn sqlx.SqlConn, c cache.CacheConf) ProductTemplateModel {
	return &defaultProductTemplateModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`product_template`",
	}
}

func (m *defaultProductTemplateModel) Insert(data *ProductTemplate) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, productTemplateRowsExpectAutoSet)
	ret, err := m.ExecNoCache(query, data.ProductID, data.Template, data.CreatedTime, data.UpdatedTime, data.DeletedTime)

	return ret, err
}

func (m *defaultProductTemplateModel) FindOne(productID string) (*ProductTemplate, error) {
	dmProductTemplateProductIDKey := fmt.Sprintf("%s%v", cacheDmProductTemplateProductIDPrefix, productID)
	var resp ProductTemplate
	err := m.QueryRow(&resp, dmProductTemplateProductIDKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `productID` = ? limit 1", productTemplateRows, m.table)
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

func (m *defaultProductTemplateModel) Update(data *ProductTemplate) error {
	dmProductTemplateProductIDKey := fmt.Sprintf("%s%v", cacheDmProductTemplateProductIDPrefix, data.ProductID)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `productID` = ?", m.table, productTemplateRowsWithPlaceHolder)
		return conn.Exec(query, data.Template, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.ProductID)
	}, dmProductTemplateProductIDKey)
	return err
}

func (m *defaultProductTemplateModel) Delete(productID string) error {

	dmProductTemplateProductIDKey := fmt.Sprintf("%s%v", cacheDmProductTemplateProductIDPrefix, productID)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `productID` = ?", m.table)
		return conn.Exec(query, productID)
	}, dmProductTemplateProductIDKey)
	return err
}

func (m *defaultProductTemplateModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheDmProductTemplateProductIDPrefix, primary)
}

func (m *defaultProductTemplateModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `productID` = ? limit 1", productTemplateRows, m.table)
	return conn.QueryRow(v, query, primary)
}
