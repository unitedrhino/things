package stores

import (
	"context"
	"database/sql/driver"
	"github.com/i-Things/things/shared/caches"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type ProjectID int64

func (t ProjectID) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) { //更新的时候会调用此接口
	stmt := db.Statement
	need, authIDs, err := caches.GatherUserAuthProjectIDs(ctx)
	if err != nil {
		stmt.Error = err
		return
	}
	expr = clause.Expr{SQL: "?", Vars: []interface{}{int64(t)}}
	if !need { //root 权限不用管
		return
	}
	if !utils.SliceIn(int64(t), authIDs...) { //如果没有权限
		stmt.Error = errors.Permissions.WithMsg("项目权限不足")
	}
	return
}
func (t *ProjectID) Scan(value interface{}) error {
	ret := utils.ToInt64(value)
	p := ProjectID(ret)
	*t = p
	return nil
}

// Value implements the driver Valuer interface.
func (t ProjectID) Value() (driver.Value, error) {
	return int64(t), nil
}

func (t ProjectID) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{AuthQueryClause{Field: f, GetAuthIDs: t.GetAuthIDs(f), AuthDataType: def.AuthDataTypeProject}}
}
func (t ProjectID) CreateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{CreateClause[ProjectID]{Field: f, GetAuthIDs: t.GetAuthIDs(f), AuthDataType: def.AuthDataTypeProject}}
}

func (t ProjectID) DeleteClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{AuthQueryClause{Field: f, GetAuthIDs: t.GetAuthIDs(f), AuthDataType: def.AuthDataTypeProject}}
}

func (t ProjectID) GetAuthIDs(f *schema.Field) GetAuthIDs {
	return func(stmt *gorm.Statement) (authIDs []int64, isRoot bool, err error) {
		need, authIDs, err := caches.GatherUserAuthProjectIDs(stmt.Context)
		if err != nil {
			return nil, false, err
		}
		if !need { //root 权限不用管
			return nil, true, nil
		}
		currProjectID := ctxs.GetMetaProjectID(stmt.Context)
		if currProjectID == 0 { //如果没有选择项目,则给他有的所有项目
			return authIDs, false, nil
		}
		if !utils.SliceIn(currProjectID, authIDs...) { //如果没有权限
			err = errors.Permissions.WithMsg(def.AuthDataTypeFieldTextMap[def.AuthDataTypeProject] + "不足")
			return
		}
		return []int64{currProjectID}, false, nil
	}
}
