package errors_test

import (
	"encoding/json"
	"fmt"
	"github.com/go-things/things/shared/errors"
	"strings"
	"testing"
)

func TestGetDetailMsg(t *testing.T) {
	err := errors.OK.AddDetail("detail1").AddDetail("detail2").GetDetailMsg()
	fmt.Println(err)
}

func TestJson(t *testing.T) {
	str := `
	{
		"bit7":72057594037927935,
		"bit2":123123,
		"str":"72057594037927935",
		"struct":{
			"stru1":123,
			"stru2":"faewfae"
		},
		"array":[
			"afw2","fwefa"
		],
		"bool":true,
		"string":"dfawefa"
	}
`
	mapStr := make(map[string]interface{}, 10)
	decoder := json.NewDecoder(strings.NewReader(str))
	decoder.UseNumber()
	err := decoder.Decode(&mapStr)

	out, _ := json.Marshal(mapStr)
	fmt.Printf("err=%v,str=%#v\n|bit7=%T|bit2=%T|struct=%T|array=%T|bool=%T|string=%T,out=%s\n",
		err, mapStr, mapStr["bit7"], mapStr["bit2"], mapStr["struct"], mapStr["array"],
		mapStr["bool"], mapStr["string"], string(out))
}
