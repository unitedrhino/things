package firmware

import (
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func firmwareInfoToApi(v *dm.FirmwareInfo) *types.OtaFirmwareIndex {
	return &types.OtaFirmwareIndex{
		FirmwareID:  v.FirmwareID,
		Name:        v.Name,
		CreatedTime: v.CreatedTime, //创建时间 只读
		ProductID:   v.ProductID,   //产品id 只读
		ProductName: v.ProductName, //产品名称
		Version:     v.Version,
		IsDiff:      v.IsDiff,
		SignMethod:  v.SignMethod,
	}
}
func firmwareReadToApi(v *dm.FirmwareInfoReadResp) *types.OtaFirmwareReadResp {
	dd := make([]*types.OtaFirmwareFileInfo, 0, len(v.Files))
	for _, fileInfo := range v.Files {
		dd = append(dd, &types.OtaFirmwareFileInfo{
			Uri:       fileInfo.FilePath,
			Name:      fileInfo.Name,
			Signature: fileInfo.Signature,
			Size:      fileInfo.Size,
		})
	}
	return &types.OtaFirmwareReadResp{
		FirmwareID:  v.FirmwareID,
		CreatedTime: v.CreatedTime, //创建时间 只读
		Name:        v.Name,
		ProductID:   v.ProductID,   //产品id 只读
		ProductName: v.ProductName, //产品名称
		IsDiff:      v.IsDiff,
		SignMethod:  v.SignMethod,
		Desc:        &v.Desc.Value,
		ExtData:     &v.ExtData.Value,
		Files:       dd,
	}
}
