package types

import "context"

type UserCtx struct{
	Uid int64  //用户id
}

//使用该函数前必须传了UserCtx
func GetUserCtx(ctx context.Context)*UserCtx{
	userCtx,ok := ctx.Value(USER_UID).(*UserCtx)
	if !ok{
		//这里线上不能获取不到
		panic("GetUserCtx get UserCtx failed")
	}
	return userCtx
}
