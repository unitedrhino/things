package scene

import (
	"context"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
)

/*
	设备告警通知模版如下,需要在params中填写body,如: 电量过低
		你好,{{.deviceAlias}}设备告警:{{.body}}
*/

type ActionScene struct {
	AreaID    int64     `json:"areaID,string,omitempty"` //仅做记录
	SceneID   int64     `json:"sceneID,omitempty"`
	SceneType SceneType `json:"sceneType,omitempty"` //类型,后端会赋值
	SceneName string    `json:"sceneName,omitempty"` //更新及创建的时候后端会赋值
	Status    def.Bool  `json:"status,omitempty"`    //如果是自动化类型则为修改状态
}

func (a *ActionScene) Validate(repo CheckRepo) error {
	s, err := repo.GetSceneInfo(repo.Ctx, a.SceneID)
	if err != nil {
		return err
	}
	if s.Type == SceneTypeAuto && utils.SliceIn(a.Status, def.False, def.False) {
		return errors.Parameter.AddMsg("自动化类型需要填写修改的状态值")
	} else {
		a.Status = 0
	}
	a.SceneType = s.Type
	a.SceneName = s.Name
	return nil
}
func (a *ActionScene) Execute(ctx context.Context, repo ActionRepo) error {
	err := repo.SceneExec(ctx, a.SceneID, a.Status)
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
