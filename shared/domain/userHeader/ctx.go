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
	return context.WithValue(ctx, UserUid, userCtx)
}

// 使用该函数前必须传了UserCtx
func GetUserCtx(ctx context.Context) *UserCtx {
	userCtx, ok := ctx.Value(UserUid).(*UserCtx)
	if !ok {
		//这里线上不能获取不到
		panic("GetUserCtx get UserCtx failed")
	}
	return userCtx
}
