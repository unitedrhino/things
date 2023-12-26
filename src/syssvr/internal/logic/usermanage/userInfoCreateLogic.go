package usermanagelogic

import (
	"context"
	"database/sql"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UserInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UiDB *relationDB.UserInfoRepo
}

func NewUserInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoCreateLogic {
	return &UserInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UiDB:   relationDB.NewUserInfoRepo(ctx),
	}
}

func (l *UserInfoCreateLogic) UserInfoInsert(in *sys.UserInfoCreateReq) (int64, error) {
	info := in.Info
	var userID int64
	//首先校验账号格式使用正则表达式，对用户账号做格式校验：只能是大小写字母，数字和下划线，减号
	err := CheckUserName(info.UserName)
	if err != nil {
		return 0, err
	}
	//校验密码强度
	err = CheckPwd(l.svcCtx, info.Password)
	if info.Role == 0 {
		info.Role = in.RoleIDs[0]
	} else if !utils.SliceIn(info.Role, in.RoleIDs...) {
		return 0, errors.Parameter.AddMsgf("用户默认角色不存在")
	}
	count, err := relationDB.NewRoleInfoRepo(l.ctx).CountByFilter(l.ctx, relationDB.RoleInfoFilter{IDs: in.RoleIDs})
	if err != nil {
		return 0, err
	}
	if int(count) != len(in.RoleIDs) {
		return 0, errors.Parameter.AddMsgf("角色有不存在的")
	}
	err = stores.GetCommonConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		uidb := relationDB.NewUserInfoRepo(tx)
		var account = []string{info.UserName}
		if info.Phone != "" {
			account = append(account, info.Phone)
		}
		if info.Email != "" {
			account = append(account, info.Email)
		}
		_, err = uidb.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Accounts: account})
		if err == nil { //已注册
			//提示重复注册
			return errors.DuplicateRegister
		}
		if !errors.Cmp(err, errors.NotFind) {
			return err
		}
		//1.生成uid
		userID = l.svcCtx.UserID.GetSnowflakeId()

		//2.对密码进行md5加密
		password := utils.MakePwd(info.Password, userID, false)
		ui := relationDB.SysTenantUserInfo{
			UserID:    userID,
			UserName:  sql.NullString{String: info.UserName, Valid: true},
			Password:  password,
			NickName:  info.NickName,
			City:      info.City,
			Country:   info.Country,
			Province:  info.Province,
			Language:  info.Language,
			HeadImg:   info.HeadImg,
			Role:      info.Role,
			Sex:       info.Sex,
			IsAllData: info.IsAllData,
		}
		if info.Email != "" {
			ui.Email = sql.NullString{String: info.Email, Valid: true}
		}
		if info.Phone != "" {
			ui.Phone = sql.NullString{String: info.Phone, Valid: true}
		}
		err = uidb.Insert(l.ctx, &ui)
		if err != nil { //并发情况下有可能重复所以需要再次判断一次
			if errors.Cmp(err, errors.NotFind) {
				return errors.DuplicateUsername.AddDetail(info.UserName)
			}
			l.Errorf("%s.Inserts err=%#v", utils.FuncName(), err)
			return err
		}
		err := relationDB.NewUserRoleRepo(tx).MultiUpdate(l.ctx, ui.UserID, in.RoleIDs)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return userID, nil
}
func (l *UserInfoCreateLogic) UserInfoCreate(in *sys.UserInfoCreateReq) (*sys.UserCreateResp, error) {
	l.Infof("%s req=%+v", utils.FuncName(), in)
	userID, err := l.UserInfoInsert(in)
	if err != nil {
		return nil, err
	}
	return &sys.UserCreateResp{UserID: userID}, nil
}
