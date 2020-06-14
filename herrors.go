package herrors

import (
	"errors"
	"net/http"
)

var (
	ErrBadRequest     = New(http.StatusBadRequest)
	ErrUnauthorized   = New(http.StatusUnauthorized)
	ErrInternalServer = New(http.StatusInternalServerError)
)

// New constructs an ErrHttp error with the code and message populated
func New(statusCode uint) error {
	e := &ErrHttp{
		code: statusCode,
	}

	e.msg = http.StatusText(int(e.Code()))

	return e
}

// Wrap will create a new ErrHttp and place the provided error
// inside of it. There is a complementary errors.Unwrap to retrieve
// the wrapped error.
func Wrap(err error, statusCode uint) error {
	e := &ErrHttp{
		code:  statusCode,
		cause: err,
	}

	e.msg = http.StatusText(int(e.Code()))

	return e
}

// Write sets the response status code to the provided
// error code. When unable to find an ErrHttp error in
// the error chain then a 500 internal error is output.
func Write(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	var target hasCode
	// don't allow unknown errors to surface to user
	if !errors.As(err, &target) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(int(target.Code()))
}

type ErrHttp struct {
	msg   string
	code  uint
	cause error
}

// Code returns http status code that the error represents
func (e *ErrHttp) Code() uint {
	if e.code == 0 {
		return http.StatusInternalServerError
	}

	return e.code
}

// Error returns http friendly error message
func (e *ErrHttp) Error() string {
	return e.msg
}

// Unwrap will return the first wrapped error if there is one
func (e *ErrHttp) Unwrap() error {
	return e.cause
}

// Is used by errors.Is for comparing ErrHttp
func (e *ErrHttp) Is(target error) bool {
	h, ok := target.(hasCode)
	if !ok {
		return false
	}

	return e.Code() == h.Code()
}

type hasCode interface {
	Code() uint
}
