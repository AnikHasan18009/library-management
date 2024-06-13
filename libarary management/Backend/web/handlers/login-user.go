package handler

import (
	"database/sql"
	"fmt"
	"library-service/config"
	"library-service/db"
	"library-service/logger"
	"library-service/web/middlewares"
	"library-service/web/utils"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CreateJWT(authClaims *middlewares.AuthClaims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, *authClaims)
	signedToken, err := token.SignedString([]byte(config.GetConfig().JwtSecret))
	if err != nil {
		slog.Error("error generating jwt", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": *authClaims,
		}))
	}
	return signedToken, err
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var (
		user db.LoggedUser
		err  error

		status       int
		id           int
		accessToken  string
		refreshToken string
	)

	utils.DecodeJSON(r.Body, &user)

	if err = utils.ValidateStructTypes(user); err == nil {

		if id, err = db.GetUserRepo().VerifyUserCredentials(user); id > 0 {
			accessTokenIat := jwt.NewNumericDate(time.Now())
			refreshTokenIat := accessTokenIat
			accessTokenExp := jwt.NewNumericDate(time.Now().Add(time.Minute * 30))
			refreshTokenExp := jwt.NewNumericDate(time.Now().Add(time.Minute * 60))
			claims := middlewares.AuthClaims{
				Id:    id,
				Email: user.Email,
				RegisteredClaims: jwt.RegisteredClaims{
					IssuedAt:  accessTokenIat,
					ExpiresAt: accessTokenExp,
				},
			}
			accessToken, err = CreateJWT(&claims)
			if err != nil {
				utils.SendError(w, http.StatusInternalServerError, err.Error(), user)
				return
			}
			claims.RegisteredClaims.ExpiresAt = refreshTokenExp
			claims.RegisteredClaims.IssuedAt = refreshTokenIat
			refreshToken, err = CreateJWT(&claims)
			if err != nil {
				utils.SendError(w, http.StatusInternalServerError, err.Error(), user)
				return
			}
			utils.SendData(w, map[string]any{
				"login-status":  "logged in",
				"access-token":  "Bearer " + accessToken,
				"refresh-token": refreshToken,
			})
			return
		}
		if err == bcrypt.ErrMismatchedHashAndPassword {
			err = fmt.Errorf("wrong password")
			status = http.StatusBadRequest
		} else if err == sql.ErrNoRows {
			err = fmt.Errorf("sign in first")
			status = http.StatusBadRequest
		}
	}

	slog.Error("login error", logger.Extra(map[string]any{
		"error":   err.Error(),
		"payload": user,
	}))
	utils.SendError(w, status, err.Error(), user)

}
