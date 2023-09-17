package productmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductCustomReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PcDB *relationDB.ProductCustomRepo
}

func NewProductCustomReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductCustomReadLogic {
	return &ProductCustomReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PcDB:   relationDB.NewProductCustomRepo(ctx),
	}
}

// 脚本管理
func (l *ProductCustomReadLogic) ProductCustomRead(in *dm.ProductCustomReadReq) (*dm.ProductCustom, error) {
	pi, err := l.PcDB.FindOneByProductID(l.ctx, in.ProductID)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return &dm.ProductCustom{
				ProductID:       in.ProductID,
				ScriptLang:      1,
				TransformScript: nil,
				CustomTopic:     nil,
			}, nil
		}
		return nil, err
	}
	var customTopic []string
	utils.SqlNullStringToAny(pi.CustomTopic, &customTopic)
	return &dm.ProductCustom{
		ProductID:       pi.ProductID,
		ScriptLang:      pi.ScriptLang,
		TransformScript: utils.ToRpcNullString(pi.TransformScript),
		CustomTopic:     customTopic,
	}, nil
}
