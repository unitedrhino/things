package index

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetServerConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetServerConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetServerConfigLogic {
	return &SetServerConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetServerConfigLogic) SetServerConfig(req *types.IndexApiSetServerConfigReq) (resp *types.IndexApiSetServerConfigResp, err error) {
	// todo: add your logic here and delete this line
	dataRecv := new(types.IndexApiSetServerConfigResp)
	strmConfig := new(types.ServerConfig)
	err = json.Unmarshal([]byte(req.Data), strmConfig)
	if err != nil {
		fmt.Println("parse Json failed:", err)
		return dataRecv, err
	}
	strmConfig.GeneralMediaServerId = req.VidmgrID
	SetDefaultConfig(l.svcCtx.Config.Mediakit.Host, l.svcCtx.Config.Mediakit.Port, strmConfig)
	//set default
	byte4, err := json.Marshal(strmConfig)
	mdata, err := proxySetMediaServer(l.ctx, SETSERVERCONFIG, req.VidmgrID, byte4)
	err = json.Unmarshal(mdata, dataRecv)
	if err != nil {
		fmt.Println("parse Json failed:", err)
		return dataRecv, err
	}
	return dataRecv, err
}
