package pkg

import "errors"

var (
	ErrBadRequest       = errors.New("bad request")
	ErrNotFound         = errors.New("not found")
	ErrFileNotFound     = errors.New("file not found")
	ErrCreateFile       = errors.New("create file")
	ErrReadFile         = errors.New("read file")
	ErrWriteFile        = errors.New("write file")
	ErrSaveFile         = errors.New("save file")
	ErrDirNotFound      = errors.New("directory not found")
	ErrCreateDir        = errors.New("create directory")
	ErrResourceNotFound = errors.New("resource not found")
	ErrUnsupportedOs    = errors.New("unsupported os")
	ErrInternalFailure  = errors.New("internal failure")
	ErrProcessFail      = errors.New("process fail")
)

type Error struct {
	appError error
	svcError error
}

func NewError(svcError, appError error) error {
	return Error{
		svcError: svcError,
		appError: appError,
	}
}

func (e Error) AppError() error {
	return e.appError
}

func (e Error) SvcError() error {
	return e.svcError
}

func (e Error) Error() string {
	return errors.Join(e.svcError, e.appError).Error()
}
