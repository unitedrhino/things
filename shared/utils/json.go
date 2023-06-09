package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func Unmarshal(data []byte, v any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	return decoder.Decode(v)
}

func UnmarshalNoErr[inT any](data string) inT {
	var ret inT
	json.Unmarshal([]byte(data), &ret)
	return ret
}

func MarshalNoErr(v any) string {
	ret, _ := json.Marshal(v)
	return string(ret)
}

// Fmt 将结构以更方便看的方式打印出来
func Fmt(v any) string {
	switch v.(type) {
	case string:
		return v.(string)
	case []byte:
		return string(v.([]byte))
	case error:
		return v.(error).Error()
	case interface{ String() string }:
		return v.(interface{ String() string }).String()
	default:
		js, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%#v", js)
		}
		return string(js)
	}
}
