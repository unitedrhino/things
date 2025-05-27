package productmanagelogic

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/gob"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductInfoMultiExportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductInfoMultiExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductInfoMultiExportLogic {
	return &ProductInfoMultiExportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

type ProductExportInfo struct {
	Info *relationDB.DmProductInfo
	Tls  string
}

func (l *ProductInfoMultiExportLogic) ProductInfoMultiExport(in *dm.ProductInfoExportReq) (*dm.ProductInfoExportResp, error) {
	pos, err := relationDB.NewProductInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.ProductFilter{ProductIDs: in.ProductIDs, WithConfig: true}, nil)
	if err != nil {
		return nil, err
	}
	var infos []ProductExportInfo
	for _, v := range pos {
		model, err := l.svcCtx.ProductSchemaRepo.GetData(l.ctx, v.ProductID)
		if err != nil {
			return nil, err
		}
		infos = append(infos, ProductExportInfo{Tls: model.String(), Info: v})
	}
	ps, err := CompressStruct(infos)
	return &dm.ProductInfoExportResp{Products: ps}, err
}

func CompressStruct(data interface{}) (string, error) {
	var b bytes.Buffer

	// 使用 gzip 包装 buffer
	gz := gzip.NewWriter(&b)

	// 使用 gob 编码器将数据写入 gzip
	enc := gob.NewEncoder(gz)
	if err := enc.Encode(data); err != nil {
		return "", err
	}

	// 关闭 gzip 写入器以刷新缓冲
	if err := gz.Close(); err != nil {
		return "", err
	}

	// 将压缩后的字节流转换为 base64 字符串
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

func DecompressStruct(compressed string, target interface{}) error {
	// 将 Base64 字符串解码为字节流
	compressedData, err := base64.StdEncoding.DecodeString(compressed)
	if err != nil {
		return err
	}

	b := bytes.NewBuffer(compressedData)

	// 创建 gzip 读取器
	gz, err := gzip.NewReader(b)
	if err != nil {
		return err
	}
	defer gz.Close()

	// 使用 gob 解码器从 gzip 读取数据
	dec := gob.NewDecoder(gz)
	return dec.Decode(target)
}
