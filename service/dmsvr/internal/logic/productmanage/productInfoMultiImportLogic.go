package productmanagelogic

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductInfoMultiImportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductInfoMultiImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductInfoMultiImportLogic {
	return &ProductInfoMultiImportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductInfoMultiImportLogic) ProductInfoMultiImport(in *dm.ProductInfoImportReq) (*dm.ImportResp, error) {
	var infos []ProductExportInfo
	err := json.Unmarshal([]byte(in.Products), &infos)
	if err != nil {
		return nil, err
	}
	var resp = dm.ImportResp{Total: int64(len(infos))}
	for _, info := range infos {
		pi := logic.ToProductInfo(l.ctx, l.svcCtx, info.Info)
		_, err := NewProductInfoCreateLogic(l.ctx, l.svcCtx).ProductInfoCreate(pi)
		if err != nil {
			if errors.Cmp(err, errors.Duplicate) {
				resp.IgnoreCount++
				continue
			}
			l.Error(pi, err)
			resp.ErrCount++
			continue
		}
		resp.SuccCount++
		if info.Info.Config != nil {
			err = relationDB.NewProductConfigRepo(l.ctx).UpdateWithProducID(l.ctx, info.Info.Config)
			if err != nil {
				l.Error(pi, err)
			}
		}
		if info.Tls != "" {
			_, err = NewProductSchemaTslImportLogic(l.ctx, l.svcCtx).ProductSchemaTslImport(&dm.ProductSchemaTslImportReq{ProductID: info.Info.ProductID, Tsl: info.Tls})
			if err != nil {
				l.Error(pi, err)
			}
		}
	}
	return &resp, nil
}
