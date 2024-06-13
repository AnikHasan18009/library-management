package handler

import (
	"encoding/json"
	"library-service/db"
	"library-service/logger"
	"library-service/web/utils"
	"log/slog"
	"net/http"
)

func GetNotApprovedReaders(w http.ResponseWriter, r *http.Request) {
	regUser := db.RegisteredUser{}

	err := json.NewDecoder(r.Body).Decode(&regUser)
	if err != nil {
		slog.Error("error decoding json", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": regUser,
		}))
		utils.SendError(w, http.StatusInternalServerError, err.Error(), nil)
		return

	}

	newUser := db.NewUser{
		RegisteredUser: regUser,
		Role:           "reader",
		Approved:       false,
	}

	statusCode, err := db.GetUserRepo().Insert(newUser)

	if err != nil {
		utils.SendError(w, statusCode, err.Error(), nil)
	} else {
		utils.SendJson(w, http.StatusOK, map[string]interface{}{
			"status":  true,
			"message": "Success",
		})
	}

}
