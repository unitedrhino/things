package users

import "github.com/i-Things/things/shared/utils"

// phone 手机号 wxOpen 微信开放平台 wxIn 微信内 wxMiniP 微信小程序 pwd 账号密码
const (
	RegEmail   = "email"   //邮箱
	RegPhone   = "phone"   //手机号
	RegWxOpen  = "wxOpen"  //微信开放平台登录
	RegWxIn    = "wxIn"    //微信内登录
	RegWxMiniP = "wxMiniP" //微信小程序
	RegPwd     = "pwd"     //账号密码注册
)

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
