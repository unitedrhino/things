package projectmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
)

func checkProject(ctx context.Context, productID int64) (*relationDB.SysProjectInfo, error) {
	po, err := relationDB.NewProjectInfoRepo(ctx).FindOne(ctx, productID)
	if err == nil {
		return po, nil
	}
	if errors.Cmp(err, errors.NotFind) {
		return nil, nil
	}
	return nil, err
}
