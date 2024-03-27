package main

import "fmt"

func main() {
	str := "{\"actions\":[{\"thenType\":\"\",\"editNumber\":0,\"delay\":300,\"device\":{\"productID\":\"k0sr6p55ktJ\",\"selector\":\"fixed\",\"areaID\":\"1787298960470122496\",\"deviceName\":\"nTJbTQo7axfcDWGWP2h3\",\"deviceAlias\":\"照明设备\",\"schemaAffordance\":\"{\\\"isUseShadow\\\":true,\\\"isNoRecord\\\":false,\\\"define\\\":{\\\"type\\\":\\\"int\\\",\\\"min\\\":\\\"2700\\\",\\\"max\\\":\\\"6500\\\",\\\"start\\\":\\\"2700\\\",\\\"step\\\":\\\"100\\\",\\\"unit\\\":\\\"k\\\"},\\\"mode\\\":\\\"rw\\\"}\",\"dataName\":\"色温\",\"dataID\":\"colorTemp\",\"value\":4524}}]}"
	fmt.Println(str)
}
