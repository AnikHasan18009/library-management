package handler

import (
	"encoding/base64"
	"encoding/json"
	"log/slog"

	"library-service/logger"
	"library-service/web/middlewares"
	"library-service/web/utils"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func RefreshToken(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("refresh-token")
	payload := strings.Split(token, ".")[1]

	claims, err := getPayloadData(payload)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	accessTokenIat := jwt.NewNumericDate(time.Now())
	refreshTokenIat := accessTokenIat
	accessTokenExp := jwt.NewNumericDate(time.Now().Add(time.Minute * 5))
	refreshTokenExp := claims.RegisteredClaims.ExpiresAt
	//creating access token
	claims.RegisteredClaims.ExpiresAt = accessTokenExp
	claims.RegisteredClaims.IssuedAt = accessTokenIat
	accessToken, err := CreateJWT(&claims)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	//creating refresh token
	claims.RegisteredClaims.ExpiresAt = refreshTokenExp
	claims.RegisteredClaims.IssuedAt = refreshTokenIat
	refreshToken, err := CreateJWT(&claims)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.SendData(w, map[string]any{
		"access-token":  "Bearer " + accessToken,
		"refresh-token": refreshToken,
	})

}

func getPayloadData(payload string) (middlewares.AuthClaims, error) {
	var claims = middlewares.AuthClaims{}
	var err error
	decodedPayload, _ := base64.RawStdEncoding.DecodeString(payload)
	if err := json.Unmarshal(decodedPayload, &claims); err != nil {
		slog.Error("error unmarshalling claims", logger.Extra(map[string]any{
			"error":         err.Error(),
			"token-payload": payload,
		}))
	}
	return claims, err
}
