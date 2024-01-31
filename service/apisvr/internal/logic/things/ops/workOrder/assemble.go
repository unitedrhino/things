package workOrder

import (
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/udsvr/pb/ud"
)

func ToOpsWorkOrderPb(in *types.OpsWorkOrder) *ud.OpsWorkOrder {
	if in == nil {
		return nil
	}
	return &ud.OpsWorkOrder{
		Id:          in.ID,
		AreaID:      in.AreaID,
		RaiseUserID: in.RaiseUserID,
		IssueDesc:   in.IssueDesc,
		Number:      in.Number,
		Type:        in.Type,
		Params:      in.Params,
		Status:      in.Status,
	}
}
