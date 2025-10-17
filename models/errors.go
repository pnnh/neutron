package models

import (
	"errors"
	"fmt"
)

var ErrNilValue = fmt.Errorf("value is nil")

func IsErrNilValue(err error) bool {
	return errors.Is(err, ErrNilValue)
}
