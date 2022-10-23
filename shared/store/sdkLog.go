package store

import "fmt"

type SDKLogStore struct {
}

func (s *SDKLogStore) GetSDKLogStableName() string {
	return fmt.Sprintf("`model_common_sdklog`")
}

func (s *SDKLogStore) GetSDKLogTableName(productID, deviceName string) string {
	return fmt.Sprintf("`device_sdklog_%s_%s`", productID, deviceName)
}
