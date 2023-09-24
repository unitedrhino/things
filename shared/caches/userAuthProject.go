package caches

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/domain/userDataAuth"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils/cast"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// 生产用户数据权限缓存key
func genUserAuthProjectKey(userID int64) string {
	return fmt.Sprintf("user:data:auth:project:userID:%v", userID)
}

// 设置用户数据权限缓存（通用，ctx不限，但需uid传参）
func SetUserAuthProject(ctx context.Context, userID int64, dataIDs []*userDataAuth.Project) error {
	ccJson, err := json.Marshal(dataIDs)
	if err != nil {
		return err
	}
	err = store.SetCtx(ctx, genUserAuthProjectKey(userID), string(ccJson))
	if err != nil {
		return err
	}
	return nil
}

// 读取用户数据权限缓存（通用，ctx不限，但需uid传参）
func GetUserAuthProject(ctx context.Context, userID int64) ([]*userDataAuth.Project, error) {
	ccJson, err := store.GetCtx(ctx, genUserAuthProjectKey(userID))
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, errors.Database.AddDetail(err)
	}
	if ccJson == "" { //没有设置过
		return nil, nil
	}
	var dataIDs []*userDataAuth.Project
	err = json.Unmarshal([]byte(ccJson), &dataIDs)
	if err != nil {
		return nil, err
	}
	return dataIDs, nil
}

// 聚合用户数据权限情况
func GatherUserAuthProjectIDs(ctx context.Context) (bool, []int64, error) {
	//检查是否有所有数据权限
	uc := ctxs.GetUserCtxOrNil(ctx)
	if uc == nil || uc.IsAllData {
		return false, nil, nil
	}
	//读取权限项目ID入参
	var authIDs []int64

	//读取用户数据权限ID
	ccAuthIDs, err := GetUserAuthProject(ctx, uc.UserID)
	if err != nil {
		return false, nil, err
	}
	if len(ccAuthIDs) == 0 {
		errMsg := "项目权限不足"
		return false, nil, errors.Permissions.WithMsg(errMsg)
	}
	for _, c := range ccAuthIDs {
		authIDs = append(authIDs, cast.ToInt64(c.ProjectID))
	}

	return true, authIDs, nil
}
