package handler

import (
	"library-service/db"
	"library-service/web/middlewares"
	"library-service/web/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	var (
		user         db.LoggedUser
		err          error
		accessToken  string
		refreshToken string
		id           int
		email        string
	)

	id, _ = strconv.Atoi(r.Header.Get("id"))
	email = r.Header.Get("email")
	accessTokenIat := jwt.NewNumericDate(time.Now())
	refreshTokenIat := accessTokenIat
	accessTokenExp := jwt.NewNumericDate(time.Now().Add(time.Minute * -1))
	refreshTokenExp := jwt.NewNumericDate(time.Now().Add(time.Minute * -1))
	claims := middlewares.AuthClaims{
		Id:    id,
		Email: email,
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
		"login-status":  "logged out",
		"access-token":  "Bearer " + accessToken,
		"refresh-token": refreshToken,
	})

}
