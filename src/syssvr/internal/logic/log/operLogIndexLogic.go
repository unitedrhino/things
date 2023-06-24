package loglogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/mysql"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperLogIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOperLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperLogIndexLogic {
	return &OperLogIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OperLogIndexLogic) OperLogIndex(in *sys.OperLogIndexReq) (*sys.OperLogIndexResp, error) {
	resp, total, err := l.svcCtx.LogModel.OperLogIndex(l.ctx, &mysql.OperLogFilter{
		Page:         &def.PageInfo{Page: in.Page.Page, Size: in.Page.Size},
		OperName:     in.OperName,
		OperUserName: in.OperUserName,
		BusinessType: in.BusinessType,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}

	info := make([]*sys.OperLogIndexData, 0, len(resp))
	for _, v := range resp {
		info = append(info, &sys.OperLogIndexData{
			UserID:       v.OperUserID,
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

	return &sys.OperLogIndexResp{Info: info, Total: total}, nil
}
