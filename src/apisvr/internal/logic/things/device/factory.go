package device

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/domain/things/device"
	"strings"
)

func NewMultiImportCsvRow(cell []string) *device.MultiImportCsvRow {
	return &device.MultiImportCsvRow{
		ProductName: strings.TrimSpace(utils.SliceIndex(cell, 0, "")),
		DeviceName:  strings.TrimSpace(utils.SliceIndex(cell, 1, "")),
		LogLevel:    strings.TrimSpace(utils.SliceIndex(cell, 2, "")),
		Tags:        strings.TrimSpace(utils.SliceIndex(cell, 3, "")),
		Position:    strings.TrimSpace(utils.SliceIndex(cell, 4, "")),
		Address:     strings.TrimSpace(utils.SliceIndex(cell, 5, "")),
	}
}
