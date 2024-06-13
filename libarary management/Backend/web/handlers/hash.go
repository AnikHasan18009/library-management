package handler

import (
	"library-service/web/utils"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Pass struct {
	Val string `json:"password"`
}

func GenerateHashedPassword(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
func HashVal(w http.ResponseWriter, r *http.Request) {

	pass := Pass{}
	utils.DecodeJSON(r.Body, &pass)
	err := bcrypt.CompareHashAndPassword([]byte("$2a$10$SSILNALcismgCm/3OtHZfuaOzAr7L0shyZtIDe3ZzOXGy2FsUVhUa"), []byte(pass.Val))
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SendData(w, pass)

}
