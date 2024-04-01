package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/luisnquin/dashdashdash/internal/config"
	"github.com/luisnquin/dashdashdash/internal/helpers/echox"
	"github.com/luisnquin/dashdashdash/internal/helpers/reasons"
	"github.com/luisnquin/dashdashdash/internal/models"
	"github.com/luisnquin/go-log"
	"golang.org/x/crypto/bcrypt"
)

func doBasicAuth(c echo.Context, config *config.Config, authRepo Repository) (signedToken string, err error) {
	ctx := c.Request().Context()

	username, password, ok := c.Request().BasicAuth()
	if !ok {
		return "", echox.ApiError{
			StatusCode: http.StatusUnauthorized,
			Data: LoginResponse{
				Success:    false,
				Reason:     "no basic auth credentials provided",
				ReasonCode: reasons.SESSION_MISSING_CREDS,
			},
		}
	}

	user, err := authRepo.FindOneUserByUsername(ctx, username)
	if err != nil {
		log.Warn().Err(err).Msg("error getting user from repository")

		if errors.Is(err, sql.ErrNoRows) {
			return "", echox.ApiError{
				StatusCode: http.StatusUnauthorized,
				Data: LoginResponse{
					Success:    false,
					Reason:     fmt.Sprintf("user '%s' not found", username),
					ReasonCode: reasons.SESSION_MISSING_USER,
				},
			}
		}

		return "", echox.ApiError{
			StatusCode: http.StatusInternalServerError,
			Data: LoginResponse{
				Success:    false,
				Reason:     fmt.Sprintf("unable to find user '%s', try it again later", username),
				ReasonCode: reasons.INTERNAL_ERROR,
			},
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		message := fmt.Sprintf("incorrect password for user '%s'", username)
		log.Debug().Err(err).Msg(message)

		return "", echox.ApiError{
			StatusCode: http.StatusUnauthorized,
			Data: LoginResponse{
				Success:    false,
				Reason:     message,
				ReasonCode: reasons.SESSION_INV_CREDS,
			},
		}
	}

	tokenDuration := config.Auth.GetJWTDuration()

	claims := &models.JWTCustomClaims{
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.GetIssuerName(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString(config.Auth.GetJWTSecret())
	if err != nil {
		log.Err(err).Msg("cannot generate jwt token, check the token or the claims")

		return "", echox.ApiError{
			StatusCode: http.StatusInternalServerError,
			Data: LoginResponse{
				Success:    false,
				Reason:     fmt.Sprintf("unable to generate token for user '%s', try it again later", username),
				ReasonCode: reasons.INTERNAL_ERROR,
			},
		}
	}

	if err := authRepo.SaveUserSession(ctx, username, signedToken, tokenDuration); err != nil {
		log.Err(err).Str("token_was", signedToken).Msg("(after generation) failed to save JWT session")

		return "", echox.ApiError{
			StatusCode: http.StatusInternalServerError,
			Data: LoginResponse{
				Success:    false,
				Reason:     "something went wrong",
				ReasonCode: reasons.INTERNAL_ERROR,
			},
		}
	}

	c.SetCookie(&http.Cookie{
		Name:     fmt.Sprintf("%s-token", config.GetIssuerName()),
		Value:    signedToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(tokenDuration),
	})

	return signedToken, nil
}
