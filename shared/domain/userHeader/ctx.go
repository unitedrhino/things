package userHeader

import "context"

type UserCtx struct {
	IsOpen    bool   //是否开放认证用户
	Uid       int64  //用户id（开放认证用户值为0）
	Role      int64  //用户角色（开放认证用户值为0）
	IsAllData bool   //是否所有数据权限（开放认证用户值为true）
	IP        string //用户的ip地址
	Os        string //操作系统
}

func SetUserCtx(ctx context.Context, userCtx *UserCtx) context.Context {
	return context.WithValue(ctx, UserUidKey, userCtx)
}

// 使用该函数前必须传了UserCtx
func GetUserCtx(ctx context.Context) *UserCtx {
	val, ok := ctx.Value(UserUidKey).(*UserCtx)
	if !ok { //这里线上不能获取不到
		panic("GetUserCtx get UserCtx failed")
	}
	return val
}

type MetadataCtx map[string][]string

func SetMetadataCtx(ctx context.Context, maps map[string][]string) context.Context {
	return context.WithValue(ctx, MetadataKey, maps)
}
func GetMetadataCtx(ctx context.Context) *MetadataCtx {
	val, ok := ctx.Value(UserUidKey).(*MetadataCtx)
	if !ok {
		return nil
	}
	return val
}
