package scene

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestThen(t1 *testing.T) {
	do := Then{
		Actions: Actions{&Action{
			Type:  ActionExecutorDelay,
			Delay: 5,
		}, &Action{
			Type: ActionExecutorDevice,
			Device: &ActionDevice{
				ProductID:  "28c6O38q2K4",
				SelectType: SelectorDeviceAll,
				Type:       ActionDeviceTypePropertyControl,
				DataID:     "switch",
				Value:      "true",
			},
		}},
	}
	text, _ := json.Marshal(do)
	fmt.Println(string(text))
}
