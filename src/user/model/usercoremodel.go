package model

import (
	"database/sql"
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
	"strings"

	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	userCoreFieldNames          = builderx.RawFieldNames(&UserCore{})
	userCoreRows                = strings.Join(userCoreFieldNames, ",")
	userCoreRowsExpectAutoSet   = strings.Join(stringx.Remove(userCoreFieldNames, "`create_time`", "`update_time`"), ",")
	userCoreRowsWithPlaceHolder = strings.Join(stringx.Remove(userCoreFieldNames, "`uid`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheUserCoreUidPrefix      = "cache#userCore#uid#"
	cacheUserCoreEmailPrefix    = "cache#userCore#email#"
	cacheUserCorePhonePrefix    = "cache#userCore#phone#"
	cacheUserCoreUserNamePrefix = "cache#userCore#userName#"
	cacheUserCoreWechatPrefix   = "cache#userCore#wechat#"
)

type (
	UserCoreModel interface {
		Insert(data UserCore) (sql.Result, error)
		FindOne(uid int64) (*UserCore, error)
		FindOneByEmail(email string) (*UserCore, error)
		FindOneByPhone(phone string) (*UserCore, error)
		FindOneByUserName(userName string) (*UserCore, error)
		FindOneByWechat(wechat string) (*UserCore, error)
		Update(data UserCore) error
		Delete(uid int64) error
	}

	defaultUserCoreModel struct {
		sqlc.CachedConn
		table string
	}

	UserCore struct {
		Uid         int64        `db:"uid"`      // 用户id
		UserName    string       `db:"userName"` // 登录用户名
		Password    string       `db:"password"` // 登录密码
		Email       string       `db:"email"`    // 邮箱
		Phone       string       `db:"phone"`    // 手机号
		Wechat      string       `db:"wechat"`   // 微信UnionID
		LastIP      string       `db:"lastIP"`   // 最后登录ip
		RegIP       string       `db:"regIP"`    // 注册ip
		CreatedTime sql.NullTime `db:"createdTime"`
		UpdatedTime sql.NullTime `db:"updatedTime"`
		DeletedTime sql.NullTime `db:"deletedTime"`
		Status      int64        `db:"status"` // 用户状态:0为未注册状态
	}
)

func NewUserCoreModel(conn sqlx.SqlConn, c cache.CacheConf) UserCoreModel {
	return &defaultUserCoreModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`user_core`",
	}
}

func (m *defaultUserCoreModel) Insert(data UserCore) (sql.Result, error) {
	userCoreWechatKey := fmt.Sprintf("%s%v", cacheUserCoreWechatPrefix, data.Wechat)
	userCoreEmailKey := fmt.Sprintf("%s%v", cacheUserCoreEmailPrefix, data.Email)
	userCorePhoneKey := fmt.Sprintf("%s%v", cacheUserCorePhonePrefix, data.Phone)
	userCoreUserNameKey := fmt.Sprintf("%s%v", cacheUserCoreUserNamePrefix, data.UserName)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, userCoreRowsExpectAutoSet)
		return conn.Exec(query, data.Uid, data.UserName, data.Password, data.Email, data.Phone, data.Wechat, data.LastIP, data.RegIP, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.Status)
	}, userCoreEmailKey, userCorePhoneKey, userCoreUserNameKey, userCoreWechatKey)
	return ret, err
}

