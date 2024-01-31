package gbsip

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/vidsip/pb/sip"
	"github.com/i-Things/things/service/vidsvr/pb/vid"

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
	jsonStr, _ := json.Marshal(req)
	fmt.Println("airgens create dev:", string(jsonStr))
	//查询流服务ID
	vidResp, err := l.svcCtx.VidmgrM.VidmgrInfoRead(l.ctx, &vid.VidmgrInfoReadReq{
		VidmgrID: req.VidmgrID,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s rpc.ManageVidmgr req=%v err=%+v", utils.FuncName(), req, er)
		return errors.MediaSipDevCreateError.AddDetail("流服务不存在！")
	}

	sipReq := &sip.SipDevCreateReq{
		DeviceID:  req.DeviceID,
		VidmgrID:  req.VidmgrID,
		Name:      req.Name,
		PWD:       req.PWD,
		MediaPort: 10000,
		MediaIP:   vidResp.VidmgrIpV4,
	}
	_, err = l.svcCtx.SipRpc.SipDeviceCreate(l.ctx, sipReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ManageVidmgr req=%v err=%v", utils.FuncName(), req, er)
		return errors.MediaSipDevCreateError.AddDetail("RPC服务调用失败！")
	}
	return nil
}
