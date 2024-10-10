package deviceTimer

import (
	"gitee.com/unitedrhino/things/service/udsvr/internal/domain/scene"
	"time"
)

// 单设备的场景控制,其中所有带设备的参数都可以不填,系统会自动填充,同时只能操作一个设备
type DeviceInfo struct {
	ID          int64           `json:"id"`
	HeadImg     string          `json:"headImg"` // 头像
	Tag         string          `json:"tag"`
	Name        string          `json:"name"`
	Desc        string          `json:"desc"`
	CreatedTime time.Time       `json:"createdTime"`
	Type        scene.SceneType `json:"type"`
	If          scene.If        `json:"if"`     //多种触发方式
	When        scene.When      `json:"when"`   //手动触发模式不生效
	Then        scene.Then      `json:"then"`   //触发后执行的动作
	Status      int64           `json:"status"` // 状态（1启用 2禁用）
}
