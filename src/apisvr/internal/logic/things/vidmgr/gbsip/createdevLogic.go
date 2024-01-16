package gbsip

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatedevLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatedevLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatedevLogic {
	return &CreatedevLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatedevLogic) Createdev(req *types.VidmgrSipCreateDevReq) error {
	// todo: add your logic here and delete this line
	vidReq := &vid.VidmgrGbsipDeviceCreateReq{
		DeviceID: req.DeviceID,
		Name:     req.Name,
		PWD:      req.PWD,
	}
	jsonStr, _ := json.Marshal(req)
	fmt.Println("airgens create dev:", string(jsonStr))
	_, err := l.svcCtx.VidmgrG.VidmgrGbsipDeviceCreate(l.ctx, vidReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ManageVidmgr req=%v err=%v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
