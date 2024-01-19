package stores

import "gorm.io/gorm"

type GetAuthIDs func(stmt *gorm.Statement) (authIDs []int64, err error)
type GetValues func(stmt *gorm.Statement) (values []any, isRoot bool, isAllData bool, err error)
