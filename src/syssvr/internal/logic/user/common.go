package userlogic

import (
	"context"
	"github.com/i-Things/things/shared/caches"
	"github.com/i-Things/things/shared/domain/userDataAuth"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
)

func checkUser(ctx context.Context, userID int64) (*relationDB.SysUserInfo, error) {
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
	projects, err := relationDB.NewUserAuthProjectRepo(ctx).FindByFilter(ctx, relationDB.UserAuthProjectFilter{UserID: userID}, nil)
	if err != nil {
		return err
	}
	return caches.SetUserAuthProject(ctx, userID, DBToAuthProjectDos(projects))
}
func InitCacheUserAuthArea(ctx context.Context, userID int64) error {
	areas, err := relationDB.NewUserAuthAreaRepo(ctx).FindByFilter(ctx, relationDB.UserAuthAreaFilter{UserID: userID}, nil)
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
