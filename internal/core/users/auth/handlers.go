package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/luisnquin/dashdashdash/internal/helpers/echox"
	"github.com/luisnquin/dashdashdash/internal/helpers/reasons"
	"github.com/luisnquin/dashdashdash/internal/models"
	"github.com/luisnquin/go-log"
)

func (m Module) GenerateTOTPUriHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, _ := c.Get("user").(*models.User)

		return c.JSON(http.StatusOK, GenerateTOPTURIResponse{
			URI: m.totp.ProvisioningUri(user.Username, m.config.GetIssuerName()),
		})
	}
}

func (m Module) ValidateTOTPCodeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		code := strings.TrimSpace(c.Param("code"))

		log.Debug().Str("provided_code", code).Str("current_totp_code", m.totp.Now()).Send()

		if !m.totp.Verify(code, time.Now().UTC().Unix()) {
			log.Warn().Msg("provided code wasn't valid :/")

			return c.JSON(http.StatusUnauthorized, ValidateTOPTCodeResponse{
				IsValid:    false,
				Reason:     "code doesn't match with the expected",
				ReasonCode: reasons.SESSION_INV_CREDS,
			})
		}

		return c.JSON(http.StatusOK, ValidateTOPTCodeResponse{
			IsValid: true,
		})
	}
}

func (m Module) LoginHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		signedToken, err := doBasicAuth(c, m.config, m.repo.auth)
		if err != nil {
			apiErr, ok := err.(echox.ApiError)
			if ok {
				return apiErr.JSON(c)
			}

			return c.JSON(http.StatusInternalServerError, LoginResponse{
				Success:    false,
				Reason:     "something went wrong",
				ReasonCode: reasons.INTERNAL_ERROR,
			})
		}

		return c.JSON(http.StatusOK, LoginResponse{
			Success: true,
			Token:   &signedToken,
		})
	}
}

func (m Module) LogoutHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		if err := m.repo.auth.RemoveUserSession(ctx, "username"); err != nil {
			log.Err(err).Msg("failed to remove JWT session from redis cache")

			return c.JSON(http.StatusInternalServerError, LoginResponse{
				Success:    false,
				Reason:     "something went wrong",
				ReasonCode: reasons.INTERNAL_ERROR,
			})
		}

		c.SetCookie(&http.Cookie{
			Name:     fmt.Sprintf("%s-token", m.config.GetIssuerName()),
			Value:    "",
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
			MaxAge:   -1,
		})

		return c.JSON(http.StatusNotImplemented, nil)
	}
}
