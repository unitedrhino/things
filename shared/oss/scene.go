package oss

import (
	"fmt"
	"github.com/golang-module/carbon/v2"
	"github.com/hashicorp/go-uuid"
	"github.com/i-Things/things/shared/errors"
	"path"
	"strings"
)

type (
	SceneInfo struct {
		Business string
		Scene    string
		FilePath string
		FileName string
	}
)

// 产品管理
const (
	BusinessProductManage = "productManage" //产品管理
	SceneProductImg       = "productImg"    //产品图片
)

var (
	//Key 是业务 value是场景
	SceneRecord = map[string]map[string]struct{}{
		BusinessProductManage: {SceneProductImg: struct{}{}},
	}
)

func GetSceneInfo(filePath string) (*SceneInfo, error) {
	paths := strings.Split(filePath, "/")
	if len(paths) < 3 {
		return nil, errors.Parameter.WithMsg("路径不对")
	}
	scene := &SceneInfo{
		Business: paths[0],
		Scene:    paths[1],
		FilePath: strings.Join(paths[2:], "/"),
		FileName: paths[len(paths)-1],
	}
	return scene, CheckSceneInfo(scene)
}

func CheckSceneInfo(info *SceneInfo) error {
	scene := SceneRecord[info.Business]
	if scene == nil {
		return errors.Parameter.WithMsgf("business is not right:%v", info.Business)
	}
	if _, ok := scene[info.Scene]; !ok {
		return errors.Parameter.WithMsgf("scene is not right:%v", info.Scene)
	}
	return nil
}

func GetFilePath(scene *SceneInfo, rename bool) (string, error) {
	if rename == true {
		ext := path.Ext(scene.FilePath)
		if ext == "" {
			return "", errors.Parameter.WithMsg("未能获取文件后缀名")
		}
		uuid, er := uuid.GenerateUUID()
		if er != nil {
			err := errors.System.AddDetail(er)
			return "", err
		}
		scene.FilePath = uuid + ext
	} else {
		spcChar := []string{`,`, `?`, `*`, `|`, `{`, `}`, `\`, `/`, `$`, `、`, `·`, "`", `'`, `"`}
		if strings.ContainsAny(scene.FilePath, strings.Join(spcChar, "")) {
			return "", errors.Parameter.WithMsg("包含特殊字符")
		}
	}
	if scene.Date == "" {
		scene.Date = carbon.Now().ToDateString()
	}
	filePath := fmt.Sprintf("%s/%s/%s", scene.Business, scene.Scene, scene.FilePath)
	return filePath, nil
}
