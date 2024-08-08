package scene

import (
	"context"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
)

/*
	设备告警通知模版如下,需要在params中填写body,如: 电量过低
		你好,{{.deviceAlias}}设备告警:{{.body}}
*/

type ActionScene struct {
	AreaID    int64  `json:"areaID,string"` //仅做记录
	SceneID   int64  `json:"sceneID"`
	SceneName string `json:"sceneName"` //更新及创建的时候后端会赋值
}

func (a *ActionScene) Validate(repo CheckRepo) error {
	s, err := repo.GetSceneInfo(repo.Ctx, a.SceneID)
	if err != nil {
		return err
	}
	if s.Type != SceneTypeManual {
		return errors.Parameter.AddMsg("只能执行manual类型的场景")
	}
	a.SceneName = s.Name
	return nil
}
func (a *ActionScene) Execute(ctx context.Context, repo ActionRepo) error {
	err := repo.SceneExec(ctx, a.SceneID)
	func() {
		er := errors.Fmt(err)
		status := int64(def.True)
		if er.GetCode() != errors.OK.GetCode() {
			status = def.False
			repo.Info.Log.Status = def.False
		}
		repo.Info.Log.ActionMutex.Lock()
		defer repo.Info.Log.ActionMutex.Unlock()
		repo.Info.Log.Actions = append(repo.Info.Log.Actions, &LogAction{
			Type: ActionExecutorScene,
			Scene: &LogActionScene{
				SceneID:   a.SceneID,
				SceneName: a.SceneName,
			},
			Status: status,
			Code:   er.GetCode(),
			Msg:    er.GetMsg(),
		})
	}()
	return err
}
