package tenantmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppCreateLogic {
	return &TenantAppCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppCreateLogic) TenantAppCreate(in *sys.TenantAppCreateReq) (*sys.Response, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	uc := ctxs.GetUserCtx(l.ctx)
	uc.AllTenant = true
	defer func() { uc.AllTenant = false }()
	conn := stores.GetTenantConn(l.ctx)
	err := conn.Transaction(func(tx *gorm.DB) error {
		//todo 需要检查租户是否存在
		err := relationDB.NewTenantAppRepo(tx).Insert(l.ctx, &relationDB.SysTenantApp{
			TenantCode: stores.TenantCode(in.Code),
			AppCode:    in.AppCode,
		})
		if err != nil {
			return err
		}
		for _, module := range in.Modules {
			err := ModuleCreate(l.ctx, tx, in.Code, in.AppCode, module.Code, module.MenuIDs)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return &sys.Response{}, err
}
func ModuleCreate(ctx context.Context, tx *gorm.DB, tenantCode, appCode string, moduleCode string, menuIDs []int64) error {
	if _, err := relationDB.NewTenantAppModuleRepo(tx).FindOneByFilter(ctx,
		relationDB.TenantAppModuleFilter{TenantCode: tenantCode, AppCode: appCode, ModuleCodes: []string{moduleCode}}); err == nil || !errors.Cmp(err, errors.NotFind) { //如果报错或者已经有了则跳过
		return err
	}
	mi, err := relationDB.NewModuleInfoRepo(tx).FindOneByFilter(ctx,
		relationDB.ModuleInfoFilter{Codes: []string{moduleCode}, WithApis: true, WithMenus: true})
	if err != nil {
		return err
	}
	var (
		menuMap = make(map[int64]*relationDB.SysModuleMenu)
		//menuTree = make(map[int64]*relationDB.SysModuleMenu)
		allMenu = false
	)
	if len(menuIDs) == 0 {
		allMenu = true
	}
	for _, m := range mi.Menus {
		menuMap[m.ID] = m
		if allMenu {
			menuIDs = append(menuIDs, m.ID)
		}
	}
	var (
		insertMenus []*relationDB.SysTenantAppMenu
	)

	for _, id := range menuIDs {
		m := menuMap[id]
		if m == nil { //模板里不存在无法添加
			continue
		}
		v := relationDB.SysTenantAppMenu{
			TempLateID: m.ID,
			TenantCode: stores.TenantCode(tenantCode),
			AppCode:    appCode, SysModuleMenu: *m}
		v.ID = 0
		insertMenus = append(insertMenus, &v)
	}
	err = relationDB.NewTenantAppMenuRepo(tx).MultiInsert(ctx, insertMenus)
	if err != nil {
		return err
	}
	err = relationDB.NewTenantAppModuleRepo(tx).Insert(ctx, &relationDB.SysTenantAppModule{
		TenantCode: stores.TenantCode(tenantCode), SysAppModule: relationDB.SysAppModule{AppCode: appCode, ModuleCode: moduleCode}})
	return err
}
