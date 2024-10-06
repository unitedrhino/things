package productmanagelogic

import (
	"archive/zip"
	"context"
	"fmt"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/oss"
	"gitee.com/i-Things/share/oss/common"
	"gitee.com/i-Things/share/systems"
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/i-Things/things/service/dmsvr/internal/svc"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"strings"
)

type ProductInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProductInfoRepo
}

func NewProductInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductInfoUpdateLogic {
	return &ProductInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
	}
}

//
//// archiver解压压缩包
//func ArchiverTest(path string) {
//	f, _ := os.Open(path)
//	format, readStream, err := archiver.Identify(path, f)
//	if err != nil {
//		return
//	}
//	extractor, ok := format.(archiver.Extractor)
//	if !ok {
//		return
//	}
//	switch extractor.(type) {
//	case archiver.Zip:
//		extractor = archiver.Zip{TextEncoding: "gbk"}
//		fmt.Println("archiver.Zip")
//	case archiver.SevenZip:
//		extractor = archiver.SevenZip{}
//		fmt.Println("archiver.SevenZip")
//	case archiver.Rar:
//		extractor = archiver.Rar{}
//		fmt.Println("archiver.Rar")
//	default:
//		fmt.Println("unsupported compression algorithm")
//		return
//	}
//
//	//fileList := []string{"file1.txt", "subfolder"}
//	ctx := context.Background()
//	handler := func(ctx context.Context, f archiver.File) error {
//		filename := f.Name()
//		newfile, err := os.Create(filename)
//		if err != nil {
//			panic(err)
//		}
//		defer newfile.Close()
//		old, err := f.Open()
//		if err != nil {
//			panic(err)
//		}
//		defer old.Close()
//		_, err = io.Copy(newfile, old)
//
//		if err != nil {
//			panic(err)
//		}
//		fmt.Printf("extracted %s \n", f.Name())
//		return nil
//	}
//
//	err = extractor.Extract(ctx, readStream, nil, handler)
//	if err != nil {
//		return
//	}
//
//}

func (l *ProductInfoUpdateLogic) setPoByPb(old *relationDB.DmProductInfo, data *dm.ProductInfo) error {
	if data.Tags != nil {
		old.Tags = data.Tags
	}
	if data.ProductName != "" {
		old.ProductName = data.ProductName
	}
	if data.TrialTime != nil {
		old.TrialTime = data.TrialTime.GetValue()
	}
	if data.SceneMode != "" {
		old.SceneMode = data.SceneMode
	}
	if data.AuthMode != def.Unknown {
		old.AuthMode = data.AuthMode
	}
	if data.Status != 0 {
		old.Status = data.Status
	}
	if data.Desc != nil {
		old.Desc = data.Desc.GetValue()
	}

	if data.AutoRegister != def.Unknown {
		old.AutoRegister = data.AutoRegister
	}

	if data.ProductName != "" {
		old.ProductName = data.ProductName
	}
	if data.ProductImg != "" && data.IsUpdateProductImg == true { //如果填了参数且不等于原来的,说明修改头像,需要处理
		if old.ProductImg != "" {
			err := l.svcCtx.OssClient.PublicBucket().Delete(l.ctx, old.ProductImg, common.OptionKv{})
			if err != nil {
				l.Errorf("Delete file err path:%v,err:%v", old.ProductImg, err)
			}
		}
		nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessProductManage, oss.SceneProductImg, fmt.Sprintf("%s/%s", data.ProductID, oss.GetFileNameWithPath(data.ProductImg)))
		path, err := l.svcCtx.OssClient.PublicBucket().CopyFromTempBucket(data.ProductImg, nwePath)
		if err != nil {
			return errors.System.AddDetail(err)
		}

		old.ProductImg = path
	}
	if data.CustomUi != nil {
		for k, v := range data.CustomUi {
			if v.IsUpdateUi == false || v.Path == "" {
				if old.CustomUi != nil && old.CustomUi[k] != nil {
					v.Version = old.CustomUi[k].Version
				}
				continue
			}
			fileName := oss.GetFileNameWithPath(v.Path)
			localFilePath := fmt.Sprintf("%s/things/product/customUi/%s", systems.TmpDir, fileName)
			defer os.RemoveAll(fileName)
			err := l.svcCtx.OssClient.TemporaryBucket().GetObjectLocal(l.ctx, v.Path, localFilePath)
			if err != nil {
				return errors.System.AddMsg("拉取文件出错").AddDetail(err)
			}
			archive, err := zip.OpenReader(localFilePath)
			if err != nil {
				return errors.System.AddDetail(err)
			}
			defer archive.Close()
			var version int64 = 1
			if old.CustomUi != nil && old.CustomUi[k] != nil {
				version += old.CustomUi[k].Version
			}
			v.Version = version
			uploadPath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessProductManage, oss.SceneProductCustomUi, fmt.Sprintf("%s/%d", data.ProductID, version))
			for _, f := range archive.File {
				if f.FileInfo().IsDir() {
					continue
				}
				path := f.Name
				if strings.HasPrefix(path, "dist/") {
					path = strings.TrimPrefix(path, "dist/")
				}
				uPath := fmt.Sprintf("%s/%s", uploadPath, path)
				r, err := f.Open()
				if err != nil {
					l.Error(err)
					return errors.System.AddDetail(err)
				}
				_, err = l.svcCtx.OssClient.PublicBucket().Upload(l.ctx, uPath, r, common.OptionKv{})
				if err != nil {
					l.Error(err)
					return errors.System.AddDetail(err)
				}
			}
			v.Path = fmt.Sprintf("%s/index.html", uploadPath)
		}
		old.CustomUi = utils.CopyMap[relationDB.ProductCustomUi](data.CustomUi)
		if old.CustomUi == nil {
			old.CustomUi = map[string]*relationDB.ProductCustomUi{}
		}
	}
	if data.AuthMode != 0 {
		old.AuthMode = data.AuthMode
	}
	if len(data.ProtocolConf) != 0 {
		old.ProtocolConf = data.ProtocolConf
	}
	if data.CategoryID != 0 {
		old.CategoryID = data.CategoryID
	}
	if data.NetType != 0 {
		old.NetType = data.NetType
	}
	if data.ProtocolCode != "" {
		old.ProtocolCode = data.ProtocolCode
	}
	if data.AutoRegister != 0 {
		old.AutoRegister = data.AutoRegister
	}
	return nil
}

// 更新设备
func (l *ProductInfoUpdateLogic) ProductInfoUpdate(in *dm.ProductInfo) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	po, err := l.PiDB.FindOneByFilter(l.ctx, relationDB.ProductFilter{ProductIDs: []string{in.ProductID}})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddDetail("not find Product_id id:" + cast.ToString(in.ProductID))
		}
		return nil, err
	}

	err = l.setPoByPb(po, in)
	if err != nil {
		return nil, err
	}
	err = l.PiDB.Update(l.ctx, po)
	if err != nil {
		l.Errorf("%s.Update err=%+v", utils.FuncName(), err)
		if errors.Cmp(err, errors.Duplicate) {
			return nil, errors.Duplicate.WithMsgf("产品名称重复:%s", in.ProductName)
		}
		return nil, err
	}
	err = l.svcCtx.ProductCache.SetData(l.ctx, in.ProductID, nil)
	if err != nil {
		l.Error(err)
	}
	return &dm.Empty{}, nil
}
