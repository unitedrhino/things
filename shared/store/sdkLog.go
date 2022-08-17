package store

import "fmt"

type SDKLogStore struct {
}

func (s *SDKLogStore) GetSDKLogStableName(productID string) string {
	return fmt.Sprintf("`model_sdklog_%s`", productID)
}

func (s *SDKLogStore) GetSDKLogTableName(productID, deviceName string) string {
	return fmt.Sprintf("`sdk_log_%s_%s`", productID, deviceName)
}
