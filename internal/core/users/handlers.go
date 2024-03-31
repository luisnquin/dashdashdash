package users

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (m Module) GetUserByUsernameHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := m.repo.users.FindOneByUsername(c.Param("username"))
		if err != nil {
			fmt.Println(err)
			c.Error(err)

			return err
		}

		return c.JSON(http.StatusOK, user)
	}
}
