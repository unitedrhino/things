package def

type Opt = int64

const (
	OptAdd    Opt = 0 //增加
	OptModify Opt = 1 //修改
	OptDel    Opt = 2 //删除
)
const Unknown = 0
const (
	OffLine = 1 //离线
	OnLine  = 2 //在线
)
