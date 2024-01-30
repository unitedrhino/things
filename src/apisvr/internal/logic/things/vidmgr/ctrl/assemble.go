package ctrl

import (
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/src/apisvr/internal/types"
)

func handleZLMediakitReq(req *types.CtrlApiReq) ([]byte, error) {
	mgr := &clients.SvcZlmedia{
		Secret: req.Secret,
		Port:   req.Port,
		IP:     req.IP,
	}
	bytetmp := make([]byte, 0)
	return clients.ProxyMediaServer(req.Cmd, mgr, bytetmp)
}
