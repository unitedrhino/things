package stores

import (
	"context"
	"database/sql/driver"
	"github.com/i-Things/things/shared/caches"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/shared/utils/cast"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type AreaID int64

func (t AreaID) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) { //更新的时候会调用此接口
	stmt := db.Statement
	need, authIDs, err := caches.GatherUserAuthAreaIDs(ctx)
	if err != nil {
		stmt.Error = err
		return
	}
	expr = clause.Expr{SQL: "?", Vars: []interface{}{int64(t)}}
	if !need { //root 权限不用管
		return
	}
	if !utils.SliceIn(int64(t), authIDs...) { //如果没有权限
		stmt.Error = errors.Permissions.WithMsg("区域权限不足")
	}
	return
}
func (t *AreaID) Scan(value interface{}) error {
	ret := cast.ToInt64(value)
	p := AreaID(ret)
	*t = p
	return nil
}

// Value implements the driver Valuer interface.
func (t AreaID) Value() (driver.Value, error) {
	return int64(t), nil
}

func (t AreaID) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{AuthQueryClause{Field: f, GetAuthIDs: t.GetAuthIDs(f), AuthDataType: def.AuthDataTypeArea}}
}

func (t AreaID) CreateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{CreateClause[AreaID]{Field: f, GetAuthIDs: t.GetAuthIDs(f), AuthDataType: def.AuthDataTypeArea}}
}

func (t AreaID) DeleteClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{AuthQueryClause{Field: f, GetAuthIDs: t.GetAuthIDs(f), AuthDataType: def.AuthDataTypeArea}}
}

func (t AreaID) GetAuthIDs(f *schema.Field) GetAuthIDs {
	return func(stmt *gorm.Statement) (authIDs []int64, isRoot bool, err error) {
		need, authIDs, err := caches.GatherUserAuthAreaIDs(stmt.Context)
		if err != nil {
			return nil, false, err
		}
		if !need { //root 权限不用管
			return nil, true, nil
		}
		return authIDs, false, nil
	}
}
