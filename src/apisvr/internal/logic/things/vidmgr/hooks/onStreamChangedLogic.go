package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/vidsvr/pb/vid"
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

// [fix-bug by wfj]   不同IP推上来的相同路径的流，可能会有问题
func (l *OnStreamChangedLogic) OnStreamChanged(req *types.HooksApiStreamChangedRep) (resp *types.HooksApiResp, err error) {
	// todo: add your logic here and delete this line
	reqStr, _ := json.Marshal(*req)
	fmt.Println("---------OnStreamChanged--------------:", string(reqStr))

	//需要先判断该流服务是否有注册过，未注册过忽略消息
	vidInfo, err := l.svcCtx.VidmgrM.VidmgrInfoRead(l.ctx, &vid.VidmgrInfoReadReq{
		VidmgrtID: req.MediaServerId,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s rpc.ManageVidmgr req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	if vidInfo != nil {
		//查找四要素：vidmgr_id  app  stream   vhost
		vidStreamIndex, err := l.svcCtx.VidmgrS.VidmgrStreamIndex(l.ctx, &vid.VidmgrStreamIndexReq{
			App:      req.App,
			VidmgrID: req.MediaServerId,
			Stream:   req.Stream,
			Vhost:    req.Vhost,
			//具体连接
			//Identifier: req.OriginSock.Identifier,
			//LocalIP:    req.OriginSock.LocalIp,
			//LocalPort:  req.OriginSock.LocalPort,
			//PeerIP:     req.OriginSock.PeerIp,
			//PeerPort:   req.OriginSock.PeerPort,
		})

		if err != nil {
			l.Errorf("%s rpc.VidmgrStreamIndex req=%v err=%+v", utils.FuncName(), req, err)
			return nil, err
		}
		//fmt.Println("vidStreamIndex :", vidStreamIndex)
		if vidStreamIndex.Total > 0 {
			//查到流之后需要修改流
			//判断Sock相同为同一流  	update
			if req.Regist {
				//对应位(bit)  置1
				vidStreamIndex.List[0].Protocol |= GetProtocol(req.Schema)
				//fmt.Printf("ADD-Read[--airgens--]   0x:%x\n", vidStreamIndex.List[0].Protocol)
				//fmt.Printf("ADD-Read[--airgnes--]  val:%d\n)", vidStreamIndex.List[0].Protocol)
			} else {
				//对应位(bit) 置0
				vidStreamIndex.List[0].Protocol = vidStreamIndex.List[0].Protocol &^ GetProtocol(req.Schema)
			}
			if vidStreamIndex.List[0].Protocol == 0 {
				vidStreamIndex.List[0].IsOnline = false
			}
			_, err := l.svcCtx.VidmgrS.VidmgrStreamUpdate(l.ctx, vidStreamIndex.List[0])
			if err != nil {
				l.Errorf("%s rpc.VidmgrStreamUpdate  err=%+v", utils.FuncName(), err)
				return nil, err
			}
		} else {
			//如果没有查到流
			if req.Regist { //注册请求
				vidStreamInfo := ToVidmgrStreamRpc(req)
				vidStreamInfo.IsOnline = true //设置状态为在线
				vidStreamInfo.Protocol = GetProtocol(req.Schema)
				//fmt.Println("ADD-Write[--airgnes--]GetProtocol(req.Schema)", req.Schema)
				//fmt.Printf("ADD-Write[--airgnes--] 0X:%x\n)", vidStreamInfo.Protocol)
				//fmt.Printf("ADD-Write[--airgnes--]val:%d\n)", vidStreamInfo.Protocol)
				_, err := l.svcCtx.VidmgrS.VidmgrStreamCreate(l.ctx, vidStreamInfo)
				if err != nil {
					l.Errorf("%s rpc.VidmgrStreamCreate  err=%+v", utils.FuncName(), err)
					return nil, err
				}
			}
		}
	}
	return &types.HooksApiResp{
		Code: 0,
		Msg:  "success",
	}, nil
}
