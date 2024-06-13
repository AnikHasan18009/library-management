package middlewares

import (
	utils "library-service/web/utils"
	"net/http"
)

func RefreshTokenVerification(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("refresh-token")

		if _, err := VerifyToken(token); err != nil {
			utils.SendError(w, http.StatusUnauthorized, "refresh token verification failed", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
