package userHeader

import "context"

type UserCtx struct {
	Uid int64  //用户id
	IP  string //用户的ip地址
}

func SetUserCtx(ctx context.Context, uid int64, IP string) context.Context {
	return context.WithValue(ctx, UserUid, &UserCtx{Uid: uid, IP: IP})
}

//使用该函数前必须传了UserCtx
func GetUserCtx(ctx context.Context) *UserCtx {
	userCtx, ok := ctx.Value(UserUid).(*UserCtx)
	if !ok {
		//这里线上不能获取不到
		panic("GetUserCtx get UserCtx failed")
	}
	return userCtx
}
