package systemd

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

func (m Module) ListUnitsHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		unitStatus, unitType := c.QueryParam("status"), c.QueryParam("type")
		scope := parseListUnitsScopeParam(c.QueryParam("scope"))

		if unitType != "" && !lo.Contains(getAllUnitTypes(), unitType) {
			return c.JSON(http.StatusBadRequest, "invalid unit type")
		} else if unitStatus != "" && !lo.Contains(getAllUnitStatuses(), unitStatus) {
			return c.JSON(http.StatusBadRequest, "invalid unit status")
		}

		ctx := c.Request().Context()

		var optFns []filterOption

		if unitStatus != "" {
			optFns = append(optFns, WithUnitStatus(unitStatus))
		}

		if unitType != "" {
			optFns = append(optFns, WithUnitType(unitType))
		}

		switch scope {
		case SYSTEM_ONLY:
			units, err := m.repo.systemd.ListSystemUnits(ctx, optFns...)
			if err != nil {
				log.Println(err)

				return c.JSON(http.StatusInternalServerError, err)
			}

			return c.JSON(http.StatusOK, units)

		case USER_ONLY:
			units, err := m.repo.systemd.ListUserUnits(ctx, optFns...)
			if err != nil {
				log.Println(err)

				return c.JSON(http.StatusInternalServerError, err)
			}

			return c.JSON(http.StatusOK, units)

		default:
			systemUnits, err := m.repo.systemd.ListSystemUnits(ctx, optFns...)
			if err != nil {
				log.Println(err)

				return c.JSON(http.StatusInternalServerError, err)
			}

			userUnits, err := m.repo.systemd.ListUserUnits(ctx, optFns...)
			if err != nil {
				log.Println(err)

				return c.JSON(http.StatusInternalServerError, err)
			}

			return c.JSON(http.StatusOK, echo.Map{
				"user":   userUnits,
				"system": systemUnits,
			})
		}
	}
}
