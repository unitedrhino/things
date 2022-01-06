package dm

import (
	"github.com/go-things/things/shared/errors"
	"github.com/spf13/cast"
	"strings"
)

type LoginDevice struct {
	ClientID   string //clientID
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
		ClientID:   ClientID,
		ProductID:  ClientID[0:PRODUCTID_LEN],
		DeviceName: ClientID[PRODUCTID_LEN:],
	}
	return lg, nil
}
