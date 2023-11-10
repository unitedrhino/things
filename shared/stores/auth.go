package stores

import "gorm.io/gorm"

type GetAuthIDs func(stmt *gorm.Statement) (authIDs []int64, isRoot bool, err error)
