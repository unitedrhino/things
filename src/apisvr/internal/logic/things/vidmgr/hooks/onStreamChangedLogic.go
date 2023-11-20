package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OnStreamChangedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOnStreamChangedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnStreamChangedLogic {
	return &OnStreamChangedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnStreamChangedLogic) OnStreamChanged(req *types.HooksApiStreamChangedRep) (resp *types.HooksApiResp, err error) {
	// todo: add your logic here and delete this line
	reqStr, _ := json.Marshal(*req)
	fmt.Println("---------OnStreamChanged--------------:", string(reqStr))
	//获取流注册时分支

	//需要先判断该流服务是否有注册过，未注册过
	vidStreamIndex, err := l.svcCtx.VidmgrS.VidmgrStreamIndex(l.ctx, &vid.VidmgrStreamIndexReq{
		App:        req.App,
		VidmgrID:   req.MediaServerId,
		Schema:     req.Schema,
		Vhost:      req.Vhost,
		Identifier: req.OriginSock.Identifier,
		LocalIP:    req.OriginSock.LocalIp,
		LocalPort:  req.OriginSock.LocalPort,
		PeerIP:     req.OriginSock.PeerIp,
		PeerPort:   req.OriginSock.PeerPort,
	})

	if err != nil {
		l.Errorf("%s rpc.VidmgrStreamIndex req=%v err=%+v", utils.FuncName(), req, err)
		return nil, err
	}
	vidStreamInfo := ToVidmgrStreamRpc(req)
	if req.Regist {
		vidStreamInfo.IsOnline = true
	} else { //注销时分支
		vidStreamInfo.IsOnline = false
	}
	//判断流路径是一样
	if len(vidStreamIndex.List) >= 1 {
		//判断Sock相同为同一流  	update
		fmt.Println("[--airgens--]we update stream:", vidStreamInfo)
		vidStreamInfo.StreamID = vidStreamIndex.List[0].StreamID
		vidStreamInfo.StreamName = vidStreamIndex.List[0].StreamName
		_, err := l.svcCtx.VidmgrS.VidmgrStreamUpdate(l.ctx, vidStreamInfo)
		if err != nil {
			l.Errorf("%s rpc.VidmgrStreamUpdate  err=%+v", utils.FuncName(), err)
			return nil, err
		}
	} else {
		fmt.Println("[--airgens--]we insert stream:", vidStreamInfo)
		_, err := l.svcCtx.VidmgrS.VidmgrStreamCreate(l.ctx, vidStreamInfo)
		if err != nil {
			l.Errorf("%s rpc.VidmgrStreamCreate  err=%+v", utils.FuncName(), err)
			return nil, err
		}
		//insert
	}

	return &types.HooksApiResp{
		Code: 0,
		Msg:  "success",
	}, nil
}
