package productmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductScriptReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductScriptReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductScriptReadLogic {
	return &ProductScriptReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 脚本管理
func (l *ProductScriptReadLogic) ProductScriptRead(in *dm.ProductScriptReadReq) (*dm.ProductScript, error) {
	pi, err := l.svcCtx.ProductScript.FindOneByProductID(l.ctx, in.ProductID)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.NotFind
		}
		return nil, err
	}
	return &dm.ProductScript{
		ProductID: pi.ProductID,
		Script:    pi.Script,
		Lang:      pi.Lang,
	}, nil
}
