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
	userInfoFieldNames          = builder.RawFieldNames(&UserInfo{})
	userInfoRows                = strings.Join(userInfoFieldNames, ",")
	userInfoRowsExpectAutoSet   = strings.Join(stringx.Remove(userInfoFieldNames, "`create_time`", "`update_time`"), ",")
	userInfoRowsWithPlaceHolder = strings.Join(stringx.Remove(userInfoFieldNames, "`uid`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheUserInfoUidPrefix = "cache#userInfo#uid#"
)

type (
	UserInfoModel interface {
		Insert(data UserInfo) (sql.Result, error)
		InsertOrUpdate(data UserInfo) error
		FindOne(uid int64) (*UserInfo, error)
		Update(data UserInfo) error
		Delete(uid int64) error
	}

	defaultUserInfoModel struct {
		sqlc.CachedConn
		table string
	}

	UserInfo struct {
		Uid         int64        `db:"uid"`        // 用户id
		UserName    string       `db:"userName"`   // 用户名
		NickName    string       `db:"nickName"`   // 用户的昵称
		InviterUid  int64        `db:"inviterUid"` // 邀请人用户id
		InviterId   string       `db:"inviterId"`  // 邀请码
		Sex         int64        `db:"sex"`        // 用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
		City        string       `db:"city"`       // 用户所在城市
		Country     string       `db:"country"`    // 用户所在国家
		Province    string       `db:"province"`   // 用户所在省份
		Language    string       `db:"language"`   // 用户的语言，简体中文为zh_CN
		HeadImgUrl  string       `db:"headImgUrl"` // 用户头像
		CreatedTime sql.NullTime `db:"createdTime"`
		UpdatedTime sql.NullTime `db:"updatedTime"`
		DeletedTime sql.NullTime `db:"deletedTime"`
	}
)

func NewUserInfoModel(conn sqlx.SqlConn, c cache.CacheConf) UserInfoModel {
	return &defaultUserInfoModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`user_info`",
	}
}

func (m *defaultUserInfoModel) Insert(data UserInfo) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, userInfoRowsExpectAutoSet)
	ret, err := m.ExecNoCache(query, data.Uid, data.UserName, data.NickName, data.InviterUid, data.InviterId, data.Sex, data.City, data.Country, data.Province, data.Language, data.HeadImgUrl, data.CreatedTime, data.UpdatedTime, data.DeletedTime)

	return ret, err
}

func (m *defaultUserInfoModel) FindOne(uid int64) (*UserInfo, error) {
	userInfoUidKey := fmt.Sprintf("%s%v", cacheUserInfoUidPrefix, uid)
	var resp UserInfo
	err := m.QueryRow(&resp, userInfoUidKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `uid` = ? limit 1", userInfoRows, m.table)
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

func (m *defaultUserInfoModel) Update(data UserInfo) error {
	data.UpdatedTime = sql.NullTime{Valid: true, Time: time.Now()}
	userInfoUidKey := fmt.Sprintf("%s%v", cacheUserInfoUidPrefix, data.Uid)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `uid` = ?", m.table, userInfoRowsWithPlaceHolder)
		return conn.Exec(query, data.UserName, data.NickName, data.InviterUid, data.InviterId, data.Sex, data.City, data.Country, data.Province, data.Language, data.HeadImgUrl, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.Uid)
	}, userInfoUidKey)
	return err
}

func (m *defaultUserInfoModel) InsertOrUpdate(data UserInfo) error {
	_, err := m.FindOne(data.Uid)
	switch err {
	case nil: //如果找到了直接更新
		err = m.Update(data)
	case ErrNotFound: //如果没找到则插入
		_, err = m.Insert(data)
	}
	return err
}

func (m *defaultUserInfoModel) Delete(uid int64) error {

	userInfoUidKey := fmt.Sprintf("%s%v", cacheUserInfoUidPrefix, uid)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `uid` = ?", m.table)
		return conn.Exec(query, uid)
	}, userInfoUidKey)
	return err
}

func (m *defaultUserInfoModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheUserInfoUidPrefix, primary)
}

func (m *defaultUserInfoModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `uid` = ? limit 1", userInfoRows, m.table)
	return conn.QueryRow(v, query, primary)
}
