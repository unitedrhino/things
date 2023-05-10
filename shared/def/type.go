package def

import "golang.org/x/exp/constraints"

type Opt = int64

const (
	OptAdd    Opt = 0 //增加
	OptModify Opt = 1 //修改
	OptDel    Opt = 2 //删除
)
const Unknown = 0

const (
	True  = 1 //是
	False = 2 //否
)

const (
	Enable  = 1 //启用
	Disable = 2 //禁用
)

type CoordinateSystem string

const (
	CoordinateSystemEarth CoordinateSystem = "WGS84" //GPS坐标系：地球系
	CoordinateSystemMars  CoordinateSystem = "GCJ02" //GPS坐标系：火星系
	CoordinateSystemBaidu CoordinateSystem = "BD09"  //GPS坐标系：百度系
)

func ToBool[boolType constraints.Integer](in boolType) bool {
	if in == True {
		return true
	}
	return false
}
func ToIntBool[boolType constraints.Integer](in bool) boolType {
	if in == true {
		return True
	}
	return False
}
