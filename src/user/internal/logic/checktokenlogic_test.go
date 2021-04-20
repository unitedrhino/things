package logic

import (
	"testing"
	"yl/src/user/user"
)

func TestCheckToken(t *testing.T) {
	 l := CheckTokenLogic{}
	 resp,err := l.CheckToken(&user.CheckTokenReq{
	 	Token: "123123",
	 })
	 t.Errorf("TestCheckToken|resp=%#v|err=%#v\n",resp,err)
}
