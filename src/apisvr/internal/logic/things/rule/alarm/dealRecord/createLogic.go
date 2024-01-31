package dealRecord

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.AlarmDealRecordCreateReq) error {
	_, err := l.svcCtx.Alarm.AlarmDealRecordCreate(l.ctx, &rule.AlarmDealRecordCreateReq{
		AlarmRecordID: req.AlarmRecordID,
		Result:        req.Result,
		Type:          1,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.AlarmDealRecordCreate req=%v err=%v", utils.FuncName(), req, er)
		return er
	}

	return nil
}
