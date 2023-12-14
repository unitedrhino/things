package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"regexp"
	"time"
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
	infoRepo := relationDB.NewVidmgrInfoRepo(l.ctx)
	vidInfo, err := infoRepo.FindOneByFilter(l.ctx, relationDB.VidmgrFilter{
		VidmgrIDs: []string{req.MediaServerId},
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s rpc.ManageVidmgr req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	if vidInfo != nil {
		//查找四要素：vidmgr_id  app  stream   vhost
		streamRepo := relationDB.NewVidmgrStreamRepo(l.ctx)
		vidStreamIndex, err := streamRepo.FindAllFilter(l.ctx, relationDB.VidmgrStreamFilter{
			App:      req.App,
			VidmgrID: req.MediaServerId,
			Stream:   req.Stream,
			Vhost:    req.Vhost,
		})

		if err != nil {
			l.Errorf("%s rpc.VidmgrStreamIndex req=%v err=%+v", utils.FuncName(), req, err)
			return nil, err
		}
		//fmt.Println("vidStreamIndex :", vidStreamIndex)
		if len(vidStreamIndex) > 0 {
			//查到流之后需要修改流
			//判断Sock相同为同一流  	update
			if req.Regist {
				//对应位(bit)  置1
				vidStreamIndex[0].Protocol |= GetProtocol(req.Schema)
				vidStreamIndex[0].LastLogin = time.Now()
				fmt.Printf("ADD-Read[--airgens--Register]   0x:%x\n", vidStreamIndex[0].Protocol)
				fmt.Printf("ADD-Read[--airgnes--Register]  val:%d\n)", vidStreamIndex[0].Protocol)
			} else {
				//对应位(bit) 置0
				fmt.Printf("ADD-Read[--airgnes--old       ]  val:%d\n)", vidStreamIndex[0].Protocol)
				vidStreamIndex[0].Protocol = vidStreamIndex[0].Protocol &^ GetProtocol(req.Schema)
				fmt.Printf("ADD-Read[--airgens--UnRegister]   0x:%x\n", GetProtocol(req.Schema))
				fmt.Printf("ADD-Read[--airgens--UnRegister]   0x:%x\n", vidStreamIndex[0].Protocol)
				fmt.Printf("ADD-Read[--airgnes--UnRegister]  val:%d\n)", vidStreamIndex[0].Protocol)
			}
			if vidStreamIndex[0].Protocol == 0 {
				vidStreamIndex[0].IsOnline = false
			}
			err := streamRepo.Update(l.ctx, vidStreamIndex[0])
			if err != nil {
				l.Errorf("%s rpc.VidmgrStreamUpdate  err=%+v", utils.FuncName(), err)
				return nil, err
			}
		} else {
			//如果没有查到流
			if req.Regist { //注册请求
				vidStreamInfo := ToVidmgrStreamRpc(req)
				vidStreamInfo.IsOnline = true //设置状态为在线
				vidStreamInfo.FirstLogin = time.Now()
				vidStreamInfo.LastLogin = time.Now()
				vidStreamInfo.Protocol = GetProtocol(req.Schema)

				if vidStreamInfo.OriginType == RTMP_PUSH {
					re := regexp.MustCompile(vidStreamInfo.Vhost)
					if vidInfo.MediasvrType == 1 { //docker 模式
						//l.svcCtx.Config.Restconf.Host
						vidStreamInfo.OriginUrl =
							re.ReplaceAllString(vidStreamInfo.OriginUrl, l.svcCtx.Config.Restconf.Host)
					} else {
						//vidStreamInfo.PeerIP
						vidStreamInfo.OriginUrl =
							re.ReplaceAllString(vidStreamInfo.OriginUrl, req.OriginSock.PeerIp)
					}
				}

				err := streamRepo.Insert(l.ctx, vidStreamInfo)
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
