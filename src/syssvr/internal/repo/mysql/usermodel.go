package mysql

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
	Value any
}

type (
	UserModel interface {
		RegisterCore(data UserInfo, key Keys) (sql.Result, error)
		Index(page def.PageInfo) ([]*UserInfo, int64, error)
	}

	userModel struct {
		sqlc.CachedConn
		userInfo string
	}
)

func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf) UserModel {
	return &userModel{
		CachedConn: sqlc.NewConn(conn, c),
		userInfo:   "`user_info_test`",
	}
}

//插入的时候检查key是否重复
func (m *userModel) RegisterCore(data UserInfo, key Keys) (result sql.Result, err error) {

	m.Transact(func(session sqlx.Session) error {
		var resp UserInfo
		var isUpdate bool = true
		query := fmt.Sprintf("select %s from %s where `%s` = ?  limit 1", userInfoRows, m.userInfo, key.Key)
		err = session.QueryRow(&resp, query, key.Value)
		if err == sqlc.ErrNotFound {
			isUpdate = false
		}
		if isUpdate == true {
			query = fmt.Sprintf("update %s set %s where `uid` = ?", m.userInfo,
				"`userName`=?,`password`=?,`email`=?,`phone`=?,`wechat`=?,`lastIP`=?,`regIP`=?,`nickName`=?,`sex`=?,`city`=?,`country`=?,`province`=?,`language`=?,`headImgUrl`=?,`role`=?")
			result, err = session.Exec(query, data.UserName, data.Password, data.Email, data.Phone, data.Wechat, data.LastIP, data.RegIP,
				data.NickName, data.Sex, data.City, data.Country, data.Province, data.Language, data.HeadImgUrl, data.Role, data.Uid)
		} else {
			query = fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.userInfo,
				"`uid`,`userName`,`password`,`email`,`phone`,`wechat`,`lastIP`,`regIP`,`nickName`,`sex`,`city`,`country`,`province`,`language`,`headImgUrl`,`role`")
			result, err = session.Exec(query, data.Uid, data.UserName, data.Password, data.Email, data.Phone, data.Wechat, data.LastIP, data.RegIP,
				data.NickName, data.Sex, data.City, data.Country, data.Province, data.Language, data.HeadImgUrl, data.Role)
		}
		return err

	})
	return result, err
}

//返回 usercore列表,总数及错误信息
func (m *userModel) Index(page def.PageInfo) ([]*UserInfo, int64, error) {
	var resp []*UserInfo
	query := fmt.Sprintf("select %s from %s  limit %d offset %d ",
		userInfoRows, m.userInfo, page.GetLimit(), page.GetOffset())
	err := m.CachedConn.QueryRowsNoCache(&resp, query)
	if err != nil {
		return nil, 0, err
	}

	count := fmt.Sprintf("select count(1) from %s", m.userInfo)
	var total int64
	err = m.CachedConn.QueryRowNoCache(&total, count)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, nil
}
