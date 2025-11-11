package models

import (
	"errors"
	"fmt"
)

var ErrNilValue = fmt.Errorf("value is nil")
var ErrNotFound = fmt.Errorf("not found")

// IsErrNotFound 判断错误是否为未找到错误
func IsErrNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsErrNilValue 判断错误是否为值为空错误
func IsErrNilValue(err error) bool {
	return errors.Is(err, ErrNilValue)
}
