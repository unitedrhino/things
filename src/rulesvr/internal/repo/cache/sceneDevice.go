package cache

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/shared/def"
	"gitee.com/i-Things/core/shared/devices"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/rulesvr/internal/domain/scene"
	"sync"
)

type (
	SceneDeviceRepo struct {
		scene      scene.Repo
		triggerMap sync.Map
	}
)

func NewSceneDeviceRepo(sceneRepo scene.Repo) *SceneDeviceRepo {
	return &SceneDeviceRepo{
		scene: sceneRepo,
	}
}

func (s *SceneDeviceRepo) genKeyWithInfo(info *scene.Info) []string {
	var keys []string
	for _, t := range info.Trigger.Device {
		if utils.SliceIn(t.Operator, scene.DeviceOperationOperatorConnected, scene.DeviceOperationOperatorDisConnected) {
			key := fmt.Sprintf("%s.%s", t.ProductID, t.Operator)
			keys = append(keys, key)
			continue
		}
		key := fmt.Sprintf("%s.%s.%s", t.ProductID, t.Operator, t.OperationSchema.DataID[0])
		keys = append(keys, key)
	}
	return keys
}
func (s *SceneDeviceRepo) genKey(device devices.Core, operator scene.DeviceOperationOperator, dataID string) string {
	if utils.SliceIn(operator, scene.DeviceOperationOperatorConnected, scene.DeviceOperationOperatorDisConnected) {
		key := fmt.Sprintf("%s.%s", device.ProductID, operator)
		return key
	}
	key := fmt.Sprintf("%s.%s.%s", device.ProductID, operator, dataID)
	return key
}

func (s *SceneDeviceRepo) Init(ctx context.Context) error {
	s.triggerMap.Range(func(key, value any) bool {
		s.triggerMap.Delete(key)
		return true
	})
	infos, err := s.scene.FindByFilter(ctx, scene.InfoFilter{
		Status:      def.Enable,
		TriggerType: scene.TriggerTypeDevice,
	}, nil)
	if err != nil {
		return err
	}
	for _, info := range infos {
		err := s.Insert(ctx, info)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SceneDeviceRepo) Insert(ctx context.Context, info *scene.Info) error {
	keys := s.genKeyWithInfo(info)
	for _, k := range keys {

		temp, ok := s.triggerMap.Load(k)
		if !ok {
			s.triggerMap.Store(k, scene.Infos{info})
			return nil
		}
		v := temp.(scene.Infos)
		v = append(v, info)
		s.triggerMap.Store(k, v)
	}
	return nil
}

func (s *SceneDeviceRepo) GetInfos(ctx context.Context, device devices.Core,
	operator scene.DeviceOperationOperator, dataID string) (scene.Infos, error) {
	key := s.genKey(device, operator, dataID)
	temp, ok := s.triggerMap.Load(key)
	if !ok {
		return nil, nil
	}
	return temp.(scene.Infos), nil
}

func (s *SceneDeviceRepo) Update(ctx context.Context, info *scene.Info) error {
	s.Delete(ctx, info.ID)
	s.Insert(ctx, info)
	return nil
}

func (s *SceneDeviceRepo) Delete(ctx context.Context, id int64) error {
	//遍历删除所有该scene
	s.triggerMap.Range(func(key, value any) bool {
		scenes := value.(scene.Infos)
		var find bool
		for i, v := range scenes {
			if v.ID == id {
				scenes = append(scenes[:i], scenes[i+1:]...)
				find = true
				continue
			}
		}
		if find {
			s.triggerMap.Store(key, scenes)
		}
		return true
	})
	return nil
}
