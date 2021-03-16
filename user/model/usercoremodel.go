package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

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
)

type (
	UserCoreModel interface {
		Insert(data UserCore) (sql.Result, error)
		FindOne(uid int64) (*UserCore, error)
		FindOneByEmail(email sql.NullString) (*UserCore, error)
		FindOneByPhone(phone sql.NullString) (*UserCore, error)
		FindOneByUserName(userName sql.NullString) (*UserCore, error)
		FindOneByWechat(wechat sql.NullString) (*UserCore, error)
		Update(data UserCore) error
		Delete(uid int64) error
	}

	defaultUserCoreModel struct {
		conn  sqlx.SqlConn
		table string
	}

	UserCore struct {
		Uid         int64          `db:"uid"`      // 用户id
		UserName    sql.NullString `db:"userName"` // 登录用户名
		Password    sql.NullString `db:"password"` // 登录密码
		Email       sql.NullString `db:"email"`    // 邮箱
		Phone       sql.NullString `db:"phone"`    // 手机号
		Wechat      sql.NullString `db:"wechat"`   // 微信openId
		LastIP      sql.NullString `db:"lastIP"`   // 最后登录ip
		RegIP       sql.NullString `db:"regIP"`    // 注册ip
		CreatedTime sql.NullTime   `db:"createdTime"`
		UpdatedTime sql.NullTime   `db:"updatedTime"`
		DeletedTime sql.NullTime   `db:"deletedTime"`
		Status      int64          `db:"status"` // 用户状态:0为未注册状态
	}
)

func NewUserCoreModel(conn sqlx.SqlConn) UserCoreModel {
	return &defaultUserCoreModel{
		conn:  conn,
		table: "`user_core`",
	}
}

func (m *defaultUserCoreModel) Insert(data UserCore) (sql.Result, error) {
	data.CreatedTime = sql.NullTime{
		Time: time.Now(),
		Valid: true,
	}
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, userCoreRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.Uid, data.UserName, data.Password, data.Email, data.Phone, data.Wechat, data.LastIP, data.RegIP, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.Status)
	return ret, err
}

func (m *defaultUserCoreModel) FindOne(uid int64) (*UserCore, error) {
	query := fmt.Sprintf("select %s from %s where `uid` = ? limit 1", userCoreRows, m.table)
	var resp UserCore
	err := m.conn.QueryRow(&resp, query, uid)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserCoreModel) FindOneByEmail(email sql.NullString) (*UserCore, error) {
	var resp UserCore
	query := fmt.Sprintf("select %s from %s where `email` = ? limit 1", userCoreRows, m.table)
	err := m.conn.QueryRow(&resp, query, email)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserCoreModel) FindOneByPhone(phone sql.NullString) (*UserCore, error) {
	var resp UserCore
	query := fmt.Sprintf("select %s from %s where `phone` = ? limit 1", userCoreRows, m.table)
	err := m.conn.QueryRow(&resp, query, phone)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserCoreModel) FindOneByUserName(userName sql.NullString) (*UserCore, error) {
	var resp UserCore
	query := fmt.Sprintf("select %s from %s where `userName` = ? limit 1", userCoreRows, m.table)
	err := m.conn.QueryRow(&resp, query, userName)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserCoreModel) FindOneByWechat(wechat sql.NullString) (*UserCore, error) {
	var resp UserCore
	query := fmt.Sprintf("select %s from %s where `wechat` = ? limit 1", userCoreRows, m.table)
	err := m.conn.QueryRow(&resp, query, wechat)
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
	data.UpdatedTime = sql.NullTime{
		Time: time.Now(),Valid: true,
	}
	query := fmt.Sprintf("update %s set %s where `uid` = ?", m.table, userCoreRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.UserName, data.Password, data.Email, data.Phone, data.Wechat, data.LastIP, data.RegIP, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.Status, data.Uid)
	return err
}

func (m *defaultUserCoreModel) Delete(uid int64) error {
	query := fmt.Sprintf("delete from %s where `uid` = ?", m.table)
	_, err := m.conn.Exec(query, uid)
	return err
}
