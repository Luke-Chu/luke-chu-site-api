package apperr

import (
	"errors"
	"fmt"
)

type AppError struct {
	HTTPStatus int
	Code       int
	Message    string
	Cause      error
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func New(httpStatus, code int, message string) *AppError {
	return &AppError{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
	}
}

func Wrap(httpStatus, code int, message string, cause error) *AppError {
	return &AppError{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
		Cause:      cause,
	}
}

func As(err error) (*AppError, bool) {
	var target *AppError
	ok := errors.As(err, &target)
	return target, ok
}
