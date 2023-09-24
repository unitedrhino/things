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
	"strconv"
)

// 生产用户数据权限缓存key
func genUserAuthAreaKey(userID int64) string {
	return fmt.Sprintf("user:data:auth:area:userID:%v", userID)
}

// 设置用户数据权限缓存（通用，ctx不限，但需uid传参）
func SetUserAuthArea(ctx context.Context, userID int64, projectID int64, dataIDs []*userDataAuth.Area) error {
	ccJson, err := json.Marshal(dataIDs)
	if err != nil {
		return err
	}
	field := strconv.FormatInt(projectID, 10)
	err = store.HsetCtx(ctx, genUserAuthAreaKey(userID), field, string(ccJson))
	if err != nil {
		return err
	}
	return nil
}

// 读取用户数据权限缓存（通用，ctx不限，但需uid传参）
func GetUserAuthArea(ctx context.Context, userID int64, projectID int64) ([]*userDataAuth.Area, error) {
	field := strconv.FormatInt(int64(projectID), 10)
	ccJson, err := store.HgetCtx(ctx, genUserAuthAreaKey(userID), field)
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, errors.Database.AddDetail(err)
	}
	var dataIDs []*userDataAuth.Area
	err = json.Unmarshal([]byte(ccJson), &dataIDs)
	if err != nil {
		return nil, err
	}
	return dataIDs, nil
}

// 聚合用户数据权限情况
func GatherUserAuthAreaIDs(ctx context.Context) (bool, []int64, error) {
	//检查是否有所有数据权限
	uc := ctxs.GetUserCtxOrNil(ctx)
	if uc == nil || uc.IsAllData || uc.InnerCtx.AllArea {
		return false, nil, nil
	}

	projectID := ctxs.GetMetaProjectID(ctx)
	//读取权限项目ID入参
	var authIDs []int64

	//读取用户数据权限ID
	ccAuthIDs, err := GetUserAuthArea(ctx, uc.UserID, projectID)
	if err != nil {
		return false, nil, err
	}
	if len(ccAuthIDs) == 0 {
		errMsg := "区域权限不足"
		return false, nil, errors.Permissions.WithMsg(errMsg)
	}
	for _, c := range ccAuthIDs {
		authIDs = append(authIDs, cast.ToInt64(c.AreaID))
	}

	return true, authIDs, nil
}
