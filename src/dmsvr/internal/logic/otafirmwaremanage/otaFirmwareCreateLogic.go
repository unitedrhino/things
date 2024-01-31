package otafirmwaremanagelogic

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/oss"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg/msgOta"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"path"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaFirmwareCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProductInfoRepo
	OfDB *relationDB.OtaFirmwareRepo
}

func NewOtaFirmwareCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaFirmwareCreateLogic {
	return &OtaFirmwareCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
		OfDB:   relationDB.NewOtaFirmwareRepo(ctx),
	}
}

func (l *OtaFirmwareCreateLogic) CheckOtaFirmware(in *dm.OtaFirmwareCreateReq) (bool, error) {
	//查询升级包是否存在
	logx.Infof("relationDB:%+v", relationDB.ProductFilter{ProductIDs: []string{in.ProductID}})
	_, err := l.PiDB.FindOneByFilter(l.ctx, relationDB.ProductFilter{ProductIDs: []string{in.ProductID}})
	if errors.Cmp(err, errors.NotFind) {
		l.Errorf("not find product id:" + in.ProductID)
		return false, nil
	} else if err != nil {
		return false, nil
	}
	fsize, err := l.OfDB.CountByFilter(l.ctx, relationDB.OtaFirmwareFilter{
		Version: in.DestVersion,
	})
	if fsize == 0 {
		return true, nil
	} else if err != nil {
		return true, err
	}
	return true, errors.Parameter.WithMsg(fmt.Sprintf("%s版本的升级包已经存在", in.DestVersion))
}

// 添加升级包
func (l *OtaFirmwareCreateLogic) OtaFirmwareCreate(in *dm.OtaFirmwareCreateReq) (*dm.OtaFirmwareResp, error) {
	//校验版本号是否重复
	logx.Infof("ctx:%+v", l.ctx)
	var fileDB = relationDB.NewOtaFirmwareFileRepo(l.ctx)
	logx.Infof("%+v", in)
	find, err := l.CheckOtaFirmware(in)
	if errors.Cmp(err, errors.Parameter) {
		return nil, err
	} else if err != nil {
		l.Errorf("AddDevice|CheckProduct|in=%v\n", in)
		return nil, errors.Database.AddDetail(err)
	} else if find == false {
		return nil, err
	}
	var total_size int64
	if len(in.FirmwareFiles) > 0 {
		for k, firmwareFile := range in.FirmwareFiles {
			//阿里云寻找scene
			logx.Infof("k:%+v", k)
			logx.Infof("firmwareFile:%+v", firmwareFile)
			logx.Infof("url:%+v", firmwareFile.FilePath)
			info, err := oss.GetSceneInfo(firmwareFile.FilePath)
			if err != nil {
				return nil, err
			}
			if !(info.Business == oss.BusinessProductManage && info.Scene == oss.SceneOta) {
				return nil, errors.Parameter.WithMsg("附件的路径不对")
			}
			info.FilePath = in.ProductID + path.Ext(info.FilePath)
			logx.Info("file_path:", info.FilePath)
			newPath, err := oss.GetFilePath(info, false)
			if err != nil {
				return nil, err
			}
			logx.Infof("newPath:%+v", newPath)
			path, err := l.svcCtx.OssClient.PrivateBucket().CopyFromTempBucket(firmwareFile.FilePath, newPath)
			if err != nil {
				logx.Error(err)
				return nil, errors.System.AddDetail(err)
			}
			otaFirmwareFileInfo, err := l.svcCtx.OssClient.PrivateBucket().GetObjectInfo(l.ctx, path)
			if err != nil {
				logx.Error(err)
				return nil, errors.System.AddDetail(err)
			}
			total_size += otaFirmwareFileInfo.Size
			in.FirmwareFiles[k].FilePath = path
			in.FirmwareFiles[k].Size = otaFirmwareFileInfo.Size
			//in.FirmwareFiles[k].FileMd5 = otaFirmwareFileInfo.Md5
			in.FirmwareFiles[k].Signature = otaFirmwareFileInfo.Md5
		}
	}
	logx.Infof("req:%+v", in)
	logx.Infof("extra:%+v", in.FirmwareUdi)
	logx.Infof("extra:%+v", in.FirmwareUdi.Value)
	di := relationDB.DmOtaFirmware{
		ProductID:  in.ProductID,
		Version:    in.DestVersion,
		Name:       in.FirmwareName,
		Desc:       in.FirmwareDesc,
		SignMethod: in.SignMethod,
		SrcVersion: in.SrcVersion,
		Module:     in.Module,
		IsDiff:     in.IsDiff,
		Extra:      in.FirmwareUdi.Value,
		TotalSize:  total_size,
	}
	//是否需要平台验证
	di.Status = msgOta.OtaFirmwareStatusNotVerified
	if !in.NeedToVerify {
		di.Status = msgOta.OtaFirmwareStatusNotRequired
	}
	//整包或差包
	if in.IsDiff == msgOta.DiffPackage {
		di.SrcVersion = in.SrcVersion
	}
	l.OfDB.Insert(l.ctx, &di)
	if err != nil {
		l.Errorf("AddFirmware.FirmwareInfo.Insert err=%+v", err)
		return nil, errors.System.AddDetail(err)
	}
	//插入file
	firmWareId := di.ID
	for _, firmWareFile := range in.FirmwareFiles {
		ff := relationDB.DmOtaFirmwareFile{
			FirmwareID: firmWareId,
			FilePath:   firmWareFile.FilePath,
			Name:       firmWareFile.Name,
			Size:       firmWareFile.Size,
			Storage:    "minio",
			Signature:  firmWareFile.Signature,
		}
		fileDB.Insert(l.ctx, &ff)
	}
	return &dm.OtaFirmwareResp{FirmwareID: firmWareId}, nil

}
