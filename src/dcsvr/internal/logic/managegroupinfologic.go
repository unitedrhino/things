package logic

import (
	"context"
	"database/sql"
	"gitee.com/godLei6/things/shared/def"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/src/dcsvr/model"
	"github.com/spf13/cast"
	"time"

	"gitee.com/godLei6/things/src/dcsvr/dc"
	"gitee.com/godLei6/things/src/dcsvr/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

type ManageGroupInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewManageGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ManageGroupInfoLogic {
	return &ManageGroupInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *ManageGroupInfoLogic) AddGroupInfo(in *dc.ManageGroupInfoReq) (*dc.GroupInfo, error) {
	info := in.Info
	GroupID := l.svcCtx.GroupID.GetSnowflakeId()
	gi := &model.GroupInfo{
		GroupID     :GroupID, // 组id
		Name        :info.Name,    // 组名
		Uid         :info.Uid,     // 管理员用户id
		CreatedTime :time.Now(),
	}
	_,err := l.svcCtx.GroupInfo.Insert(*gi)
	if err != nil {
		l.Errorf("AddGroupInfo|GroupInfo|Insert|err=%+v", err)
		return nil, errors.System.AddDetail(err.Error())
	}
	return DBToRPCFmt(gi).(*dc.GroupInfo), nil

}
func (l *ManageGroupInfoLogic) ModifyGroupInfo(in *dc.ManageGroupInfoReq) (*dc.GroupInfo, error) {
	gi,err := l.svcCtx.GroupInfo.FindOne(in.Info.GroupID)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.Parameter.AddDetail("not find GroupID id:" + cast.ToString(in.Info.GroupID))
		}
		return nil, errors.System.AddDetail(err.Error())
	}
	gi.Name=in.Info.Name
	gi.Uid=in.Info.Uid
	gi.UpdatedTime = sql.NullTime{Valid: true, Time: time.Now()}
	err = l.svcCtx.GroupInfo.Update(*gi)
	if err != nil {
		l.Errorf("ModifyGroupInfo|GroupInfo|Update|err=%+v", err)
		return nil, errors.System.AddDetail(err.Error())
	}
	return DBToRPCFmt(gi).(*dc.GroupInfo), nil
}
func (l *ManageGroupInfoLogic) DelGroupInfo(in *dc.ManageGroupInfoReq) (*dc.GroupInfo, error) {
	err := l.svcCtx.GroupInfo.Delete(in.Info.GroupID)
	if err != nil {
		l.Errorf("DelGroupInfo|GroupInfo|Delete|err=%+v", err)
		return nil, errors.System.AddDetail(err.Error())
	}
	return &dc.GroupInfo{}, nil
}
// 管理组
func (l *ManageGroupInfoLogic) ManageGroupInfo(in *dc.ManageGroupInfoReq) (*dc.GroupInfo, error) {
	l.Infof("ManageProduct|opt=%d|req=%+v", in.Opt, in)
	switch in.Opt {
	case def.OPT_ADD:
		if in.Info == nil {
			return nil, errors.Parameter.WithMsg("add opt need info")
		}
		return l.AddGroupInfo(in)
	case def.OPT_MODIFY:
		return l.ModifyGroupInfo(in)
	case def.OPT_DEL:
		return l.DelGroupInfo(in)
	default:
		return nil, errors.Parameter.AddDetail("not suppot opt:" + cast.ToString(in.Opt))
	}

	return &dc.GroupInfo{}, nil
}
