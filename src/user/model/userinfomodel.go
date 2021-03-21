package model

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	userInfoFieldNames          = builderx.RawFieldNames(&UserInfo{})
	userInfoRows                = strings.Join(userInfoFieldNames, ",")
	userInfoRowsExpectAutoSet   = strings.Join(stringx.Remove(userInfoFieldNames, "`create_time`", "`update_time`"), ",")
	userInfoRowsWithPlaceHolder = strings.Join(stringx.Remove(userInfoFieldNames, "`uid`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	UserInfoModel interface {
		Insert(data UserInfo) (sql.Result, error)
		FindOne(uid int64) (*UserInfo, error)
		Update(data UserInfo) error
		Delete(uid int64) error
	}

	defaultUserInfoModel struct {
		conn  sqlx.SqlConn
		table string
	}

	UserInfo struct {
		Uid         int64          `db:"uid"`        // 用户id
		NickName    sql.NullString `db:"nickName"`   // 用户的昵称
		InviterUid  int64          `db:"inviterUid"` // 邀请人用户id
		InviterId   sql.NullString `db:"inviterId"`  // 邀请码
		Sex         int64          `db:"sex"`        // 用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
		City        sql.NullString `db:"city"`       // 用户所在城市
		Country     sql.NullString `db:"country"`    // 用户所在国家
		Province    sql.NullString `db:"province"`   // 用户所在省份
		Language    sql.NullString `db:"language"`   // 用户的语言，简体中文为zh_CN
		Headimgurl  sql.NullString `db:"headimgurl"` // 用户头像
		CreatedTime sql.NullTime   `db:"createdTime"`
		UpdatedTime sql.NullTime   `db:"updatedTime"`
		DeletedTime sql.NullTime   `db:"deletedTime"`
	}
)

func NewUserInfoModel(conn sqlx.SqlConn) UserInfoModel {
	return &defaultUserInfoModel{
		conn:  conn,
		table: "`user_info`",
	}
}

func (m *defaultUserInfoModel) Insert(data UserInfo) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, userInfoRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.Uid, data.NickName, data.InviterUid, data.InviterId, data.Sex, data.City, data.Country, data.Province, data.Language, data.Headimgurl, data.CreatedTime, data.UpdatedTime, data.DeletedTime)
	return ret, err
}

func (m *defaultUserInfoModel) FindOne(uid int64) (*UserInfo, error) {
	query := fmt.Sprintf("select %s from %s where `uid` = ? limit 1", userInfoRows, m.table)
	var resp UserInfo
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

func (m *defaultUserInfoModel) Update(data UserInfo) error {
	query := fmt.Sprintf("update %s set %s where `uid` = ?", m.table, userInfoRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.NickName, data.InviterUid, data.InviterId, data.Sex, data.City, data.Country, data.Province, data.Language, data.Headimgurl, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.Uid)
	return err
}

func (m *defaultUserInfoModel) Delete(uid int64) error {
	query := fmt.Sprintf("delete from %s where `uid` = ?", m.table)
	_, err := m.conn.Exec(query, uid)
	return err
}
