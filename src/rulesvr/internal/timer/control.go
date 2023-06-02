package timer

import (
	"github.com/i-Things/things/src/rulesvr/internal/domain/scene"
)

type SceneControl interface {
	Create(info *scene.Info) error
	Update(info *scene.Info) error
	IsRunning() bool
	Delete(id int64) error
}
