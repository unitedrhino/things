package scene

import (
	"context"
	"gitee.com/i-Things/share/errors"
)

/*
	设备告警通知模版如下,需要在params中填写body,如: 电量过低
		你好,{{.deviceAlias}}设备告警:{{.body}}
*/

type ActionScene struct {
	AreaID  int64 `json:"areaID,string"` //仅做记录
	SceneID int64 `json:"sceneID"`
}

func (a *ActionScene) Validate(repo ValidateRepo) error {
	s, err := repo.GetSceneInfo(repo.Ctx, a.SceneID)
	if err != nil {
		return err
	}
	if s.Type != SceneTypeManual {
		return errors.Parameter.AddMsg("只能执行manual类型的场景")
	}
	return nil
}
func (a *ActionScene) Execute(ctx context.Context, repo ActionRepo) error {
	err := repo.SceneExec(ctx, a.SceneID)
	return err
}
