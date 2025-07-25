package schema

import (
	"fmt"
	"gitee.com/unitedrhino/share/utils"
	"strings"
)

func GenArray(identifier string, num int) string {
	return fmt.Sprintf("%s.%d", identifier, num)
}

func GetArray(identifier string) (ident string, num int, ok bool) {
	b, a, ok := strings.Cut(identifier, ".")
	if !ok {
		return identifier, 0, false
	}
	return b, utils.ToInt(a), ok
}
