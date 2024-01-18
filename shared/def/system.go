package def

const DeviceGroupLevel = 3

type RoleID int64

const (
	RoleIDSuper RoleID = 1
)

type Auth = int64

const (
	AuthRead      Auth = 1 //读权限,只能读,不能写
	AuthWrite     Auth = 2 //写权限，可以写，不能读
	AuthReadWrite Auth = 3 //读写权限
	AuthAdmin     Auth = 4 //管理权限,可以修改别人的权限
)

type CoordinateSystem string

var SystemCoordinate = CoordinateSystemBaidu //默认坐标系

const (
	CoordinateSystemEarth CoordinateSystem = "WGS84" //GPS坐标系：地球系
	CoordinateSystemMars  CoordinateSystem = "GCJ02" //GPS坐标系：火星系
	CoordinateSystemBaidu CoordinateSystem = "BD09"  //GPS坐标系：百度系
)

// 坐标，
type Point struct {
	CoordinateSystem CoordinateSystem `json:"coordinateSystem,omitempty"` //坐标系：WGS84(地球系)，GCJ02(火星系)，BD09(百度系)<br/>参考解释：https://www.cnblogs.com/bigroc/p/16423120.html
	Longitude        float64          `json:"longitude,range=[0:180]"`    //经度
	Latitude         float64          `json:"latitude,range=[0:90]"`      //纬度
}

// 用户数据权限-数据类型
type AuthDataType int64

const (
	AuthDataTypeProject AuthDataType = iota + 1 //项目权限类型
	AuthDataTypeArea                            //区域权限类型
)

var AuthDataTypeFieldTextMap = map[AuthDataType]string{
	AuthDataTypeProject: "项目数据权限",
	AuthDataTypeArea:    "区域数据权限",
}

var AuthDataTypeFieldIDsMap = map[AuthDataType][]string{
	AuthDataTypeProject: {"ProjectID", "ProjectIDs"},
	AuthDataTypeArea:    {"AreaID", "AreaIDs"},
}
