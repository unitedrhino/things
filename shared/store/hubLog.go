package store

import "fmt"

type HubLogStore struct {
}

func (h *HubLogStore) GetLogStableName(productID string) string {
	return fmt.Sprintf("`model_hublog_%s`", productID)
}

func (h *HubLogStore) GetLogTableName(productID, deviceName string) string {
	return fmt.Sprintf("`hub_log_%s_%s`", productID, deviceName)
}