package workOrder

import (
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/udsvr/pb/ud"
)

func ToOpsWorkOrderPb(in *types.OpsWorkOrder) *ud.OpsWorkOrder {
	return utils.Copy[ud.OpsWorkOrder](in)
}
