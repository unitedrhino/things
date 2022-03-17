package logic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/dcsvr/dc"
	"github.com/i-Things/things/src/dcsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupInfoLogic {
	return &GetGroupInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取组信息
func (l *GetGroupInfoLogic) GetGroupInfo(in *dc.GetGroupInfoReq) (resp *dc.GetGroupInfoResp, err error) {
	l.Infof("GetGroupInfo|req=%+v", in)
	var info []*dc.GroupInfo
	var size int64
	if in.Page == nil || in.Page.Page == 0 { //只获取一个
		gi, err := l.svcCtx.GroupInfo.FindOne(in.GroupID)
		if err != nil {
			return nil, err
		}
		info = append(info, DBToRPCFmt(gi).(*dc.GroupInfo))
	} else {
		size, err = l.svcCtx.DcDB.GetCountByGroupInfo()
		if err != nil {
			return nil, err
		}
		di, err := l.svcCtx.DcDB.FindByGroupInfo(def.PageInfo{PageSize: in.Page.PageSize, Page: in.Page.Page})
		if err != nil {
			return nil, err
		}
		info = make([]*dc.GroupInfo, 0, len(di))
		for _, v := range di {
			info = append(info, DBToRPCFmt(v).(*dc.GroupInfo))
		}
	}
	return &dc.GetGroupInfoResp{Info: info, Total: size}, nil
}
