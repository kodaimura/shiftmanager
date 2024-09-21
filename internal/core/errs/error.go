package errs

import (
	"fmt"
)


/////////////////////////////////////////////////////////////////////////
type UniqueConstraintError struct {
	Column string
}

func NewUniqueConstraintError(column string) error {
	return UniqueConstraintError{Column: column}
}

func (e UniqueConstraintError) Error() string {
	return fmt.Sprintf("UniqueConstraintError: %s", e.Column)
}

/////////////////////////////////////////////////////////////////////////
type NotFoundError struct {}

func NewNotFoundError() error {
	return NotFoundError{}
}

func (e NotFoundError) Error() string {
	return "NotFoundError"
}

/////////////////////////////////////////////////////////////////////////
type AlreadyRegisteredError struct {}

func NewAlreadyRegisteredError() error {
	return AlreadyRegisteredError{}
}

func (e AlreadyRegisteredError) Error() string {
	return "AlreadyRegisteredError"
}