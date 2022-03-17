package tests

import (
	"github.com/i-Things/things/src/usersvr/internal/logic"
	"github.com/i-Things/things/src/usersvr/user"
	"testing"
)

func TestCheckToken(t *testing.T) {
	l := logic.CheckTokenLogic{}
	resp, err := l.CheckToken(&user.CheckTokenReq{
		Token: "123123",
	})
	t.Errorf("TestCheckToken|resp=%#v|err=%#v\n", resp, err)
}
