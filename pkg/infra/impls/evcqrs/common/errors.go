// All of the domain level errors (errors used by and known to the domain)
// first digit identifies the layer (3 = infra)
// the first two digit identify the domain the error
// was created for 99 refers to a non domain specific error

package common

import "github.com/betalixt/gorr"

const (
	FailedToAssertContextTypeErrorCode    = 3_99_000
	FailedToAssertContextTypeErrorMessage = "FailedToAssertContextTypeError"

	FailedToAssertDatabaseCtxTypeErrorCode    = 3_99_000
	FailedToAssertDatabaseCtxTypeErrorMessage = "FailedToAssertContextTypeError"
)

func NewFailedToAssertContextTypeError() *gorr.Error {
	return gorr.NewError(
		gorr.ErrorCode{
			Code:    FailedToAssertContextTypeErrorCode,
			Message: FailedToAssertContextTypeErrorMessage,
		},
		403,
		"",
	)
}

func NewFailedToAssertDatabaseCtxTypeError() *gorr.Error {
	return gorr.NewError(
		gorr.ErrorCode{
			Code:    FailedToAssertDatabaseCtxTypeErrorCode,
			Message: FailedToAssertDatabaseCtxTypeErrorMessage,
		},
		403,
		"",
	)
}
