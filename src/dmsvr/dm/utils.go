package dm

import "gitee.com/godLei6/things/shared/utils"


//字符串类型的产品id有11个字节,不够的需要在前面补0
func GetStrProductID(id int64)string{
	str:= utils.DecimalToAny(id,62)
	return utils.ToLen(str,11)
}


func GetInt64ProductID(id string)int64{
	return utils.AnyToDecimal(id,62)
}