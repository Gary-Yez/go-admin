package utils

import (
	"errors"
	"gorm.io/gorm"
)

func IsForeignKeyConstraintError(err error) bool {
	return errors.Is(err, gorm.ErrForeignKeyViolated)
}
