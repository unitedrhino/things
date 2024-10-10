package productmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

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
