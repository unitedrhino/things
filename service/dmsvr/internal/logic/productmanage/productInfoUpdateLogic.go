package productmanagelogic

import (
	"archive/zip"
	"context"
	"fmt"
	"os"
	"strings"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/oss"
	"gitee.com/unitedrhino/share/oss/common"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/systems"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"gitee.com/unitedrhino/things/share/topics"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
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

func (l *ProductInfoUpdateLogic) setPoByPb(old *relationDB.DmProductInfo, data *dm.ProductInfo) (funcs []func(tx *stores.DB) error, err error) {
	if data.Tags != nil {
		old.Tags = data.Tags
	}
	if data.DeviceType != 0 {
		old.DeviceType = data.DeviceType
	}
	if data.BindLevel != 0 {
		old.BindLevel = data.BindLevel
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
	if data.SubProtocolCode != nil {
		old.SubProtocolCode = data.SubProtocolCode.GetValue()
	}

	if data.AutoRegister != def.Unknown {
		old.AutoRegister = data.AutoRegister
	}
	if data.OnlineHandle != def.Unknown {
		old.OnlineHandle = data.OnlineHandle
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
			return nil, errors.System.AddDetail(err)
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
				return nil, errors.System.AddMsg("拉取文件出错").AddDetail(err)
			}
			archive, err := zip.OpenReader(localFilePath)
			if err != nil {
				return nil, errors.System.AddDetail(err)
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
					return nil, errors.System.AddDetail(err)
				}
				_, err = l.svcCtx.OssClient.PublicBucket().Upload(l.ctx, uPath, r, common.OptionKv{})
				if err != nil {
					l.Error(err)
					return nil, errors.System.AddDetail(err)
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
	if data.ProtocolConf != nil {
		old.ProtocolConf = data.ProtocolConf
	}
	if data.SubProtocolConf != nil {
		old.SubProtocolConf = data.SubProtocolConf
	}
	if data.CategoryID != 0 && data.CategoryID != old.CategoryID {
		var schemas []*relationDB.DmProductSchema
		if data.CategoryID != def.NotClassified {
			var categoryIDs = []int64{def.RootNode}
			if data.CategoryID != def.RootNode {
				pcs, err := relationDB.NewProductCategoryRepo(l.ctx).FindOne(l.ctx, data.CategoryID)
				if err != nil {
					return nil, err
				}
				if pcs.IDPath != "" {
					categoryIDs = append(categoryIDs, utils.GetIDPath(pcs.IDPath)...)
				}
			}
			pcss, err := relationDB.NewCommonSchemaRepo(l.ctx).FindByFilter(l.ctx, relationDB.CommonSchemaFilter{
				ProductCategoryIDs: categoryIDs,
			}, nil)
			if err != nil {
				return nil, err
			}
			for _, pcs := range pcss {
				pcs.Tag = schema.TagRequired
				schemas = append(schemas, &relationDB.DmProductSchema{
					TenantCode:   old.TenantCode,
					ProductID:    old.ProductID,
					Identifier:   pcs.Identifier,
					DmSchemaCore: pcs.DmSchemaCore,
				})
			}
		}
		funcs = append(funcs, func(tx *stores.DB) error {
			err = relationDB.NewProductSchemaRepo(tx).UpdateTag(l.ctx, []string{data.ProductID}, nil, schema.TagRequired, schema.TagOptional)
			if err != nil {
				return err
			}
			if len(schemas) > 0 {
				err = relationDB.NewProductSchemaRepo(tx).MultiInsert(l.ctx, schemas)
				if err != nil {
					return err
				}
			}
			return nil
		})
		old.CategoryID = data.CategoryID
	}
	if data.NetType != 0 {
		old.NetType = data.NetType
	}
	if data.ProtocolCode != "" {
		old.ProtocolCode = data.ProtocolCode
	}
	if data.AutoRegister != def.Unknown {
		old.AutoRegister = data.AutoRegister
	}
	if data.OnlineHandle != def.Unknown {
		old.OnlineHandle = data.OnlineHandle
	}
	if data.DeviceSchemaMode != 0 {
		old.DeviceSchemaMode = data.DeviceSchemaMode
	}
	return
}

// 更新设备
func (l *ProductInfoUpdateLogic) ProductInfoUpdate(in *dm.ProductInfo) (*dm.Empty, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	po, err := l.PiDB.FindOneByFilter(l.ctx, relationDB.ProductFilter{ProductIDs: []string{in.ProductID}})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddDetail("not find Product_id id:" + cast.ToString(in.ProductID))
		}
		return nil, err
	}
	txFuncs, err := l.setPoByPb(po, in)
	if err != nil {
		return nil, err
	}
	if len(txFuncs) == 0 {
		err = l.PiDB.Update(l.ctx, po)
	} else {
		stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
			for _, txFunc := range txFuncs {
				err := txFunc(tx)
				if err != nil {
					return err
				}
			}
			return relationDB.NewProductInfoRepo(tx).Update(l.ctx, po)
		})
	}
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
	err = l.svcCtx.FastEvent.Publish(l.ctx, topics.DmProductInfoUpdate, in.ProductID)
	if err != nil {
		l.Error(err)
	}
	return &dm.Empty{}, nil
}
