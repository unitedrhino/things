package scene

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestTrigger(t1 *testing.T) {
	do := Trigger{
		Type:   TriggerTypeAuto,
		Timers: Timers{&Timer{ExecAt: 123, ExecRepeat: 0b1111111}, &Timer{ExecAt: 223, ExecRepeat: 0b1111111}},
	}
	text, _ := json.Marshal(do)
	fmt.Println(string(text))
}
