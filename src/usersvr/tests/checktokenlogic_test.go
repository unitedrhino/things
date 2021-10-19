package tests

import (
	"gitee.com/godLei6/things/src/usersvr/internal/logic"
	"gitee.com/godLei6/things/src/usersvr/user"
	"testing"
)

func TestCheckToken(t *testing.T) {
	l := logic.CheckTokenLogic{}
	resp, err := l.CheckToken(&user.CheckTokenReq{
		Token: "123123",
	})
	t.Errorf("TestCheckToken|resp=%#v|err=%#v\n", resp, err)
}
