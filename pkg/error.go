package pkg

import "fmt"

// Error represents an error with a code and message.
const (
	ErrCodeNotFound         = 1001
	ErrCodeFileNotFound     = 1002
	ErrCodeCreateFile       = 1003
	ErrCodeReadFile         = 1004
	ErrCodeWriteFile        = 1005
	ErrCodeSaveFile         = 1006
	ErrCodeResourceNotFound = 1007
	ErrCodeDirNotFound      = 1008
	ErrCodeUnsupportedOs    = 1009
	ErrCodeProcessFail      = 1010
)

type Error struct {
	Code    int
	Message string
}

// Error implements the error interface for Error.
func (e *Error) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// Error creates a new Error with the given code and message.
func ErrorStatus(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}
