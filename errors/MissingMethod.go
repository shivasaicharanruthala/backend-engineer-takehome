package errors

import "fmt"

type MethodMissing struct {
	StatusCode int `default:"405"`
	Method     string
	URL        string
}

func (m MethodMissing) Error() string {
	return fmt.Sprintf("Method '%s' for '%s' not defined yet", m.Method, m.URL)
}
