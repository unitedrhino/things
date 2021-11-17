package user

import (
	"gitee.com/godLei6/things/src/usersvr/user"
	"gitee.com/godLei6/things/src/webapi/internal/types"
)


func UserCoreToApi(core *user.UserCore)*types.UserCore{

	return &types.UserCore{
		Uid:                core.Uid,
		UserName:           core.UserName,
		Email:              core.Email,
		Phone:              core.Phone,
		Wechat:             core.Wechat,
		LastIP:             core.LastIP,
		RegIP:              core.RegIP,
		CreatedTime:        core.CreatedTime,
		Status:             core.Status,
	}
}

