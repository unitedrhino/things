package mysql

import (
	"errors"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

var ErrNotFound = sqlx.ErrNotFound
var ErrDuplicate = errors.New("sql: duplicate key")
