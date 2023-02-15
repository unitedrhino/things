package scene

import (
	"github.com/i-Things/things/shared/errors"
	"time"
)

type InfoDo struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Desc        string    `json:"desc"`
	CreatedTime time.Time `json:"createdTime"`
	Trigger     *Trigger  `json:"trigger"`
	When        []*Term   `json:"when"`
	Then        []*Action `json:"then"`
}

type CreateInfoDto struct {
	Name        string    `json:"name"`
	Desc        string    `json:"desc"`
	createdTime time.Time `json:"createdTime"`
	Trigger     Trigger   `json:"trigger"`
	When        Term      `json:"when"`
	Action      Action    `json:"action"`
}

func NewInfo(dto CreateInfoDto) (*InfoDo, error) {
	if dto.Name == "" {
		return nil, errors.Parameter.AddMsg("场景名不能为空")
	}
	return nil, nil
}
