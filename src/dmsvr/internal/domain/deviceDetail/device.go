package deviceDetail

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/spf13/cast"
	"strings"
)

const (
	PRODUCTID_LEN = 11
)

type LOG_LEVEL = int64

const (
	LOG_CLOSE LOG_LEVEL = 1 //关闭
	LOG_ERROR LOG_LEVEL = 2 //错误
	LOG_WARN  LOG_LEVEL = 3 //告警
	LOG_INFO  LOG_LEVEL = 4 //信息
	LOG_DEBUG LOG_LEVEL = 5 //调试
)

type LoginDevice struct {
	ProductID  string //产品id
	DeviceName string //设备名称
	SdkAppID   int64  //appid 直接填 12010126
	ConnID     string //随机6字节字符串 帮助查bug
	Expiry     int64  //过期时间 unix时间戳
}

func GetLoginDevice(userName string) (*LoginDevice, error) {
	keys := strings.Split(userName, ";")
	if len(keys) != 4 || len(keys[0]) < 11 {
		return nil, errors.Parameter.AddDetail("userName not right")
	}
	lg, err := GetClientIDInfo(keys[0])
	if err != nil {
		return nil, err
	}
	lg.SdkAppID = cast.ToInt64(keys[1])
	lg.ConnID = keys[2]
	lg.Expiry = cast.ToInt64(keys[3])
	return lg, nil
}

func GetClientIDInfo(ClientID string) (*LoginDevice, error) {
	if len(ClientID) < PRODUCTID_LEN {
		return nil, errors.Parameter.AddDetail("clientID length not enough")
	}
	lg := &LoginDevice{
		ProductID:  ClientID[0:PRODUCTID_LEN],
		DeviceName: ClientID[PRODUCTID_LEN:],
	}
	return lg, nil
}

//字符串类型的产品id有11个字节,不够的需要在前面补0
func GetStrProductID(id int64) string {
	str := utils.DecimalToAny(id, 62)
	return utils.ToLen(str, 11)
}

func GetInt64ProductID(id string) int64 {
	return utils.AnyToDecimal(id, 62)
}
