package helpers

import (
	"strings"
	"errors"

	"gorm.io/gorm"
)

func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, gorm.ErrDuplicatedKey) ||
		strings.Contains(err.Error(), "23505") ||
		strings.Contains(err.Error(), "duplicate key")
}
