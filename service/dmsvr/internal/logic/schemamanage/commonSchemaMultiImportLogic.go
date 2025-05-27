package schemamanagelogic

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/share/errors"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommonSchemaMultiImportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommonSchemaMultiImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommonSchemaMultiImportLogic {
	return &CommonSchemaMultiImportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CommonSchemaMultiImportLogic) CommonSchemaMultiImport(in *dm.CommonSchemaImportReq) (*dm.ImportResp, error) {
	var schemas []*dm.CommonSchemaInfo
	err := json.Unmarshal([]byte(in.Schemas), &schemas)
	if err != nil {
		return nil, err
	}
	var resp = dm.ImportResp{Total: int64(len(schemas))}
	h := NewCommonSchemaCreateLogic(l.ctx, l.svcCtx)
	for _, v := range schemas {
		_, err := h.CommonSchemaCreate(&dm.CommonSchemaCreateReq{Info: v})
		if err != nil {
			if errors.Cmp(err, errors.Duplicate) {
				resp.IgnoreCount++
				continue
			}
			l.Error(v, err)
			resp.ErrCount++
			continue
		}
		resp.SuccCount++
	}
	return &resp, nil
}
