package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/luisnquin/dashdashdash/internal/models"
	"github.com/luisnquin/go-log"
	"golang.org/x/crypto/bcrypt"
)

func (m Module) GenerateTOTPUriHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, GenerateTOPTURIResponse{
			URI: m.totp.ProvisioningUri("luisnquin", m.config.GetOPTIssuer()),
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
				IsValid: false,
				Reason:  "code doesn't match with the expected",
			})
		}

		return c.JSON(http.StatusOK, ValidateTOPTCodeResponse{
			IsValid: true,
		})
	}
}

func (m Module) LoginHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		username, password, ok := c.Request().BasicAuth()
		if !ok {
			return c.JSON(http.StatusUnauthorized, LoginResponse{
				Success: false,
				Reason:  "no basic auth credentials provided",
			})
		}

		ctx := c.Request().Context()

		user, err := m.repo.auth.FindOneUserByUsername(ctx, username)
		if err != nil {
			log.Warn().Err(err).Msg("error getting user from repository")

			if errors.Is(err, sql.ErrNoRows) {
				return c.JSON(http.StatusUnauthorized, LoginResponse{
					Success: false,
					Reason:  fmt.Sprintf("user '%s' not found", username),
				})
			} else {
				return c.JSON(http.StatusInternalServerError, LoginResponse{
					Success: false,
					Reason:  fmt.Sprintf("unable to find user '%s', try it again later", username),
				})
			}
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			message := fmt.Sprintf("incorrect password for user '%s'", username)
			log.Debug().Err(err).Msg(message)

			return c.JSON(http.StatusUnauthorized, LoginResponse{
				Success: false,
				Reason:  message,
			})
		}

		tokenDuration := m.config.GetJWTDuration()

		claims := &models.JWTCustomClaims{
			Username: user.Username,
			Email:    user.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    m.config.GetJWTIssuer(),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration)),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		signedToken, err := token.SignedString(m.config.GetJWTSecret())
		if err != nil {
			log.Err(err).Msg("cannot generate jwt token, check the token or the claims")

			return c.JSON(http.StatusInternalServerError, LoginResponse{
				Success: false,
				Reason:  fmt.Sprintf("unable to generate token for user '%s', try it again later", username),
			})
		}

		if err := m.repo.auth.SaveUserSession(ctx, username, signedToken, tokenDuration); err != nil {
			log.Err(err).Str("token_was", signedToken).Msg("(after generation) failed to save JWT session")

			return c.JSON(http.StatusInternalServerError, LoginResponse{
				Success: false,
				Reason:  "something went wrong",
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
		return c.JSON(http.StatusNotImplemented, nil)
	}
}
