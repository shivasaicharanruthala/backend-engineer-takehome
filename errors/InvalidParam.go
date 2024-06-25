package errors

import (
	"fmt"
	"net/http"
	"time"
)

type InvalidParam struct {
	Param      string    `json:"-"`
	Msg        string    `json:"msg"`
	StatusCode int       `json:"-" default:"400"`
	TimeStamp  time.Time `json:"timestamp"`
}

func NewInvalidParam(err error) InvalidParam {
	return InvalidParam{
		Msg:        err.Error(),
		StatusCode: http.StatusBadRequest,
		TimeStamp:  time.Now().UTC(),
	}
}
func (e InvalidParam) Error() string {
	return fmt.Sprintf("Incorrect value for parameter: " + e.Param)
}
