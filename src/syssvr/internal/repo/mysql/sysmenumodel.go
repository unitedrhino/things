package mysql

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/store"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	MenuModel interface {
		Index(ctx context.Context, in *MenuIndexFilter) ([]*SysMenuInfo, error)
		DeleteMenu(ctx context.Context, menuId int64) error
	}

	menuModel struct {
		sqlc.CachedConn
		menuInfo string
		roleMenu string
	}

	MenuIndexFilter struct {
		Role    int64
		Name    string
		Path    string
		MenuIds []int64
	}
)

func NewMenuModel(conn sqlx.SqlConn, c cache.CacheConf) MenuModel {
	return &menuModel{
		CachedConn: sqlc.NewConn(conn, c),
		menuInfo:   "`sys_menu_info`",
		roleMenu:   "`sys_role_menu`",
	}
}

func (m *MenuIndexFilter) FmtSql(sql sq.SelectBuilder) sq.SelectBuilder {
	//支持账号模糊匹配
	if m.Name != "" {
		sql = sql.Where("`name` like ?", "%"+m.Name+"%")
	}
	if m.Path != "" {
		sql = sql.Where("`path` like ?", "%"+m.Path+"%")
	}
	if len(m.MenuIds) != 0 {
		sql = sql.Where(fmt.Sprintf("`id` in(%v)", store.ArrayToSql(m.MenuIds)))
	}
	return sql
}

func (m *menuModel) Index(ctx context.Context, in *MenuIndexFilter) ([]*SysMenuInfo, error) {
	var resp []*SysMenuInfo
	sql := sq.Select(sysMenuInfoRows).From(m.menuInfo)
	sql = in.FmtSql(sql)
	query, arg, err := sql.ToSql()
	if err != nil {
		return nil, err
	}
	err = m.CachedConn.QueryRowsNoCacheCtx(ctx, &resp, query, arg...)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *menuModel) DeleteMenu(ctx context.Context, MenuId int64) error {
	return m.Transact(func(session sqlx.Session) error {
		//1.从菜单表删除角色
		query := fmt.Sprintf("delete from %s where id = %d",
			m.menuInfo, MenuId)
		_, err := session.Exec(query)
		if err != nil {
			return err
		}
		//2.从角色菜单关系表删除关联菜单项
		query = fmt.Sprintf("delete from %s where  menuID = %d",
			m.roleMenu, MenuId)
		_, err = session.ExecCtx(ctx, query)
		if err != nil {
			return err
		}

		return nil
	})
}
