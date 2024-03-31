package nix

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/luisnquin/go-log"
)

func (m Module) GetInstalledUserEnvPackagesHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		packages, err := m.repo.nix.GetInstalledUserEnvPackages(c.Request().Context())
		if err != nil {
			log.Err(err).Msg("could not get installed user env packages")

			return c.JSON(http.StatusInternalServerError, "could not get installed user env packages")
		}

		return c.JSON(http.StatusOK, packages)
	}
}
