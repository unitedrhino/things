package model

import (
	"database/sql"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Keys struct {
	Key   string
	Value interface{}
}

type (
	UserModel interface {
		RegisterCore(data UserCore, key Keys) (sql.Result, error)
		GetUserCoreList(page def.PageInfo) ([]*UserCore, int64, error)
	}

	userModel struct {
		sqlc.CachedConn
		userCore string
		userInfo string
	}
)

func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf) UserModel {
	return &userModel{
		CachedConn: sqlc.NewConn(conn, c),
		userCore:   "`user_core`",
		userInfo:   "`user_info`",
	}
}

//插入的时候检查key是否重复
func (m *userModel) RegisterCore(data UserCore, key Keys) (result sql.Result, err error) {
	userCoreUidKey := fmt.Sprintf("%s%v", cacheUserCoreUidPrefix, data.Uid)
	userCoreEmailKey := fmt.Sprintf("%s%v", cacheUserCoreEmailPrefix, data.Email)
	userCorePhoneKey := fmt.Sprintf("%s%v", cacheUserCorePhonePrefix, data.Phone)
	userCoreUserNameKey := fmt.Sprintf("%s%v", cacheUserCoreUserNamePrefix, data.UserName)
	userCoreWechatKey := fmt.Sprintf("%s%v", cacheUserCoreWechatPrefix, data.Wechat)
	m.Transact(func(session sqlx.Session) error {
		var resp UserCore
		var isUpdate bool = false
		query := fmt.Sprintf("select %s from %s where `%s` = ?  limit 1", userCoreRows, m.userCore, key.Key)
		err = session.QueryRow(&resp, query, key.Value)
		if !(err == sqlc.ErrNotFound) {
			if resp.Status == def.NomalStatus {
				return ErrDuplicate
			}
			isUpdate = true
		}
		if isUpdate == true {
			query = fmt.Sprintf("update %s set %s where `uid` = ?", m.userCore, userCoreRowsWithPlaceHolder)
			result, err = session.Exec(query, data.UserName, data.Password, data.Email, data.Phone, data.Wechat, data.LastIP, data.RegIP, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.Status, data.Uid)
		} else {
			query = fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.userCore, userCoreRowsExpectAutoSet)
			result, err = session.Exec(query, data.Uid, data.UserName, data.Password, data.Email, data.Phone, data.Wechat, data.LastIP, data.RegIP, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.Status)
		}
		return err

	})
	m.DelCache(userCoreUidKey, userCoreEmailKey, userCorePhoneKey, userCoreUserNameKey, userCoreWechatKey)
	return result, err
}

//返回 usercore列表,总数及错误信息
func (m *userModel) GetUserCoreList(page def.PageInfo) ([]*UserCore, int64, error) {
	var resp []*UserCore
	query := fmt.Sprintf("select %s from %s  limit %d offset %d ",
		userCoreRows, m.userCore, page.GetLimit(), page.GetOffset())
	err := m.CachedConn.QueryRowsNoCache(&resp, query)
	if err != nil {
		return nil, 0, err
	}

	count := fmt.Sprintf("select count(1) from %s", m.userCore)
	var total int64
	err = m.CachedConn.QueryRowNoCache(&total, count)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, nil
}
