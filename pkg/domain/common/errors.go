// All of the domain level errors (errors used by and known to the domain)
// first digit identifies the layer (2 = domain)
// the first two digit identify the domain the error
// was created for 99 refers to a non domain specific error

package common

import "github.com/betalixt/gorr"

const (
	InvalidUserTypeForTaskErrorCode    = 2_00_000
	InvalidUserTypeForTaskErrorMessage = "InvalidUserTypeForTaskError"

	TaskMissingErrorCode = 2_00_001
	TaskMissingErrorMessage = "TaskMissingError"
)

func NewInvalidUserTypeForTaskError() *gorr.Error {
	return gorr.NewError(
		gorr.ErrorCode{
			Code:    InvalidUserTypeForTaskErrorCode,
			Message: InvalidUserTypeForTaskErrorMessage,
		},
		403,
		"only users are allowed to create task",
	)
}

func NewTaskMissingError() *gorr.Error {
	return gorr.NewError(
		gorr.ErrorCode{
			Code:    TaskMissingErrorCode,
			Message: TaskMissingErrorMessage,
		},
		403,
		"only users are allowed to create task",
	)
}
