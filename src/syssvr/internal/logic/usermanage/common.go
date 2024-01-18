package usermanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/caches"
	"github.com/i-Things/things/shared/domain/userDataAuth"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"regexp"
)

func checkUser(ctx context.Context, userID int64) (*relationDB.SysTenantUserInfo, error) {
	po, err := relationDB.NewUserInfoRepo(ctx).FindOne(ctx, userID)
	if err == nil {
		return po, nil
	}
	if errors.Cmp(err, errors.NotFind) {
		return nil, nil
	}
	return nil, err
}

func InitCacheUserAuthProject(ctx context.Context, userID int64) error {
	projects, err := relationDB.NewUserProjectRepo(ctx).FindByFilter(ctx, relationDB.UserProjectFilter{UserID: userID}, nil)
	if err != nil {
		return err
	}
	return caches.SetUserAuthProject(ctx, userID, DBToAuthProjectDos(projects))
}
func InitCacheUserAuthArea(ctx context.Context, userID int64) error {
	areas, err := relationDB.NewUserAreaRepo(ctx).FindByFilter(ctx, relationDB.UserAreaFilter{UserID: userID}, nil)
	if err != nil {
		return err
	}
	var areaMap = map[int64][]*userDataAuth.Area{}
	for _, v := range areas {
		areaMap[int64(v.ProjectID)] = append(areaMap[int64(v.ProjectID)], DBToAuthAreaDo(v))
	}
	for projectID, areas := range areaMap {
		caches.SetUserAuthArea(ctx, userID, projectID, areas)
	}
	return nil
}

func CheckPwd(svcCtx *svc.ServiceContext, pwd string) error {
	if svcCtx.Config.UserOpt.NeedPassWord &&
		utils.CheckPasswordLever(pwd) < svcCtx.Config.UserOpt.PassLevel {
		return errors.PasswordLevel
	}
	return nil
}
func CheckUserName(userName string) error {
	if ret, _ := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9_-]{6,19}$", userName); !ret {
		return errors.UsernameFormatErr.AddDetail("账号必须以字母开头，且只能包含大小写字母和数字下划线和减号。 长度为6到20位之间")
	}
	return nil
}
