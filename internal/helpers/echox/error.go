package echox

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type ApiError struct {
	StatusCode int
	Data       any
}

func (a ApiError) Error() string {
	return fmt.Sprintf("%d: %v", a.StatusCode, a.Data)
}

func (a ApiError) JSON(c echo.Context) error {
	return c.JSON(a.StatusCode, a.Data)
}
