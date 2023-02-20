package scene

import (
	"time"
)

type Info struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Desc        string    `json:"desc"`
	CreatedTime time.Time `json:"createdTime"`
	Trigger     *Trigger  `json:"trigger"`
	When        Terms     `json:"when"` //只有设备触发时才有用
	Then        Actions   `json:"then"`
}

func (i *Info) Validate() error {
	err := i.Trigger.Validate()
	if err != nil {
		return err
	}
	err = i.When.Validate()
	if err != nil {
		return err
	}
	err = i.Then.Validate()
	if err != nil {
		return err
	}
	return nil
}
