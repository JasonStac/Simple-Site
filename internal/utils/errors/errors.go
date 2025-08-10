package errors

import (
	"errors"
)

const (
	notFoundMessage string = "nothing found"
)

// type ErrNotFound struct {
// 	Code    int
// 	Message string
// }

// func (e *ErrNotFound) Error() string {
// 	return fmt.Sprintf("Error Code: %d: %s", e.Code, e.Message)
// }

// func (e *ErrNotFound) Is(target error) bool {
// 	if
// }

var ErrNotFound = errors.New(notFoundMessage)
