package define

type UserInfoType uint8

const (
	Uid        UserInfoType = iota //用户UID
	InviterUid                     //邀请人用户id
	UserName                       //用户登录名
	GroupId                        //用户组id
	Email                          //邮箱
	Phone                          //手机号
	Wechat						  //微信
	InfoMax						  //结束
	AuthId						  //权限id
)


type UserStatus = int64
const (
	NotRegistStatus UserStatus = iota   //未注册完成状态只注册了第一步
	NomalStatus							//正常状态
)