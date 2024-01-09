package stores

import (
	"context"
	"database/sql/driver"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type DeletedBy int64
type CreatedBy int64
type UpdatedBy int64

func (t CreatedBy) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) { //更新的时候会调用此接口
	expr = clause.Expr{SQL: "?", Vars: []interface{}{int64(t)}}
	return
}

func (t *CreatedBy) Scan(value interface{}) error {
	ret := utils.ToInt64(value)
	p := CreatedBy(ret)
	*t = p
	return nil
}

// Value implements the driver Valuer interface.
func (t CreatedBy) Value() (driver.Value, error) {
	return int64(t), nil
}

func (t CreatedBy) CreateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{UserByClause[CreatedBy]{Field: f}}
}

func (t UpdatedBy) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) { //更新的时候会调用此接口
	expr = clause.Expr{SQL: "?", Vars: []interface{}{int64(t)}}
	return
}

func (t *UpdatedBy) Scan(value interface{}) error {
	ret := utils.ToInt64(value)
	p := UpdatedBy(ret)
	*t = p
	return nil
}

// Value implements the driver Valuer interface.
func (t UpdatedBy) Value() (driver.Value, error) {
	return int64(t), nil
}

func (t UpdatedBy) UpdateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{UserByClause[UpdatedBy]{Field: f}}
}

func (t DeletedBy) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) { //更新的时候会调用此接口
	expr = clause.Expr{SQL: "?", Vars: []interface{}{int64(t)}}
	return
}

func (t *DeletedBy) Scan(value interface{}) error {
	ret := utils.ToInt64(value)
	p := DeletedBy(ret)
	*t = p
	return nil
}

// Value implements the driver Valuer interface.
func (t DeletedBy) Value() (driver.Value, error) {
	return int64(t), nil
}

type UserByClause[keyT ~int64] struct {
	Field *schema.Field
}

func (sd UserByClause[keyT]) Name() string {
	return ""
}

func (sd UserByClause[keyT]) Build(clause.Builder) {
}

func (sd UserByClause[keyT]) MergeClause(*clause.Clause) {
}

func (sd UserByClause[keyT]) ModifyStatement(stmt *gorm.Statement) { //查询的时候会调用此接口
	ctx := stmt.Context
	uc := ctxs.GetUserCtx(ctx)
	if uc == nil {
		return
	}
	var userID = keyT(uc.UserID)
	stmt.SetColumn(sd.Field.DBName, userID, true)
}
