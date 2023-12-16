package loglogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperLogIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OlDB *relationDB.OperLogRepo
}

func NewOperLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperLogIndexLogic {
	return &OperLogIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OlDB:   relationDB.NewOperLogRepo(ctx),
	}
}

func (l *OperLogIndexLogic) OperLogIndex(in *sys.OperLogIndexReq) (*sys.OperLogIndexResp, error) {
	f := relationDB.OperLogFilter{
		OperName:     in.OperName,
		OperUserName: in.OperUserName,
		BusinessType: in.BusinessType,
	}
	resp, err := l.OlDB.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	total, err := l.OlDB.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	info := make([]*sys.OperLogInfo, 0, len(resp))
	for _, v := range resp {
		info = append(info, &sys.OperLogInfo{
			UserID:       v.OperUserID,
			AppCode:      v.AppCode,
			OperUserName: v.OperUserName,
			OperName:     v.OperName,
			BusinessType: v.BusinessType,
			Uri:          v.Uri,
			OperIpAddr:   v.OperIpAddr,
			OperLocation: v.OperLocation,
			Req:          v.Req.String,
			Resp:         v.Resp.String,
			Code:         v.Code,
			Msg:          v.Msg,
			CreatedTime:  v.CreatedTime.Unix(),
		})
	}

	return &sys.OperLogIndexResp{List: info, Total: total}, nil
}
