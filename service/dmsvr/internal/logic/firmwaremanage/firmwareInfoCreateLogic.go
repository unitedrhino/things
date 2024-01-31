package firmwaremanagelogic

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/shared/utils"
	"path"

	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/oss"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type FirmwareInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProductInfoRepo
	OfDB *relationDB.OtaFirmwareRepo
}

func NewFirmwareInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FirmwareInfoCreateLogic {
	return &FirmwareInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
		OfDB:   relationDB.NewOtaFirmwareRepo(ctx),
	}
}

func (l *FirmwareInfoCreateLogic) CheckFirmware(in *dm.Firmware) (bool, error) {
	_, err := l.PiDB.FindOneByFilter(l.ctx, relationDB.ProductFilter{ProductIDs: []string{in.ProductID}})
	if errors.Cmp(err, errors.NotFind) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	fsize, err := l.OfDB.CountByFilter(l.ctx, relationDB.OtaFirmwareFilter{
		Version: in.Version,
	})
	if fsize == 0 {
		return true, nil
	} else if err != nil {
		return true, err
	}
	return true, errors.Parameter.WithMsg(fmt.Sprintf("%s版本的升级包已存在", in.Version))
}

// 新增固件升级包
func (l *FirmwareInfoCreateLogic) FirmwareInfoCreate(in *dm.Firmware) (*dm.FirmwareResp, error) {
	var fileDB = relationDB.NewOtaFirmwareFileRepo(l.ctx)
	find, err := l.CheckFirmware(in)
	if errors.Cmp(err, errors.Parameter) {
		return nil, err
	} else if err != nil {
		l.Errorf("AddDevice|CheckProduct|in=%v\n", in)
		return nil, errors.Database.AddDetail(err)
	} else if find == false {
		return nil, errors.Parameter.AddDetail("not find product id:" + utils.ToString(in.ProductID))
	}
	if len(in.Files) > 0 {
		for k, file := range in.Files {
			si, err := oss.GetSceneInfo(file.FilePath)
			if err != nil {
				return nil, err
			}
			if !(si.Business == oss.BusinessProductManage && si.Scene == oss.SceneOta) {
				return nil, errors.Parameter.WithMsg("附件的路径不对")
			}
			si.FilePath = in.ProductID + path.Ext(si.FilePath)
			logx.Info("file_path:", si.FilePath)
			nwePath, err := oss.GetFilePath(si, false)
			if err != nil {
				return nil, err
			}
			path, err := l.svcCtx.OssClient.PrivateBucket().CopyFromTempBucket(file.FilePath, nwePath)
			if err != nil {
				return nil, errors.System.AddDetail(err)
			}
			in.Files[k].FilePath = path
		}
	}

	di := relationDB.DmOtaFirmware{
		ProductID: in.ProductID, // 产品id
		IsDiff:    int64(in.IsDiff),
		Version:   in.Version,
		Name:      in.Name,
		Desc:      in.Desc.Value,
	}
	err = l.OfDB.Insert(l.ctx, &di)
	if err != nil {
		l.Errorf("AddFirmware.FirmwareInfo.Insert err=%+v", err)
		return nil, errors.System.AddDetail(err)
	}
	//插入file
	firmwareID := di.ID
	for _, file := range in.Files {
		ff := relationDB.DmOtaFirmwareFile{
			FirmwareID: firmwareID,
			//SignMethod:  in.SignMethod,
			FilePath: file.FilePath,
			Name:     file.Name,
		}
		fileDB.Insert(l.ctx, &ff)
	}
	return &dm.FirmwareResp{FirmwareID: firmwareID}, nil
}
