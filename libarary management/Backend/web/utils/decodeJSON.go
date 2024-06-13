package utils

import (
	"encoding/json"
	"io"
	"library-service/logger"
	"log/slog"
)

func DecodeJSON(requestBody io.ReadCloser, userDefinedType interface{}) {

	err := json.NewDecoder(requestBody).Decode(userDefinedType)
	if err != nil {
		slog.Error("error decoding json", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": userDefinedType,
		}))

	}
}
