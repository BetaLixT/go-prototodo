package common

import "github.com/betalixt/gorr"

// - Errors
// All of the domain level errors
// first digit identifies the layer (2 = domain)
// the first two digit identify the domain the error
// was created for 99 refers to a non domain specific error

func NewInvalidUserTypeForTaskError() *gorr.Error {
	return gorr.NewError(
		gorr.ErrorCode{
			Code:    2_00_000,
			Message: "InvalidUserTypeForTaskError",
		},
		404,
		"only users are allowed to create task",
	)
}
