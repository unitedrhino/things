package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/pb/vid"
	"time"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OnServerKeepaliveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOnServerKeepaliveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnServerKeepaliveLogic {
	return &OnServerKeepaliveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnServerKeepaliveLogic) OnServerKeepalive(req *types.HooksApiServerKeepaliveReq) (resp *types.HooksApiResp, err error) {
	// todo: add your logic here and delete this line
	reqStr, _ := json.Marshal(*req)
	fmt.Println("---------OnServerKeepalive--------------:", string(reqStr))
	//从hook data中解析出 MediaserverID值。
	//hookactive中是保持在线状态  需要更新对应的数据库
	vidmgrInfo, err := l.svcCtx.VidmgrM.VidmgrInfoRead(l.ctx, &vid.VidmgrInfoReadReq{
		VidmgrtID: req.MediaServerId,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s rpc.ManageVidmgr req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	if vidmgrInfo != nil {
		//update info
		//UPDATE
		vidReq := &vid.VidmgrInfo{
			VidmgrID:     vidmgrInfo.VidmgrID,
			VidmgrStatus: def.DeviceStatusOnline,
			LastLogin:    time.Now().Unix(),
		}
		_, err = l.svcCtx.VidmgrM.VidmgrInfoUpdate(l.ctx, vidReq)
		if err != nil {
			er := errors.Fmt(err)
			l.Errorf("%s.rpc.ManageVidmgr req=%v err=%v", utils.FuncName(), req, er)
		}
	}
	return &types.HooksApiResp{
		Code: 0,
		Msg:  "success",
	}, nil
}
