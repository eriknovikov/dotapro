package errs

import "errors"

var (
	UNIMPLEMENTED = errors.New("UNIMPLEMENTED")
	NOT_FOUND     = errors.New("NOT FOUND")
)
