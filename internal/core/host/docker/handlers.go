package docker

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

func (m Module) GetContainersHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		status := c.QueryParam("status")

		if status != "" && !lo.Contains(getContainerStatusList(), status) {
			return c.JSON(http.StatusBadRequest, "invalid container status")
		}

		containers, err := m.repo.docker.ListContainers(ctx, status)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "unable to list containers")
		}

		return c.JSON(http.StatusOK, containers)
	}
}
