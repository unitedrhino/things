package def

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
	NotRegistStatus UserStatus = iota //未注册完成状态只注册了第一步
	NomalStatus                       //正常状态
)

type OPT = int64

const (
	OPT_ADD    OPT = 0 //增加
	OPT_MODIFY OPT = 1 //修改
	OPT_DEL    OPT = 2 //删除
)
const UNKNOWN = 0