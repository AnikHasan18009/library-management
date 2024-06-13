package utils

import (
	"net/http"
	"strconv"
)

func GetUserIdFromPath(r *http.Request) int {
	parameters := []string{"id"}
	parameterValues := GetPathParameterValues(r, parameters)
	id, _ := strconv.Atoi(parameterValues["id"])
	return id
}
