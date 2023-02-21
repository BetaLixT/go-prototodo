// All of the app level errors
// first digit identifies the layer (3 = infra)
// the first two digit identify the domain the error
// was created for 99 refers to a non domain specific error

package common

import "github.com/betalixt/gorr"

const (
	InvalidContextProvidedToHandlerErrorCode    = 1_99_000
	InvalidContextProvidedToHandlerErrorMessage = "InvalidContextProvidedToHandlerError"
)

func NewInvalidContextProvidedToHandlerError() *gorr.Error {
	return gorr.NewError(
		gorr.ErrorCode{
			Code:    InvalidContextProvidedToHandlerErrorCode,
			Message: InvalidContextProvidedToHandlerErrorMessage,
		},
		500,
		"",
	)
}
