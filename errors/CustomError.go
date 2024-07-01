package errors

import "time"

type CustomError struct {
	Err        error     `json:"-"`
	Msg        string    `json:"msg"`
	StatusCode int       `json:"-" default:"500"`
	TimeStamp  time.Time `json:"timestamp"`
}

func NewCustomError(err error, statusCode ...int) CustomError {
	var code int
	if len(statusCode) == 0 {
		code = 500
	} else {
		code = statusCode[0]
	}

	if err.Error() == "dial tcp [::1]:5432: connect: connection refused" {
		code = 503
	}

	return CustomError{
		Err:        err,
		Msg:        err.Error(),
		StatusCode: code,
		TimeStamp:  time.Now().UTC(),
	}
}

func (ce CustomError) Error() string {
	if ce.Msg != "" {
		return ce.Msg
	}

	return ce.Err.Error()
}
