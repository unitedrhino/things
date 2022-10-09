package logic

import (
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/disvr/pb/di"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
)

func ToDmPageRpc(in *types.PageInfo) *dm.PageInfo {
	if in == nil {
		return nil
	}
	return &dm.PageInfo{
		Page: in.Page,
		Size: in.Size,
	}
}

func ToDiPageRpc(in *types.PageInfo) *di.PageInfo {
	if in == nil {
		return nil
	}
	return &di.PageInfo{
		Page: in.Page,
		Size: in.Size,
	}
}
