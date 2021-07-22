package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"gitee.com/godLei6/things/shared/def"
	"github.com/tal-tech/go-zero/core/logx"
	"net"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
)

func MD5V(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

/*
检测用户名是否符合规范 只可以使用字母数字及下划线 最多30个字符
*/
func CheckUserName(name string) error {
	if len(name) > 30 {
		return errors.New("pwd len more than 30")
	}
	if IsMobile(name) {
		return errors.New("pwd can't be phone number")
	}
	if IsEmail(name) {
		return errors.New("pwd can't be email")
	}
	return nil
}

/*
检测密码是否符合规范 需要至少8位 并且需要包含数字和字母
*/
//密码强度必须为字⺟⼤⼩写+数字+符号，9位以上
func CheckPasswordLever(ps string) int32 {
	level := int32(0)
	if len(ps) < 8 {
		return 0
	}
	num := `[0-9]{1}`
	a_z := `[a-z]{1}`
	A_Z := `[A-Z]{1}`
	symbol := `[!@#~$%^&*()+|_]{1}`

	if b, err := regexp.MatchString(num, ps); b && err == nil {
		level++
	}
	if b, err := regexp.MatchString(a_z, ps); b && err == nil {
		level++
	}
	if b, err := regexp.MatchString(A_Z, ps); b && err == nil {
		level++
	}
	if b, err := regexp.MatchString(symbol, ps); b && err == nil {
		level++
	}
	return level
}

// 识别手机号码
func IsMobile(mobile string) bool {
	result, _ := regexp.MatchString(`^(1[0-9][0-9]\d{4,8})$`, mobile)
	if result {
		return true
	} else {
		return false
	}
}

func IsEmail(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

/*
将密码的md5和uid进行md5
*/
func MakePwd(pwd string, uid int64, isMd5 bool) string {
	if isMd5 == false {
		pwd = MD5V([]byte(pwd))
	}
	strUid := strconv.FormatInt(uid, 8)
	return MD5V([]byte(pwd + strUid + "god17052709767"))
}

func GetLoginNameType(userName string) def.UserInfoType {
	if IsMobile(userName) == true {
		return def.Phone
	}
	return def.UserName
}

// 获取正在运行的函数名
func FuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

//func GetPos()string{
//	pc, file, line, _ := runtime.Caller(2)
//	f := runtime.FuncForPC(pc)
//
//	fmt.Sprintf("%s:%d:%s\n\n\n",file,line,f.Name())
//	return fmt.Sprintf("%s:%d:%s",file,line,f.Name())
//}

func HandleThrow(p interface{}) {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	logx.Errorf("THROW_ERROR|func=%s|error=%#v|stack=%s\n", f, p, string(debug.Stack()))
	os.Exit(-1)
}

func Ip2binary(ip string) string {
	str := strings.Split(ip, ".")
	var ipstr string
	for _, s := range str {
		i, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			fmt.Println(err)
		}
		ipstr = ipstr + fmt.Sprintf("%08b", i)
	}
	return ipstr
}

//测试IP地址和地址端是否匹配 变量ip为字符串，例子"192.168.56.4" iprange为地址端"192.168.56.64/26"
func MatchIP(ip, iprange string) bool {
	ipb := Ip2binary(ip)
	if strings.Contains(iprange, "/") { //如果是ip段
		ipr := strings.Split(iprange, "/")
		masklen, err := strconv.ParseUint(ipr[1], 10, 32)
		if err != nil {
			return false
		}
		iprb := Ip2binary(ipr[0])
		return strings.EqualFold(ipb[0:masklen], iprb[0:masklen])
	} else {
		return ip == iprange
	}

}

// @Summary 获取真实的源ip
func GetIP(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	return "", errors.New("no valid ip found")
}
