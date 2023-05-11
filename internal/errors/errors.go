package errors

import "errors"

var ErrValidation error = errors.New("validation error")
var ErrDatabaseUnreachable error = errors.New("database unreachable")
