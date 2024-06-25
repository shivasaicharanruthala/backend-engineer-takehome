package errors

import (
	"fmt"
	"net/http"
	"time"
)

type EntityNotFound struct {
	Entity     string    `json:"-"`
	ID         string    `json:"-"`
	Msg        string    `json:"msg"`
	StatusCode int       `json:"-" default:"400"`
	TimeStamp  time.Time `json:"timestamp"`
}

func NewEntityNotFound(err error) EntityNotFound {
	return EntityNotFound{
		Msg:        err.Error(),
		StatusCode: http.StatusNotFound,
		TimeStamp:  time.Now().UTC(),
	}

}

func (e EntityNotFound) Error() string {
	return fmt.Sprintf("No '%v' found for Id: '%v'", e.Entity, e.ID)
}
