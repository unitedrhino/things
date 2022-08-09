package mysql

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserInfoModel = (*customUserInfoModel)(nil)

type (
	// UserInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserInfoModel.
	UserInfoModel interface {
		userInfoModel
		FindOneByPhone(ctx context.Context, phone string) (*UserInfo, error)
		FindOneByWechat(ctx context.Context, weChat string) (*UserInfo, error)
		InsertOrUpdate(ctx context.Context, data UserInfo) error
	}

	customUserInfoModel struct {
		*defaultUserInfoModel
	}
)

// NewUserInfoModel returns a model for the database table.
func NewUserInfoModel(conn sqlx.SqlConn) UserInfoModel {
	return &customUserInfoModel{
		defaultUserInfoModel: newUserInfoModel(conn),
	}
}

func (m *defaultUserInfoModel) FindOneByPhone(ctx context.Context, phone string) (*UserInfo, error) {
	query := fmt.Sprintf("select %s from %s where `phone` = ? limit 1", userInfoRows, m.table)
	var resp UserInfo
	err := m.conn.QueryRowCtx(ctx, &resp, query, phone)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserInfoModel) FindOneByWechat(ctx context.Context, weChat string) (*UserInfo, error) {
	query := fmt.Sprintf("select %s from %s where `weChat` = ? limit 1", userInfoRows, m.table)
	var resp UserInfo
	err := m.conn.QueryRowCtx(ctx, &resp, query, weChat)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserInfoModel) InsertOrUpdate(ctx context.Context, data UserInfo) error {
	_, err := m.FindOne(ctx, data.Uid)
	switch err {
	case nil: //如果找到了直接更新
		err = m.Update(ctx, &data)
	case ErrNotFound: //如果没找到则插入
		_, err = m.Insert(ctx, &data)
	}
	return err
}
