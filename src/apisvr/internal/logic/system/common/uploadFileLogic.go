package common

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/oss"
	"github.com/i-Things/things/shared/oss/common"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/spf13/cast"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *UploadFileLogic {
	return &UploadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *UploadFileLogic) UploadFile() (resp *types.UploadFileResp, err error) {
	business := l.r.FormValue("business")
	scene := l.r.FormValue("scene")
	filePath := l.r.FormValue("filePath")
	rename := cast.ToBool(l.r.FormValue("rename"))
	if business == "" || scene == "" || filePath == "" {
		return nil, errors.Parameter.WithMsg("business,scene,fileName是必填项")
	}
	fileDir, err := oss.GetFilePath(&oss.SceneInfo{
		Business: business,
		Scene:    scene,
		FilePath: filePath}, rename)
	if err != nil {
		l.Errorf("%s.GetFilePath err:%v", utils.FuncName(), err)
		return nil, err
	}
	file, _, err := l.r.FormFile("file")
	if err != nil {
		return resp, err
	}
	defer file.Close()
	newFilePath, err := l.svcCtx.OssClient.TemporaryBucket().Upload(l.ctx, fileDir, file, common.OptionKv{})
	if err != nil {
		return resp, err
	}
	return &types.UploadFileResp{
		FileUrl: newFilePath,
		FileDir: fileDir,
	}, err
}
