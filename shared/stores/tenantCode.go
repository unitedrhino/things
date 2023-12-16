package stores

type TenantCode = string

//func (t TenantCode) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) { //更新的时候会调用此接口
//	stmt := db.Statement
//	uc := ctxs.GetUserCtx(ctx)
//	if uc == nil || uc.TenantCode == "" {
//		stmt.Error = errors.Parameter.AddDetail("tenantCode is empty")
//		return
//	}
//	expr = clause.Expr{SQL: "?", Vars: []interface{}{uc.TenantCode}}
//	return
//}
//func (t *TenantCode) Scan(value interface{}) error {
//	ret := cast.ToString(value)
//	p := TenantCode(ret)
//	*t = p
//	return nil
//}
//
//// Value implements the driver Valuer interface.
//func (t TenantCode) Value() (driver.Value, error) {
//	return string(t), nil
//}
//
//func (t TenantCode) QueryClauses(f *schema.Field) []clause.Interface {
//	return []clause.Interface{AuthQueryClause{Field: f, GetAuthIDs: t.GetAuthIDs(f), AuthDataType: def.AuthDataTypeProject}}
//}
//
//func (t TenantCode) CreateClauses(f *schema.Field) []clause.Interface {
//	return []clause.Interface{CreateClause[TenantCode]{Field: f, GetAuthIDs: t.GetAuthIDs(f), AuthDataType: def.AuthDataTypeProject}}
//}
//
//func (t TenantCode) DeleteClauses(f *schema.Field) []clause.Interface {
//	return []clause.Interface{AuthQueryClause{Field: f, GetAuthIDs: t.GetAuthIDs(f), AuthDataType: def.AuthDataTypeProject}}
//}
//
//func (t TenantCode) GetAuthIDs(f *schema.Field) GetAuthIDs {
//	return func(stmt *gorm.Statement) (authIDs []int64, isRoot bool, err error) {
//		need, authIDs, err := caches.GatherUserAuthProjectIDs(stmt.Context)
//		if err != nil {
//			return nil, false, err
//		}
//		if !need { //root 权限不用管
//			return nil, true, nil
//		}
//		currProjectID := ctxs.GetMetaProjectID(stmt.Context)
//		if currProjectID == 0 { //如果没有选择项目,则给他有的所有项目
//			return authIDs, false, nil
//		}
//		if !utils.SliceIn(currProjectID, authIDs...) { //如果没有权限
//			err = errors.Permissions.WithMsg(def.AuthDataTypeFieldTextMap[def.AuthDataTypeProject] + "不足")
//			return
//		}
//		return []int64{currProjectID}, false, nil
//	}
//}
