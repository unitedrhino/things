package devicegrouplogic

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/oss"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	GiDB *relationDB.GroupInfoRepo
}

func NewGroupInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInfoUpdateLogic {
	return &GroupInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		GiDB:   relationDB.NewGroupInfoRepo(ctx),
	}
}

// 更新分组
func (l *GroupInfoUpdateLogic) GroupInfoUpdate(in *dm.GroupInfo) (*dm.Empty, error) {
	po, err := l.GiDB.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	po.Desc = in.Desc
	po.Name = in.Name
	po.Tags = in.Tags
	if po.Tags == nil {
		po.Tags = map[string]string{}
	}
	if in.Files != nil {
		if po.Files == nil {
			po.Files = map[string]string{}
		}
		var files = map[string]string{}
		for key, v := range in.Files {
			if v == "" {
				files[key] = ""
			}
			if !oss.IsFilePath(l.svcCtx.Config.OssConf, v) { //传入的不是file path,不更新
				files[key] = po.Files[key]
				continue
			}
			nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessDeviceGroup, oss.SceneFile, fmt.Sprintf("%d/%s/%s", po.ID, key, oss.GetFileNameWithPath(v)))
			path, err := l.svcCtx.OssClient.PrivateBucket().CopyFromTempBucket(v, nwePath)
			if err != nil {
				l.Error(err)
				files[key] = po.Files[key]
				continue
			}
			files[key] = path
		}
		po.Files = files
	}
	err = l.GiDB.Update(l.ctx, po)
	if err != nil {
		return nil, errors.Parameter.AddMsg(err.Error())
	}

	return &dm.Empty{}, nil
}
