package logic

import (
	"testing"
	"gitee.com/godLei6/things/src/usersvr/user"
)

func TestCheckToken(t *testing.T) {
	 l := CheckTokenLogic{}
	 resp,err := l.CheckToken(&user.CheckTokenReq{
	 	Token: "123123",
	 })
	 t.Errorf("TestCheckToken|resp=%#v|err=%#v\n",resp,err)
}
