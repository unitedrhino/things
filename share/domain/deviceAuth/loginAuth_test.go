package deviceAuth

import (
	"fmt"
	"testing"
)

func TestNewPwdInfo(t *testing.T) {
	type args struct {
		signature  string
		signMethod string
	}
	ag := args{
		signature:  "a5bf92adfd52c7f0f39e9d4745b57356828f65b4",
		signMethod: "hmacsha1",
	}
	got, err := NewPwdInfo(ag.signature, ag.signMethod)
	fmt.Println(got, err)
	sign := fmt.Sprintf("%v%v;%v;%v", "01F", "subdevice1", 256, 1728446804)

	err = got.CmpPwd(sign, "mzxe12OY8z/im7S3DNhHsCXdB4o=")
	fmt.Println(err)
}
