package errors

import (
	"fmt"
	"net/http"
	"time"
)

type MissingParam struct {
	Param      string    `json:"-"`
	Msg        string    `json:"msg"`
	StatusCode int       `json:"-" default:"400"`
	TimeStamp  time.Time `json:"timestamp"`
}

func NewMissingParam(err error) MissingParam {
	return MissingParam{
		Msg:        err.Error(),
		StatusCode: http.StatusBadRequest,
		TimeStamp:  time.Now().UTC(),
	}
}
func (e MissingParam) Error() string {
	return fmt.Sprintf("Parameter " + e.Param + " is required for this request")
}
