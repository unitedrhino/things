package relationDB

import (
	"github.com/i-Things/things/shared/store"
	"gorm.io/gorm"
)

type ApiAuthRepo struct {
	db *gorm.DB
}

func NewApiAuthRepo(in any) *ApiAuthRepo {
	return &ApiAuthRepo{db: store.GetCommonConn(in)}
}
