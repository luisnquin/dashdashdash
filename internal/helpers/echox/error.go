package echox

import "fmt"

type ApiError struct {
	StatusCode int
	Data       any
}

func (a ApiError) Error() string {
	return fmt.Sprintf("%d: %v", a.StatusCode, a.Data)
}
