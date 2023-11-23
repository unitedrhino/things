package otafirmwaremanagelogic

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/oss"
	"github.com/i-Things/things/shared/utils/cast"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"path"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOTAFirmwareLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProductInfoRepo
	OfDB *relationDB.OtaFirmwareRepo
}

func NewCreateOTAFirmwareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOTAFirmwareLogic {
	return &CreateOTAFirmwareLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
		OfDB:   relationDB.NewOtaFirmwareRepo(ctx),
	}
}

func (l *CreateOTAFirmwareLogic) CheckOtaFirmware(in *dm.OtaFirmwareReq) (bool, error) {
	//查询升级包是否存在
	logx.Infof("in:%+v", in)
	logx.Infof("productId:%+v", in.ProductID)
	logx.Infof("lctx:%+v", l.ctx)
	logx.Infof("relationDB:%+v", relationDB.ProductFilter{ProductIDs: []string{in.ProductID}})
	logx.Infof("PiDB:%+v", l.PiDB)
	_, err := l.PiDB.FindOneByFilter(l.ctx, relationDB.ProductFilter{ProductIDs: []string{in.ProductID}})
	if errors.Cmp(err, errors.NotFind) {
		l.Errorf("not find product id:" + cast.ToString(in.ProductID))
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
func (l *CreateOTAFirmwareLogic) CreateOTAFirmware(in *dm.OtaFirmwareReq) (*dm.OtaFirmwareResp, error) {
	// todo: add your logic here and delete this line
	//校验版本号是否重复
	logx.Infof("ctx:%+v", l.ctx)
	var fileDB = relationDB.NewOtaFirmwareFileRepo(l.ctx)
	logx.Infof("%+v", &in)
	find, err := l.CheckOtaFirmware(in)
	if errors.Cmp(err, errors.Parameter) {
		return nil, err
	} else if err != nil {
		l.Errorf("AddDevice|CheckProduct|in=%v\n", in)
		return nil, errors.Database.AddDetail(err)
	} else if find == false {
		return nil, err
	}
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
			in.FirmwareFiles[k].FilePath = path
		}
	}
	di := relationDB.DmOtaFirmware{
		ProductID:  in.ProductID,
		Version:    in.DestVersion,
		Name:       in.FirmwareName,
		Desc:       in.FirmwareDesc,
		SignMethod: in.SignMethod,
		Module:     in.Module,
		IsDiff:     in.IsDiff,
	}
	di.Status = 0
	if !in.NeedToVerify {
		di.Status = -1
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
		}
		fileDB.Insert(l.ctx, &ff)
	}
	return &dm.OtaFirmwareResp{FirmwareID: firmWareId}, nil
}
