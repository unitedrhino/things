package mysql

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Keys struct {
	Key   string
	Value any
}

type (
	UserModel interface {
		Register(ctx context.Context, UserInfoModel SysUserInfoModel, data SysUserInfo, key Keys) error
		Index(in *UserIndexFilter) ([]*SysUserInfo, int64, error)
	}

	userModel struct {
		sqlc.CachedConn
		userInfo string
	}
	UserIndexFilter struct {
		Page     *def.PageInfo
		UserName string
		Phone    string
		Email    string
	}
)

func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf) UserModel {
	return &userModel{
		CachedConn: sqlc.NewConn(conn, c),
		userInfo:   "`sys_user_info`",
	}
}

//插入的时候检查key是否重复
func (m *userModel) Register(ctx context.Context, UserInfoModel SysUserInfoModel, data SysUserInfo, key Keys) (err error) {

	return m.Transact(func(session sqlx.Session) error {
		var resp SysUserInfo
		var isUpdate bool = true
		query := fmt.Sprintf("select %s from %s where `%s` = ?  limit 1", sysUserInfoRows, m.userInfo, key.Key)
		err = session.QueryRow(&resp, query, key.Value)
		if err == sqlc.ErrNotFound {
			isUpdate = false
		}

		if isUpdate == true {
			err = UserInfoModel.Update(ctx, &data)
			if err != nil {
				return err
			}
		} else {
			_, err = UserInfoModel.Insert(ctx, &data)
			if err != nil {
				return err
			}
		}
		return nil

	})
}

//返回 usercore列表,总数及错误信息
func (m *userModel) Index(in *UserIndexFilter) ([]*SysUserInfo, int64, error) {
	var resp []*SysUserInfo
	page := def.PageInfo{}
	copier.Copy(&page, in.Page)
	//支持账号模糊匹配
	sqlWhere := ""
	if in.UserName != "" {
		sqlWhere += "where userName like '%" + in.UserName + "%'"
	}
	query := fmt.Sprintf("select %s from %s %s limit %d offset %d ",
		sysUserInfoRows, m.userInfo, sqlWhere, page.GetLimit(), page.GetOffset())
	err := m.CachedConn.QueryRowsNoCache(&resp, query)
	if err != nil {
		return nil, 0, err
	}

	count := fmt.Sprintf("select count(1) from %s %s", m.userInfo, sqlWhere)
	var total int64
	err = m.CachedConn.QueryRowNoCache(&total, count)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, nil
}
