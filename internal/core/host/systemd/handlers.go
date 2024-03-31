package systemd

import (
	"context"
	"log"
	"net/http"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/labstack/echo/v4"
)

func (m Module) ListServicesHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		conn, err := dbus.NewWithContext(c.Request().Context())
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "unable to connect to dbus")
		}

		// dbus.NewUserConnectionContext()

		defer conn.Close()

		units, err := conn.ListUnitsContext(context.TODO())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "unabke to list units")
		}

		return c.JSON(http.StatusOK, units)
	}
}
