package project

import (
	"github.com/i-Things/things/src/viewsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/viewsvr/internal/types"
)

func ToProjectInfoTypes(p *relationDB.ViewProjectInfo) *types.ProjectInfo {
	return &types.ProjectInfo{
		IndexImage:    p.IndexImage,
		Name:          p.Name,
		Desc:          p.Desc,
		CreatedUserID: p.CreatedUserID,
		Status:        p.Status,
	}
}
