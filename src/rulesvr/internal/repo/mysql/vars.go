package mysql

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var ErrNotFound = sqlx.ErrNotFound

func ToError(err error) error {
	if err == ErrNotFound {
		return errors.NotFind.AddMsg("没有找到")
	}
	return errors.Database.AddDetail(err)

}