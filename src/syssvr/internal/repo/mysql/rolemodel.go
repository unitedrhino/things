package mysql

import (
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	RoleModel interface {
		Index(in *sys.RoleIndexReq) ([]*RoleInfo, int64, error)
		IndexRoleIDMenuID(RoleId int64) ([]int64, error)
		UpdateRoleIDMenuID(RoleId int64, MenuId []int64) error
		DeleteRole(RoleId int64) error
	}

	roleModel struct {
		sqlc.CachedConn
		roleInfo string
		roleMenu string
	}
)

func NewRoleModel(conn sqlx.SqlConn, c cache.CacheConf) RoleModel {
	return &roleModel{
		CachedConn: sqlc.NewConn(conn, c),
		roleInfo:   "`role_info`",
		roleMenu:   "`role_menu`",
	}
}

func (m *roleModel) Index(in *sys.RoleIndexReq) ([]*RoleInfo, int64, error) {
	var resp []*RoleInfo
	page := def.PageInfo{}
	copier.Copy(&page, in.Page)
	//支持账号模糊匹配
	sql_where := ""
	if in.Name != "" || in.Status != 0 {
		sql_where += " where "
		if in.Name != "" && in.Status != 0 {
			sql_where += "name like '%" + in.Name + "%' and status=" + cast.ToString(in.Status)
		} else if in.Name != "" && in.Status == 0 {
			sql_where += "name like '%" + in.Name + "%'"
		} else {
			sql_where += "status=" + cast.ToString(in.Status)
		}
	}

	query := fmt.Sprintf("select %s from %s %s limit %d offset %d ",
		roleInfoRows, m.roleInfo, sql_where, page.GetLimit(), page.GetOffset())
	err := m.CachedConn.QueryRowsNoCache(&resp, query)
	if err != nil {
		return nil, 0, err
	}

	count := fmt.Sprintf("select count(1) from %s %s", m.roleInfo, sql_where)
	var total int64
	err = m.CachedConn.QueryRowNoCache(&total, count)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, nil
}

func (m *roleModel) IndexRoleIDMenuID(RoleId int64) ([]int64, error) {
	var resp_menuID []int64
	var resp []*RoleMenu
	query := fmt.Sprintf("select %s from %s where roleID = %d",
		roleMenuRows, m.roleMenu, RoleId)
	err := m.CachedConn.QueryRowsNoCache(&resp, query)
	if err != nil {
		return nil, err
	}
	for _, v := range resp {
		resp_menuID = append(resp_menuID, v.MenuID.Int64)
	}

	return resp_menuID, nil
}

func (m *roleModel) UpdateRoleIDMenuID(RoleId int64, MenuId []int64) error {
	return m.Transact(func(session sqlx.Session) error {
		//1.删除所有roleID对应的menuID
		query := fmt.Sprintf("delete from %s where  roleID = %d",
			m.roleMenu, RoleId)
		_, err := session.Exec(query)
		if err != nil {
			return err
		}

		//2.重新插入menuID
		for _, v := range MenuId {
			query := fmt.Sprintf("insert into %s (roleID, menuID) values (%d, %d)",
				m.roleMenu, RoleId, v)
			_, err := session.Exec(query)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (m *roleModel) DeleteRole(RoleId int64) error {
	return m.Transact(func(session sqlx.Session) error {
		//1.从角色表删除角色
		query := fmt.Sprintf("delete from %s where id = %d",
			m.roleInfo, RoleId)
		_, err := session.Exec(query)
		if err != nil {
			return err
		}
		//2.从角色菜单关系表删除角色
		query = fmt.Sprintf("delete from %s where  roleID = %d",
			m.roleMenu, RoleId)
		_, err = session.Exec(query)
		if err != nil {
			return err
		}

		return nil
	})
}