func (m *defaultUserCoreModel) FindOne(uid int64) (*UserCore, error) {
	userCoreUidKey := fmt.Sprintf("%s%v", cacheUserCoreUidPrefix, uid)
	var resp UserCore
	err := m.QueryRow(&resp, userCoreUidKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `uid` = ? limit 1", userCoreRows, m.table)
		return conn.QueryRow(v, query, uid)
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

func (m *defaultUserCoreModel) FindOneByEmail(email string) (*UserCore, error) {
	userCoreEmailKey := fmt.Sprintf("%s%v", cacheUserCoreEmailPrefix, email)
	var resp UserCore
	err := m.QueryRowIndex(&resp, userCoreEmailKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `email` = ? limit 1", userCoreRows, m.table)
		if err := conn.QueryRow(&resp, query, email); err != nil {
			return nil, err
		}
		return resp.Uid, nil
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

func (m *defaultUserCoreModel) FindOneByPhone(phone string) (*UserCore, error) {
	userCorePhoneKey := fmt.Sprintf("%s%v", cacheUserCorePhonePrefix, phone)
	var resp UserCore
	err := m.QueryRowIndex(&resp, userCorePhoneKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `phone` = ? limit 1", userCoreRows, m.table)
		if err := conn.QueryRow(&resp, query, phone); err != nil {
			return nil, err
		}
		return resp.Uid, nil
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

func (m *defaultUserCoreModel) FindOneByUserName(userName string) (*UserCore, error) {
	userCoreUserNameKey := fmt.Sprintf("%s%v", cacheUserCoreUserNamePrefix, userName)
	var resp UserCore
	err := m.QueryRowIndex(&resp, userCoreUserNameKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `userName` = ? limit 1", userCoreRows, m.table)
		if err := conn.QueryRow(&resp, query, userName); err != nil {
			return nil, err
		}
		return resp.Uid, nil
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

func (m *defaultUserCoreModel) FindOneByWechat(wechat string) (*UserCore, error) {
	userCoreWechatKey := fmt.Sprintf("%s%v", cacheUserCoreWechatPrefix, wechat)
	var resp UserCore
	err := m.QueryRowIndex(&resp, userCoreWechatKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `wechat` = ? limit 1", userCoreRows, m.table)
		if err := conn.QueryRow(&resp, query, wechat); err != nil {
			return nil, err
		}
		return resp.Uid, nil
	}, m.queryPrimary)
	logx.Errorf("resp=%#v,err=%#v",resp,err)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserCoreModel) Update(data UserCore) error {
	userCoreUidKey := fmt.Sprintf("%s%v", cacheUserCoreUidPrefix, data.Uid)
	userCoreEmailKey := fmt.Sprintf("%s%v", cacheUserCoreEmailPrefix, data.Email)
	userCorePhoneKey := fmt.Sprintf("%s%v", cacheUserCorePhonePrefix, data.Phone)
	userCoreUserNameKey := fmt.Sprintf("%s%v", cacheUserCoreUserNamePrefix, data.UserName)
	userCoreWechatKey := fmt.Sprintf("%s%v", cacheUserCoreWechatPrefix, data.Wechat)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `uid` = ?", m.table, userCoreRowsWithPlaceHolder)
		return conn.Exec(query, data.UserName, data.Password, data.Email, data.Phone, data.Wechat, data.LastIP, data.RegIP, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.Status, data.Uid)
	}, userCoreWechatKey, userCoreUidKey, userCoreEmailKey, userCorePhoneKey, userCoreUserNameKey)
	return err
}

func (m *defaultUserCoreModel) Delete(uid int64) error {
	data, err := m.FindOne(uid)
	if err != nil {
		return err
	}

	userCoreUidKey := fmt.Sprintf("%s%v", cacheUserCoreUidPrefix, uid)
	userCoreEmailKey := fmt.Sprintf("%s%v", cacheUserCoreEmailPrefix, data.Email)
	userCorePhoneKey := fmt.Sprintf("%s%v", cacheUserCorePhonePrefix, data.Phone)
	userCoreUserNameKey := fmt.Sprintf("%s%v", cacheUserCoreUserNamePrefix, data.UserName)
	userCoreWechatKey := fmt.Sprintf("%s%v", cacheUserCoreWechatPrefix, data.Wechat)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `uid` = ?", m.table)
		return conn.Exec(query, uid)
	}, userCoreEmailKey, userCorePhoneKey, userCoreUserNameKey, userCoreWechatKey, userCoreUidKey)
	return err
}

func (m *defaultUserCoreModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheUserCoreUidPrefix, primary)
}

func (m *defaultUserCoreModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `uid` = ? limit 1", userCoreRows, m.table)
	return conn.QueryRow(v, query, primary)
}
