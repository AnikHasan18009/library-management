package handler

import (
	"encoding/json"
	"library-service/db"
	"library-service/logger"
	"library-service/web/utils"
	"log/slog"
	"net/http"
)

func ApproveReader(w http.ResponseWriter, r *http.Request) {

	data := struct {
		Email string `json:"email"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		slog.Error("error decoding json", logger.Extra(map[string]any{
			"error": err.Error(),
		}))

	}

	statusCode, err := db.GetUserRepo().ApproveReader(data.Email)

	if err != nil {
		utils.SendError(w, statusCode, err.Error(), nil)
	} else {
		utils.SendJson(w, statusCode, map[string]interface{}{
			"status":  true,
			"message": "Success",
		})
	}

}
