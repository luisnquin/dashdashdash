package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/luisnquin/dashdashdash/internal/config"
	"github.com/luisnquin/dashdashdash/internal/helpers/echox"
	"github.com/luisnquin/dashdashdash/internal/helpers/reasons"
	"github.com/luisnquin/go-log"
	"github.com/redis/go-redis/v9"
)

const UNFORTUNATELLY_ONLY_SESSION_OWNER_ALLOWED = "unfortunatelly only the session owner is allowed to be log in but don't worry, he'll be notified :)"

func (m Module) AuthCheckMiddleware() echo.MiddlewareFunc {
	return authCheckMiddleware(m.config, m.repo.auth.db, m.repo.auth.cache)
}

func authCheckMiddleware(config *config.Config, db *sqlx.DB, cache *redis.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authRepo := NewRepository(db, cache)

			var cookie *http.Cookie

			getHasBasicAuth := func() bool {
				username, _, ok := c.Request().BasicAuth()

				return ok && username == os.Getenv("USER")
			}

			if config.IsDevelopment && getHasBasicAuth() { // otherwise TOTP could be skipped
				signedToken, err := doBasicAuth(c, config, authRepo)
				if err != nil {
					apiErr, ok := err.(echox.ApiError)
					if ok {
						return apiErr.JSON(c)
					}

					return c.JSON(http.StatusInternalServerError, AuthMiddlewareResponse{
						Success:    false,
						Reason:     "something went wrong",
						ReasonCode: reasons.INTERNAL_ERROR,
					})
				}

				cookie = &http.Cookie{
					Value: signedToken,
				}
			} else if authHeaderValue := c.Request().Header.Get("Authorization"); authHeaderValue != "" {
				cookie = &http.Cookie{
					Value: strings.TrimPrefix(authHeaderValue, "Bearer "),
				}
			} else {
				var err error

				cookie, err = c.Cookie(fmt.Sprintf("%s-token", config.GetIssuerName()))
				if err != nil && !errors.Is(err, http.ErrNoCookie) {
					log.Err(err).Msg("error getting session cookie")

					return c.JSON(http.StatusInternalServerError, AuthMiddlewareResponse{
						Success:    false,
						Reason:     "something went wrong",
						ReasonCode: reasons.INTERNAL_ERROR,
					})
				}
			}

			if cookie == nil || cookie.Value == "" {
				return c.JSON(http.StatusUnauthorized, AuthMiddlewareResponse{
					Success:    false,
					Reason:     "no auth token provided",
					ReasonCode: reasons.TOKEN_MISSING,
				})
			}

			token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
				}

				return config.Auth.GetJWTSecret(), nil
			})
			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					l := log.Debug()

					if claims, ok := token.Claims.(jwt.MapClaims); ok {
						l.Any("claims", claims)
					}

					l.Msg("token has expired")

					return c.JSON(http.StatusUnauthorized, AuthMiddlewareResponse{
						Success:    false,
						Reason:     "token has expired",
						ReasonCode: reasons.TOKEN_EXPIRED,
					})
				}

				log.Err(err).Msg("failed to parse JWT token")

				return c.JSON(http.StatusUnauthorized, AuthMiddlewareResponse{
					Success:    false,
					Reason:     "invalid token",
					ReasonCode: reasons.TOKEN_INVALID,
				})
			}

			username, ok := token.Claims.(jwt.MapClaims)["username"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, AuthMiddlewareResponse{
					Success:    false,
					Reason:     "invalid auth token",
					ReasonCode: reasons.TOKEN_INVALID,
				})
			}

			if username != os.Getenv("USER") {
				return c.JSON(http.StatusUnauthorized, AuthMiddlewareResponse{
					Success:    false,
					Reason:     UNFORTUNATELLY_ONLY_SESSION_OWNER_ALLOWED,
					ReasonCode: reasons.SESSION_IMPOSIBLE,
				})
			}

			user, err := authRepo.FindOneUserByUsername(c.Request().Context(), username)
			if err != nil {
				log.Warn().Err(err).Msg("error getting user from repository")

				if errors.Is(err, sql.ErrNoRows) {
					return c.JSON(http.StatusUnauthorized, AuthMiddlewareResponse{
						Success:    false,
						Reason:     fmt.Sprintf("user '%s' not found", username),
						ReasonCode: reasons.SESSION_IMPOSIBLE,
					})
				}

				return c.JSON(http.StatusInternalServerError, AuthMiddlewareResponse{
					Success:    false,
					Reason:     "something went wrong",
					ReasonCode: reasons.INTERNAL_ERROR,
				})
			}

			c.Set("user", &user)

			return next(c)
		}
	}
}
