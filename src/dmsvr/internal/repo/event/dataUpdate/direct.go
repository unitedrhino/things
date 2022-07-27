package dataUpdate

import (
	"context"
	"github.com/i-Things/things/src/dmsvr/internal/domain/schema"
)

type (
	Direct struct {
	}
)

func NewDirect() (*Direct, error) {
	return &Direct{}, nil
}

func (d *Direct) TempModelUpdate(ctx context.Context, info *schema.SchemaInfo) error {
	return nil
}

func (d *Direct) Subscribe(handle Handle) error {
	return nil
}
