package productmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

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
				LoginAuthScript: nil,
				TransformScript: nil,
				CustomTopics:    nil,
			}, nil
		}
		return nil, err
	}
	return &dm.ProductCustom{
		ProductID:       pi.ProductID,
		ScriptLang:      pi.ScriptLang,
		TransformScript: utils.ToRpcNullString(pi.TransformScript),
		LoginAuthScript: utils.ToRpcNullString(pi.LoginAuthScript),
		CustomTopics:    logic.ToCustomTopicsPb(pi.CustomTopics),
	}, nil
}
