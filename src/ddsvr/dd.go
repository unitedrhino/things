package main

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/ddsvr/dddirect"
)

func main() {
	defer utils.Recover(context.Background())
	dddirect.NewDd()
}
