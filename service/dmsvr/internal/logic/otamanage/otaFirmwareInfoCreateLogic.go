package otamanagelogic

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/domain/deviceMsg/msgOta"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/oss"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaFirmwareInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProductInfoRepo
	OfDB *relationDB.OtaFirmwareInfoRepo
}

func NewOtaFirmwareInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaFirmwareInfoCreateLogic {
	return &OtaFirmwareInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
		OfDB:   relationDB.NewOtaFirmwareInfoRepo(ctx),
	}
}

func (l *OtaFirmwareInfoCreateLogic) CheckOtaFirmwareInfo(in *dm.OtaFirmwareInfoCreateReq) (bool, error) {
	//查询升级包是否存在
	logx.Infof("relationDB:%+v", relationDB.ProductFilter{ProductIDs: []string{in.ProductID}})
	_, err := l.svcCtx.ProductCache.GetData(l.ctx, in.ProductID)
	if errors.Cmp(err, errors.NotFind) {
		l.Errorf("not find product id:" + in.ProductID)
		return false, nil
	} else if err != nil {
		return false, nil
	}
	fsize, err := l.OfDB.CountByFilter(l.ctx, relationDB.OtaFirmwareInfoFilter{
		Version:   in.Version,
		ProductID: in.ProductID,
	})
	if fsize == 0 {
		return true, nil
	} else if err != nil {
		return true, err
	}
	if in.ModuleCode != "" && in.ModuleCode != msgOta.ModuleCodeDefault {
		module, err := relationDB.NewOtaModuleInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.OtaModuleInfoFilter{Code: in.ModuleCode})
		if err != nil {
			return false, err
		}
		if module.ProductID != in.ProductID {
			return false, errors.Parameter.AddMsg("选择的模块产品和升级包的产品不一致")
		}
	} else {
		in.ModuleCode = msgOta.ModuleCodeDefault
	}

	return true, errors.Parameter.WithMsg(fmt.Sprintf("%s版本的升级包已经存在", in.Version))
}

// 添加升级包
func (l *OtaFirmwareInfoCreateLogic) OtaFirmwareInfoCreate(in *dm.OtaFirmwareInfoCreateReq) (*dm.WithID, error) {
	//todo debug
	//if err := ctxs.IsRoot(l.ctx); err != nil {
	//	return nil, err
	//}
	l.ctx = ctxs.WithRoot(l.ctx)
	//校验版本号是否重复
	logx.Infof("ctx:%+v", l.ctx)
	var fileDB = relationDB.NewOtaFirmwareFileRepo(l.ctx)
	logx.Infof("%+v", in)
	find, err := l.CheckOtaFirmwareInfo(in)
	if errors.Cmp(err, errors.Parameter) {
		return nil, err
	} else if err != nil {
		l.Errorf("AddDevice|CheckProduct|in=%v\n", in)
		return nil, err
	} else if find == false {
		return nil, err
	}
	var totalSize int64
	var files []*relationDB.DmOtaFirmwareFile
	if len(in.FilePaths) > 0 {
		for _, filePath := range in.FilePaths {
			nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessOta, oss.SceneFirmware, fmt.Sprintf("%s/%s/%s", in.ProductID, in.Version, oss.GetFileNameWithPath(filePath)))
			path, err := l.svcCtx.OssClient.PublicBucket().CopyFromTempBucket(filePath, nwePath)
			if err != nil {
				return nil, errors.System.AddDetail(err)
			}

			OtaFirmwareInfoFileInfo, err := l.svcCtx.OssClient.PublicBucket().GetObjectInfo(l.ctx, path)
			if err != nil {
				logx.Error(err)
				return nil, err
			}
			totalSize += OtaFirmwareInfoFileInfo.Size
			files = append(files, &relationDB.DmOtaFirmwareFile{
				Name:      oss.GetFileNameWithPath(filePath),
				FilePath:  path,
				Size:      OtaFirmwareInfoFileInfo.Size,
				Signature: OtaFirmwareInfoFileInfo.Md5,
				FileMd5:   OtaFirmwareInfoFileInfo.Md5,
			})
		}
	}
	di := relationDB.DmOtaFirmwareInfo{
		ProductID:      in.ProductID,
		Version:        in.Version,
		Name:           in.Name,
		Desc:           in.Desc,
		SignMethod:     in.SignMethod,
		SrcVersion:     in.SrcVersion,
		IsDiff:         in.IsDiff,
		Extra:          in.Extra.GetValue(),
		TotalSize:      totalSize,
		IsNeedToVerify: in.IsNeedToVerify,
		ModuleCode:     in.ModuleCode,
	}
	//是否需要平台验证
	di.Status = msgOta.OtaFirmwareStatusNotVerified
	if in.IsNeedToVerify != def.True {
		di.Status = msgOta.OtaFirmwareStatusNotRequired
	}
	//整包或差包
	if in.IsDiff == msgOta.DiffPackage {
		di.SrcVersion = in.SrcVersion
	}
	err = l.OfDB.Insert(l.ctx, &di)
	if err != nil {
		l.Errorf("AddFirmware.FirmwareInfo.Insert err=%+v", err)
		return nil, err
	}
	//插入file
	firmWareId := di.ID
	for i := range files {
		files[i].FirmwareID = di.ID
	}
	err = fileDB.MultiInsert(l.ctx, files)
	if err != nil {
		l.Errorf("fileDB.MultiInsert err=%+v", err)
		return nil, err
	}
	return &dm.WithID{Id: firmWareId}, nil

}
