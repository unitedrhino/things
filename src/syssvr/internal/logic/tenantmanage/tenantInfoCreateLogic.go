package tenantmanagelogic

import (
	"context"
	"database/sql"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/shared/utils"
	usermanagelogic "github.com/i-Things/things/src/syssvr/internal/logic/usermanage"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"gorm.io/gorm"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantInfoCreateLogic {
	return &TenantInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 新增区域
func (l *TenantInfoCreateLogic) TenantInfoCreate(in *sys.TenantInfoCreateReq) (*sys.WithID, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	ctxs.GetUserCtx(l.ctx).AllTenant = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllTenant = false
	}()
	userInfo := in.AdminUserInfo
	//首先校验账号格式使用正则表达式，对用户账号做格式校验：只能是大小写字母，数字和下划线，减号
	err := usermanagelogic.CheckUserName(userInfo.UserName)
	if err != nil {
		return nil, err
	}
	//校验密码强度
	err = usermanagelogic.CheckPwd(l.svcCtx, userInfo.Password)
	if err != nil {
		return nil, err
	}
	//1.生成uid
	userID := l.svcCtx.UserID.GetSnowflakeId()
	//2.对密码进行md5加密
	password := utils.MakePwd(userInfo.Password, userID, false)
	ui := relationDB.SysTenantUserInfo{
		TenantCode: stores.TenantCode(in.Info.Code),
		UserID:     userID,
		Phone:      utils.AnyToNullString(userInfo.Phone),
		Email:      utils.AnyToNullString(userInfo.Email),
		UserName:   sql.NullString{String: userInfo.UserName, Valid: true},
		Password:   password,
		NickName:   userInfo.NickName,
		City:       userInfo.City,
		Country:    userInfo.Country,
		Province:   userInfo.Province,
		Language:   userInfo.Language,
		HeadImg:    userInfo.HeadImg,
		Role:       userInfo.Role,
		Sex:        userInfo.Sex,
		IsAllData:  def.True,
	}

	po := ToTenantInfoPo(in.Info)
	err = stores.GetCommonConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		ri := relationDB.SysTenantRoleInfo{TenantCode: stores.TenantCode(in.Info.Code), Name: "超级管理员"}
		err = relationDB.NewRoleInfoRepo(tx).Insert(l.ctx, &ri)
		ui.Role = ri.ID
		err = relationDB.NewUserInfoRepo(tx).Insert(l.ctx, &ui)
		if err != nil {
			return err
		}
		po.AdminUserID = ui.UserID
		err = relationDB.NewTenantInfoRepo(l.ctx).Insert(l.ctx, po)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &sys.WithID{Id: po.ID}, err
}
