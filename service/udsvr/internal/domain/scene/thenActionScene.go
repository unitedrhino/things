package scene

import (
	"context"
)

/*
	设备告警通知模版如下,需要在params中填写body,如: 电量过低
		你好,{{.deviceAlias}}设备告警:{{.body}}
*/

type ActionScene struct {
	SceneID int64 `json:"sceneID"`
}

func (a *ActionScene) Validate(repo ValidateRepo) error {
	if a == nil {
		return nil
	}
	return nil
}
func (a *ActionScene) Execute(ctx context.Context, repo ActionRepo) error {
	err := repo.SceneExec(ctx, a.SceneID)
	return err
}
