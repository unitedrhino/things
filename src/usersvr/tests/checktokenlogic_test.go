package tests

import (
	"github.com/go-things/things/src/usersvr/internal/logic"
	"github.com/go-things/things/src/usersvr/user"
	"testing"
)

func TestCheckToken(t *testing.T) {
	l := logic.CheckTokenLogic{}
	resp, err := l.CheckToken(&user.CheckTokenReq{
		Token: "123123",
	})
	t.Errorf("TestCheckToken|resp=%#v|err=%#v\n", resp, err)
}
