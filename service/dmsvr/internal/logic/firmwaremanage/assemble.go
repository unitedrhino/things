package firmwaremanagelogic

import (
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func ToFirmwareInfo(di *relationDB.DmOtaFirmware, df ...*relationDB.DmOtaFirmwareFile) *dm.FirmwareInfo {
	dd := make([]*dm.OtaFirmwareFile, 0, len(df))
	for _, fileInfo := range df {
		dd = append(dd, &dm.OtaFirmwareFile{
			FilePath: fileInfo.FilePath,
			Name:     fileInfo.Name,
		})
	}
	return &dm.FirmwareInfo{
		FirmwareID:  di.ID,
		Name:        di.Name,
		Version:     di.Version,
		ProductID:   di.ProductID,
		ProductName: di.ProductID,
		IsDiff:      int32(di.IsDiff),
		CreatedTime: di.CreatedTime.Unix(),
		Desc:        &wrapperspb.StringValue{Value: di.Desc},
		//ExtData:     &wrapperspb.StringValue{Value: di.Extra.String},
		Files:      dd,
		SignMethod: di.SignMethod,
	}
}
func ToFirmwareRespInfo(di *relationDB.DmOtaFirmware, df ...*relationDB.DmOtaFirmwareFile) *dm.FirmwareInfoReadResp {
	dd := make([]*dm.OtaFirmwareFileResp, 0, len(df))
	for _, fileInfo := range df {
		dd = append(dd, &dm.OtaFirmwareFileResp{
			FileID:    fileInfo.ID,
			FilePath:  fileInfo.FilePath,
			Name:      fileInfo.Name,
			Size:      fileInfo.Size,
			Storage:   fileInfo.Storage,
			Host:      fileInfo.Host,
			Signature: fileInfo.Signature,
			//SignMethod: fileInfo.SignMethod,
		})
	}
	return &dm.FirmwareInfoReadResp{
		FirmwareID:  di.ID,
		Name:        di.Name,
		Version:     di.Version,
		ProductID:   di.ProductID,
		ProductName: di.ProductID,
		CreatedTime: di.CreatedTime.Unix(),
		IsDiff:      int32(di.IsDiff),
		Desc:        &wrapperspb.StringValue{Value: di.Desc},
		//ExtData:     &wrapperspb.StringValue{Value: di.Extra.String},
		SignMethod: di.SignMethod,
		Files:      dd,
	}
}
func ToFirmwareFileResp(df *relationDB.DmOtaFirmwareFile) *dm.OtaFirmwareFileInfo {
	return &dm.OtaFirmwareFileInfo{
		FileID:     df.ID,
		FirmwareID: df.FirmwareID,
		FilePath:   df.FilePath,
		Name:       df.Name,
		Size:       df.Size,
		Storage:    df.Storage,
		Host:       df.Host,
		Signature:  df.Signature,
		//SignMethod: df.SignMethod,
	}
}
