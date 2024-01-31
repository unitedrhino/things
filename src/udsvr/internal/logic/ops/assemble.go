package opslogic

import (
	"gitee.com/i-Things/core/shared/stores"
	"github.com/i-Things/things/src/udsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/udsvr/pb/ud"
)

func ToOpsWorkOrderPo(in *ud.OpsWorkOrder) *relationDB.UdOpsWorkOrder {
	if in == nil {
		return nil
	}
	return &relationDB.UdOpsWorkOrder{
		ID:          in.Id,
		AreaID:      stores.AreaID(in.AreaID),
		RaiseUserID: in.RaiseUserID,
		IssueDesc:   in.IssueDesc,
		Number:      in.Number,
		Type:        in.Type,
		Params:      in.Params,
		Status:      in.Status,
	}
}
