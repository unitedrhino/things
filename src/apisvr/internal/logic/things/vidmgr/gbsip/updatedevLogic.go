package gbsip

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/vidsip/pb/sip"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatedevLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatedevLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatedevLogic {
	return &UpdatedevLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatedevLogic) Updatedev(req *types.VidmgrSipUpdateDevReq) error {
	// todo: add your logic here and delete this line
	vidReq := &sip.SipDevUpdateReq{
		DeviceID: req.DeviceID,
		Name:     req.Name,
		PWD:      req.PWD,
	}
	jsonStr, _ := json.Marshal(req)
	fmt.Println("airgens Updatedev:", string(jsonStr))
	_, err := l.svcCtx.SipRpc.SipDeviceUpdate(l.ctx, vidReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ManageVidmgr req=%v err=%v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
