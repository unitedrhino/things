package modulemanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
)

func CheckModule(ctx context.Context, moduleCode string) error {
	c, err := relationDB.NewModuleInfoRepo(ctx).CountByFilter(ctx, relationDB.ModuleInfoFilter{Codes: []string{moduleCode}})
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.Parameter.AddMsgf("moduleCode not find:%v", moduleCode)
	}
	return nil
}
