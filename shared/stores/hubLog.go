package stores

import "fmt"

type HubLogStore struct {
}

func (h *HubLogStore) GetLogStableName() string {
	return fmt.Sprintf("`model_common_hublog`")
}

func (h *HubLogStore) GetLogTableName(productID, deviceName string) string {
	return fmt.Sprintf("`device_hublog_%s_%s`", productID, deviceName)
}
