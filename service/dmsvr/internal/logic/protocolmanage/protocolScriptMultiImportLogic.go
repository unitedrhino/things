package protocolmanagelogic

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolScriptMultiImportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProtocolScriptMultiImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolScriptMultiImportLogic {
	return &ProtocolScriptMultiImportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProtocolScriptMultiImportLogic) ProtocolScriptMultiImport(in *dm.ProtocolScriptImportReq) (*dm.ImportResp, error) {
	var scripts []*relationDB.DmProtocolScript
	err := json.Unmarshal([]byte(in.Scripts), &scripts)
	if err != nil {
		return nil, err
	}
	var resp = dm.ImportResp{Total: int64(len(scripts))}
	for _, s := range scripts {
		ds := s.Devices
		s.ID = 0
		s.Devices = nil
		s.CreatedBy = 0
		s.UpdatedBy = 0
		err = relationDB.NewProtocolScriptRepo(l.ctx).Insert(l.ctx, s)
		if err != nil {
			l.Error(s, err)
			resp.ErrCount++
			continue
		}
		if len(ds) != 0 {
			for _, d := range ds {
				d.ID = 0
				d.ScriptID = s.ID
			}
			err = relationDB.NewProtocolScriptDeviceRepo(l.ctx).MultiInsert(l.ctx, ds)
			if err != nil {
				l.Error(s, err)
				resp.ErrCount++
				continue
			}
		}
		resp.SuccCount++
	}
	return &resp, nil
}
