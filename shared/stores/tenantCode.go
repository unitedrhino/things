package stores

import (
	"context"
	"database/sql/driver"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type TenantCode string

func (t TenantCode) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) { //更新的时候会调用此接口
	stmt := db.Statement
	uc := ctxs.GetUserCtx(ctx)
	if uc == nil { //系统初始化的时候会掉用这里
		expr = clause.Expr{SQL: "?", Vars: []interface{}{string(t)}}
		return
	}
	if uc.TenantCode == "" {
		stmt.Error = errors.Parameter.AddDetail("tenantCode is empty")
		return
	}
	expr = clause.Expr{SQL: "?", Vars: []interface{}{uc.TenantCode}}
	return
}
func (t *TenantCode) Scan(value interface{}) error {
	ret := cast.ToString(value)
	p := TenantCode(ret)
	*t = p
	return nil
}

// Value implements the driver Valuer interface.
func (t TenantCode) Value() (driver.Value, error) {
	return string(t), nil
}

func (t TenantCode) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{StringQueryClause{Field: f, GetValues: t.GetAuthIDs(f), Key: "tenant"}}
}

func (t TenantCode) CreateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{StringCreateClause[TenantCode]{Field: f, GetValues: t.GetAuthIDs(f), AuthDataType: def.AuthDataTypeProject}}
}

func (t TenantCode) DeleteClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{StringQueryClause{Field: f, GetValues: t.GetAuthIDs(f), Key: "tenant"}}
}

func (t TenantCode) GetAuthIDs(f *schema.Field) GetValues {
	return func(stmt *gorm.Statement) (authIDs []any, isRoot bool, allData bool, err error) {
		uc := ctxs.GetUserCtx(stmt.Context)
		if uc == nil {
			return nil, false, false, nil
		}
		if uc.TenantCode == def.TenantCodeDefault { //只有core租户的可以修改其他租户的租户号
			isRoot = true
		}
		return []any{TenantCode(uc.TenantCode)}, isRoot, uc.AllTenant, nil
	}
}
