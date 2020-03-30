package exception

import (
	"fmt"
)

//Exception interface
type Exception interface {
	Code() int
	Message() string
}

//exception data model
type exception struct {
	code    int
	message string
}

func (e *exception) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.code, e.message)
}

func (e *exception) Code() int {
	return e.code
}

func (e *exception) Message() string {
	return e.message
}

//New exception
func New(code int, message string, params ...interface{}) error {
	if params != nil && len(params) > 0 {
		message = fmt.Sprintf(message, params...)
	}
	return &exception{code, message}
}

//IsException is error of type exception
func IsException(err error) (Exception, bool) {
	if err == nil {
		return nil, false
	}

	exc, ok := err.(*exception)
	return exc, ok
}
