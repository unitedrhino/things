package users

import "github.com/i-Things/things/shared/utils"

type UserInfoType uint8

const (
	Uid        UserInfoType = iota //用户UID
	InviterUid                     //邀请人用户id
	UserName                       //用户登录名
	GroupId                        //用户组id
	Email                          //邮箱
	Phone                          //手机号
	Wechat                         //微信
	InfoMax                        //结束
	AuthId                         //权限id
)

type UserStatus = int64

const (
	NotRegisterStatus UserStatus = iota //未注册完成状态只注册了第一步
	NormalStatus                        //正常状态
)

func GetLoginNameType(userName string) UserInfoType {
	if utils.IsMobile(userName) == true {
		return Phone
	}
	return UserName
}
