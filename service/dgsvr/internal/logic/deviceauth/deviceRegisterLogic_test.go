package deviceauthlogic

import "testing"

func Test_getSignature(t *testing.T) {
	type args struct {
		secret string
		dest   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{args: args{
			secret: "HaWo5qNOmisSLb/36oRUfwAY43A=",
			dest:   "deviceName=test&nonce=2428685019&productID=66&timestamp=1756780254",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSignature(tt.args.secret, tt.args.dest); got != tt.want {
				t.Errorf("getSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
