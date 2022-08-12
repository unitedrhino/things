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
	NotRegisterStatus UserStatus = iota //未注册完成状态只注册了第一步
	NormalStatus                        //正常状态
)

type OPT = int64

const (
	OptAdd    OPT = 0 //增加
	OptModify OPT = 1 //修改
	OptDel    OPT = 2 //删除
)
const UNKNOWN = 0

const (
	SUB = "SUB" //订阅
	PUB = "PUB" //发布
)

const (
	PropertyMethod = "property"
	EventMethod    = "event"
	ActionMethod   = "action"

	ReportType = "report" //设备上报的信息
	InfoType   = "info"   //信息
	AlertType  = "alert"  //告警
	FaultType  = "fault"  //故障
)
