package store

import (
	"github.com/i-Things/things/shared/errors"
	"gorm.io/gorm"
	"strings"
)

func ErrFmt(err error) error {
	if err == nil {
		return nil
	}
	if err == gorm.ErrRecordNotFound {
		return errors.NotFind.WithStack()
	}
	if strings.Contains(err.Error(), "Duplicate entry") {
		return errors.Duplicate.AddDetail(err)
	}
	return errors.Database.AddDetail(err)
}
