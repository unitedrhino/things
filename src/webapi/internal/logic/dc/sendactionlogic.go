package dc

import (
	"context"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/dcsvr/dc"

	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
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

func (l *SendActionLogic) SendAction(req types.SendActionReq) (*types.SendActionResp, error) {
	l.Infof("SendAction|req=%+v", req)
	resp,err := l.svcCtx.DcRpc.SendAction(l.ctx,&dc.SendActionReq{
		MemberID:             req.MemberID,
		MemberType:           req.MemberType,
		ProductID:            req.ProductID,
		DeviceName:           req.DeviceName,
		ActionId:             req.ActionId,
		InputParams:          req.InputParams,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.SendAction|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.SendActionResp{
		ClientToken:resp.ClientToken,  //调用id
		OutputParams:resp.OutputParams, //输出参数 注意：此字段可能返回 null，表示取不到有效值。
		Status:resp.Status,       //返回状态
		Code:resp.Code,  //设备返回状态码
	}, nil
}
