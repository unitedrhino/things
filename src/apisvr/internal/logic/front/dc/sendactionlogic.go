package dc

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dcsvr/dc"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) SendActionLogic {
	return SendActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendActionLogic) SendAction(req types.SendDcActionReq) (*types.SendDcActionResp, error) {
	l.Infof("SendDcAction|req=%+v", req)
	resp, err := l.svcCtx.DcRpc.SendAction(l.ctx, &dc.SendActionReq{
		MemberID:    req.MemberID,
		MemberType:  req.MemberType,
		ProductID:   req.ProductID,
		DeviceName:  req.DeviceName,
		ActionId:    req.ActionId,
		InputParams: req.InputParams,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.SendAction|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.SendDcActionResp{
		ClientToken:  resp.ClientToken,  //调用id
		OutputParams: resp.OutputParams, //输出参数 注意：此字段可能返回 null，表示取不到有效值。
		Status:       resp.Status,       //返回状态
		Code:         resp.Code,         //设备返回状态码
	}, nil
}
