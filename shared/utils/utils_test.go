package utils

import (
	"context"
	"net/http"
	"testing"
)

func TestCheckPasswordLever(t *testing.T) {
	type args struct {
		ps string
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPasswordLever(tt.args.ps); got != tt.want {
				t.Errorf("CheckPasswordLever() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckUserName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckUserName(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("CheckUserName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFuncName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"test1", "utils.TestFuncName.func1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FuncName(); got != tt.want {
				t.Errorf("FuncName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIP(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetIP(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetIP() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandleThrow(t *testing.T) {
	type args struct {
		ctx context.Context
		p   any
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleThrow(tt.args.ctx, tt.args.p)
		})
	}
}

func TestIp2binary(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ip2binary(tt.args.ip); got != tt.want {
				t.Errorf("Ip2binary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmail(tt.args.email); got != tt.want {
				t.Errorf("IsEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMobile(t *testing.T) {
	type args struct {
		mobile string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMobile(tt.args.mobile); got != tt.want {
				t.Errorf("IsMobile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMD5V(t *testing.T) {
	type args struct {
		str []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MD5V(tt.args.str); got != tt.want {
				t.Errorf("MD5V() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakePwd(t *testing.T) {
	type args struct {
		pwd   string
		uid   int64
		isMd5 bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakePwd(tt.args.pwd, tt.args.uid, tt.args.isMd5); got != tt.want {
				t.Errorf("MakePwd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatchIP(t *testing.T) {
	type args struct {
		ip      string
		iprange string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MatchIP(tt.args.ip, tt.args.iprange); got != tt.want {
				t.Errorf("MatchIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecover(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Recover(tt.args.ctx)
		})
	}
}
