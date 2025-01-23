package devices

type Auth = int64

const (
	AuthNone   = 0 //啥权限没有
	AuthNormal = 1 //普通功能权限
	AuthSystem = 2 //系统功能权限
	AuthShare  = 3 //只有分享的权限,需要进一步判断
	AuthAll    = 4 //拥有全部权限
)
