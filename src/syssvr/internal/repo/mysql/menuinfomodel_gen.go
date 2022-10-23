// Code generated by goctl. DO NOT EDIT!

package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	menuInfoFieldNames          = builder.RawFieldNames(&MenuInfo{})
	menuInfoRows                = strings.Join(menuInfoFieldNames, ",")
	menuInfoRowsExpectAutoSet   = strings.Join(stringx.Remove(menuInfoFieldNames, "`id`", "`updatedTime`", "`deletedTime`", "`createdTime`"), ",")
	menuInfoRowsWithPlaceHolder = strings.Join(stringx.Remove(menuInfoFieldNames, "`id`", "`updatedTime`", "`deletedTime`", "`createdTime`"), "=?,") + "=?"
)

type (
	menuInfoModel interface {
		Insert(ctx context.Context, data *MenuInfo) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*MenuInfo, error)
		FindOneByName(ctx context.Context, name string) (*MenuInfo, error)
		Update(ctx context.Context, data *MenuInfo) error
		Delete(ctx context.Context, id int64) error
	}

	defaultMenuInfoModel struct {
		conn  sqlx.SqlConn
		table string
	}

	MenuInfo struct {
		Id            int64        `db:"id"`            // 编号
		ParentID      int64        `db:"parentID"`      // 父菜单ID，一级菜单为1
		Type          int64        `db:"type"`          // 类型   1：目录   2：菜单   3：按钮
		Order         int64        `db:"order"`         // 左侧table排序序号
		Name          string       `db:"name"`          // 菜单名称
		Path          string       `db:"path"`          // 系统的path
		Component     string       `db:"component"`     // 页面
		Icon          string       `db:"icon"`          // 图标
		Redirect      string       `db:"redirect"`      // 路由重定向
		BackgroundUrl string       `db:"backgroundUrl"` // 后台地址
		HideInMenu    int64        `db:"hideInMenu"`    // 是否隐藏菜单 1-是 2-否
		CreatedTime   time.Time    `db:"createdTime"`   // 创建时间
		UpdatedTime   time.Time    `db:"updatedTime"`   // 更新时间
		DeletedTime   sql.NullTime `db:"deletedTime"`
	}
)

func newMenuInfoModel(conn sqlx.SqlConn) *defaultMenuInfoModel {
	return &defaultMenuInfoModel{
		conn:  conn,
		table: "`menu_info`",
	}
}

func (m *defaultMenuInfoModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultMenuInfoModel) FindOne(ctx context.Context, id int64) (*MenuInfo, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", menuInfoRows, m.table)
	var resp MenuInfo
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultMenuInfoModel) FindOneByName(ctx context.Context, name string) (*MenuInfo, error) {
	var resp MenuInfo
	query := fmt.Sprintf("select %s from %s where `name` = ? limit 1", menuInfoRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, name)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultMenuInfoModel) Insert(ctx context.Context, data *MenuInfo) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, menuInfoRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.ParentID, data.Type, data.Order, data.Name, data.Path, data.Component, data.Icon, data.Redirect, data.BackgroundUrl, data.HideInMenu)
	return ret, err
}

func (m *defaultMenuInfoModel) Update(ctx context.Context, newData *MenuInfo) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, menuInfoRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.ParentID, newData.Type, newData.Order, newData.Name, newData.Path, newData.Component, newData.Icon, newData.Redirect, newData.BackgroundUrl, newData.HideInMenu, newData.Id)
	return err
}

func (m *defaultMenuInfoModel) tableName() string {
	return m.table
}
