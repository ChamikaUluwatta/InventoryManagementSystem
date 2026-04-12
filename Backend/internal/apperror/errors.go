package apperror

import (
	"net/http"
)

type AppError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	err        error  `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.err
}

func NotFound(message string, err error) *AppError {
	return &AppError{
		StatusCode: http.StatusNotFound,
		Message:    message,
		err:        err,
	}
}

func Internal(message string, err error) *AppError {
	return &AppError{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
		err:        err,
	}
}

func BadRequest(message string, err error) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
		err:        err,
	}
}

func Conflict(message string, err error) *AppError {
	return &AppError{
		StatusCode: http.StatusConflict,
		Message:    message,
		err:        err,
	}
}
