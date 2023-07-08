package relationDB

import (
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
)

type ApiAuthRepo struct {
	db *gorm.DB
}

func NewApiAuthRepo(in any) *ApiAuthRepo {
	return &ApiAuthRepo{db: stores.GetCommonConn(in)}
}
